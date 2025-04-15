package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/ws"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = pingPeriod * 2

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512KB

	// Maximum time to wait for connection recovery.
	maxRecoveryWait = 6 * time.Hour
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
	service  ws.Service
}

func NewWebSocketHandler(upgrader websocket.Upgrader, service ws.Service) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: upgrader,
		service:  service,
	}
}

func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to upgrade connection: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	ctx := ctxutil.NewBackgroundContext(r.Context())

	done := make(chan struct{})
	defer func() {
		logger.Logger.Info("Closing connection")
		close(done)
		conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Sugar().Debugf("Connection closed: %v", err)
			break
		}

		var msg websocketv1.Message
		if err := proto.Unmarshal(data, &msg); err != nil {
			logger.Logger.Sugar().Errorf("Failed to unmarshal message: %v", err)
			break
		}

		switch msg.Type.(type) {
		case *websocketv1.Message_InitializeHost:
			instanceID := r.Header.Get("X-Instance-Id")
			hostInstance, err := h.service.InitializeHost(ctx, conn, instanceID, &msg)
			if err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}

			defer func() {
				if err := h.updateHostInstanceStatus(ctx, hostInstance.ID, model.HostInstanceStatusOffline); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host instance status offline: %v", err)
				}
				ws.GetConnManager().DisconnectHost(hostInstance.ID)
			}()

			go h.pingPongHostInstanceLoop(ctx, conn, done, hostInstance)
		case *websocketv1.Message_InitializeClient:
			if err := h.service.InitializeClient(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RenderWidget:
			if err := h.service.RenderWidget(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_RerunPage:
			if err := h.service.RerunPage(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_CloseSession:
			if err := h.service.CloseSession(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_ScriptFinished:
			if err := h.service.ScriptFinished(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		case *websocketv1.Message_Exception:
			if err := h.service.Exception(ctx, conn, &msg); err != nil {
				ws.SendErrResponse(ctx, conn, msg.Id, err)
				continue
			}
		default:
			logger.Logger.Sugar().Errorf("Unknown method: %s", msg.Type)
			continue
		}
	}
}

func (h *WebSocketHandler) updateHostInstanceStatus(ctx context.Context, hostInstanceID uuid.UUID, status model.HostInstanceStatus) error {
	if _, err := h.service.UpdateStatus(ctx, dto.UpdateHostInstanceStatusInput{
		ID:     hostInstanceID.String(),
		Status: status,
	}); err != nil {
		return err
	}

	return nil
}

func (h *WebSocketHandler) pingPongHostInstanceLoop(ctx context.Context, conn *websocket.Conn, done chan struct{}, hostInstance *model.HostInstance) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Logger.Info("Stopped ping ticker")
		ticker.Stop()
	}()

	var firstFailureTime *time.Time
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				logger.Logger.Sugar().Errorf("Failed to write ping message: %v", err)
				// Record the first failure time if not already set
				if firstFailureTime == nil {
					now := time.Now()
					firstFailureTime = &now
					logger.Logger.Sugar().Infof("Recording first ping failure time: %v", now)
				}

				// Check if we've exceeded the maximum recovery wait time
				if time.Since(*firstFailureTime) > maxRecoveryWait {
					logger.Logger.Sugar().Infof("Connection unrecoverable after %v", maxRecoveryWait)
					return
				}
				if hostInstance.Status != model.HostInstanceStatusUnreachable {
					if err := h.updateHostInstanceStatus(ctx, hostInstance.ID, model.HostInstanceStatusUnreachable); err != nil {
						logger.Logger.Sugar().Errorf("Failed to update host instance status unreachable: %v", err)
					}
				}
				continue
			}
			// Reset failure time if ping succeeds
			if firstFailureTime != nil {
				logger.Logger.Info("Connection recovered, resetting failure time")
				firstFailureTime = nil
			}
			if hostInstance.Status != model.HostInstanceStatusOnline {
				if err := h.updateHostInstanceStatus(ctx, hostInstance.ID, model.HostInstanceStatusOnline); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host instance status online: %v", err)
				}
			}
		}
	}
}
