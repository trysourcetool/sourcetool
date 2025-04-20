package wsmanager

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/logger"
)

const (
	// Maximum number of reconnection attempts for subscribers.
	maxReconnectAttempts = 5

	// Base delay for exponential backoff (in milliseconds).
	baseReconnectDelay = 100
)

// subscribeToHostMessages starts a goroutine to subscribe to messages for hosts from the pub/sub system.
// It handles reconnection logic with exponential backoff.
func (m *manager) subscribeToHostMessages() {
	defer m.wg.Done()

	for attempt := range maxReconnectAttempts {
		if err := m.subscribeToHostMessagesWithRetry(); err != nil {
			if m.ctx.Err() != nil {
				logger.Logger.Sugar().Info("Host message subscriber stopping due to context cancellation during retry.")
				return // Context canceled, stop trying
			}
			if attempt == maxReconnectAttempts-1 {
				logger.Logger.Sugar().Errorf("Failed to subscribe to host messages after %d attempts: %v", maxReconnectAttempts, err)
				return
			}
			// Exponential backoff
			delay := time.Duration(baseReconnectDelay*(1<<attempt)) * time.Millisecond
			logger.Logger.Sugar().Warnf("Retrying host message subscription in %v (attempt %d/%d)", delay, attempt+1, maxReconnectAttempts)
			time.Sleep(delay)
			continue
		}
		// Successful subscription or context canceled during subscription
		return
	}
}

// subscribeToHostMessagesWithRetry attempts to subscribe to the host message channel and process messages.
// It returns an error if the subscription fails or the channel closes unexpectedly.
func (m *manager) subscribeToHostMessagesWithRetry() error {
	ch, err := m.pubsubClient.Subscribe(m.ctx, "host_messages")
	if err != nil {
		return fmt.Errorf("failed to subscribe to host messages: %w", err)
	}
	logger.Logger.Sugar().Info("Subscribed to host messages")

	for {
		select {
		case <-m.ctx.Done():
			logger.Logger.Sugar().Info("Host message subscriber stopping due to context cancellation.")
			return nil // Normal shutdown
		case msg, ok := <-ch:
			if !ok {
				// Check context again to differentiate between unexpected close and shutdown
				if m.ctx.Err() != nil {
					logger.Logger.Sugar().Info("Host message channel closed due to context cancellation.")
					return nil
				}
				return fmt.Errorf("host message channel closed unexpectedly")
			}

			if err := m.processHostMessage(msg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to process host message: %v", err)
				// Continue processing other messages even if one fails
				continue
			}
		}
	}
}

// processHostMessage unmarshals a message from pub/sub and forwards it to the appropriate connected host.
func (m *manager) processHostMessage(msg *port.Message) error {
	// The pubsub layer already unwraps redisv1.RedisMessage and gives us the ID and Payload.
	hostInstanceID, err := uuid.FromString(msg.ID)
	if err != nil {
		return fmt.Errorf("invalid host instance ID: %w", err)
	}

	var protoMsg websocketv1.Message
	if err := proto.Unmarshal(msg.Payload, &protoMsg); err != nil {
		return fmt.Errorf("failed to unmarshal websocket message: %w", err)
	}

	// logger.Logger.Sugar().Debugf("Received message for host: %s, type: %T", hostInstanceID, protoMsg.GetPayload()) // Example Debugging

	m.hostsMutex.RLock()
	host, ok := m.connectedHosts[hostInstanceID]
	m.hostsMutex.RUnlock()

	if !ok {
		// This can happen normally if the host disconnected between message publish and processing
		logger.Logger.Sugar().Debugf("Host %s not found for message ID %s, likely disconnected.", hostInstanceID, protoMsg.Id)
		return nil
	}

	logger.Logger.Sugar().Debugf("Sending message %s to host %s", protoMsg.Id, host.hostInstance.ID)

	data, err := proto.Marshal(&protoMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message for host %s: %w", host.hostInstance.ID, err)
	}

	if err := host.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		// Log the error but attempt to send anyway, or disconnect?
		logger.Logger.Sugar().Warnf("Failed to set write deadline for host %s before sending message %s: %v", host.hostInstance.ID, protoMsg.Id, err)
		// Let the WriteMessage handle the error more definitively
	}

	if err := host.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send message %s to host %s, disconnecting: %v", protoMsg.Id, host.hostInstance.ID, err)
		m.DisconnectHost(hostInstanceID) // Disconnect on write failure
		return fmt.Errorf("failed to send message to host %s: %w", host.hostInstance.ID, err)
	}

	return nil
}

// subscribeToClientMessages starts a goroutine to subscribe to messages for clients from the pub/sub system.
// It handles reconnection logic with exponential backoff.
func (m *manager) subscribeToClientMessages() {
	defer m.wg.Done()

	for attempt := range maxReconnectAttempts {
		if err := m.subscribeToClientMessagesWithRetry(); err != nil {
			if m.ctx.Err() != nil {
				logger.Logger.Sugar().Info("Client message subscriber stopping due to context cancellation during retry.")
				return // Context canceled, stop trying
			}
			if attempt == maxReconnectAttempts-1 {
				logger.Logger.Sugar().Errorf("Failed to subscribe to client messages after %d attempts: %v", maxReconnectAttempts, err)
				return
			}
			// Exponential backoff
			delay := time.Duration(baseReconnectDelay*(1<<attempt)) * time.Millisecond
			logger.Logger.Sugar().Warnf("Retrying client message subscription in %v (attempt %d/%d)", delay, attempt+1, maxReconnectAttempts)
			time.Sleep(delay)
			continue
		}
		// Successful subscription or context canceled during subscription
		return
	}
}

// subscribeToClientMessagesWithRetry attempts to subscribe to the client message channel and process messages.
// It returns an error if the subscription fails or the channel closes unexpectedly.
func (m *manager) subscribeToClientMessagesWithRetry() error {
	ch, err := m.pubsubClient.Subscribe(m.ctx, "client_messages")
	if err != nil {
		return fmt.Errorf("failed to subscribe to client messages: %w", err)
	}
	logger.Logger.Sugar().Info("Subscribed to client messages")

	for {
		select {
		case <-m.ctx.Done():
			logger.Logger.Sugar().Info("Client message subscriber stopping due to context cancellation.")
			return nil // Normal shutdown
		case msg, ok := <-ch:
			if !ok {
				// Check context again to differentiate between unexpected close and shutdown
				if m.ctx.Err() != nil {
					logger.Logger.Sugar().Info("Client message channel closed due to context cancellation.")
					return nil
				}
				return fmt.Errorf("client message channel closed unexpectedly")
			}

			if err := m.processClientMessage(msg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to process client message: %v", err)
				// Continue processing other messages even if one fails
				continue
			}
		}
	}
}

// processClientMessage unmarshals a message from pub/sub and forwards it to the appropriate connected client.
func (m *manager) processClientMessage(msg *port.Message) error {
	sessionID, err := uuid.FromString(msg.ID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	var protoMsg websocketv1.Message
	if err := proto.Unmarshal(msg.Payload, &protoMsg); err != nil {
		return fmt.Errorf("failed to unmarshal websocket message: %w", err)
	}

	// logger.Logger.Sugar().Debugf("Received message for client: %s, type: %T", sessionID, protoMsg.GetPayload()) // Example Debugging

	m.clientsMutex.RLock()
	client, ok := m.connectedClients[sessionID]
	m.clientsMutex.RUnlock()

	if !ok {
		// This can happen normally if the client disconnected between message publish and processing
		logger.Logger.Sugar().Debugf("Client %s not found for message ID %s, likely disconnected.", sessionID, protoMsg.Id)
		return nil
	}

	logger.Logger.Sugar().Debugf("Sending message %s to client %s", protoMsg.Id, client.session.ID)

	data, err := proto.Marshal(&protoMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message for client %s: %w", client.session.ID, err)
	}

	if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		// Log the error but attempt to send anyway, or disconnect?
		logger.Logger.Sugar().Warnf("Failed to set write deadline for client %s before sending message %s: %v", client.session.ID, protoMsg.Id, err)
		// Let the WriteMessage handle the error more definitively
	}

	if err := client.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.Logger.Sugar().Errorf("Failed to send message %s to client %s, disconnecting: %v", protoMsg.Id, client.session.ID, err)
		m.DisconnectClient(sessionID) // Disconnect on write failure
		return fmt.Errorf("failed to send message to client %s: %w", client.session.ID, err)
	}

	return nil
}
