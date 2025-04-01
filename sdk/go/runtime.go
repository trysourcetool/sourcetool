package sourcetool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	exceptionv1 "github.com/trysourcetool/sourcetool/proto/go/exception/v1"
	pagev1 "github.com/trysourcetool/sourcetool/proto/go/page/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"
	"go.uber.org/zap"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/errdefs"
	"github.com/trysourcetool/sourcetool-go/internal/logger"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

type runtime struct {
	wsClient       websocket.Client
	sessionManager *session.SessionManager
	pageManager    *pageManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page) (*runtime, error) {
	r := &runtime{
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	wsClient, err := websocket.NewClient(websocket.Config{
		URL:            endpoint,
		APIKey:         apiKey,
		InstanceID:     uuid.Must(uuid.NewV4()),
		PingInterval:   1 * time.Second,
		ReconnectDelay: 1 * time.Second,
		OnReconnecting: func() {
			logger.Log.Info("Reconnecting...")
		},
		OnReconnected: func() {
			logger.Log.Info("Reconnected!")
			r.sendInitializeHost(apiKey, pages)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %v", err)
	}

	r.wsClient = wsClient
	wsClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch t := msg.Type.(type) {
		case *websocketv1.Message_InitializeClient:
			if err := r.handleInitializeClient(t.InitializeClient); err != nil {
				r.sendException(msg.Id, *t.InitializeClient.SessionId, err)
			}
			return nil
		case *websocketv1.Message_RerunPage:
			if err := r.handleRerunPage(t.RerunPage); err != nil {
				r.sendException(msg.Id, t.RerunPage.SessionId, err)
			}
			return nil
		case *websocketv1.Message_CloseSession:
			if err := r.handleCloseSession(t.CloseSession); err != nil {
				r.sendException(msg.Id, t.CloseSession.SessionId, err)
			}
			return nil
		default:
			return fmt.Errorf("unknown message type: %T", t)
		}
	})

	r.sendInitializeHost(apiKey, pages)

	return r, nil
}

func (r *runtime) sendInitializeHost(apiKey string, pages map[uuid.UUID]*page) {
	pagesPayload := make([]*pagev1.Page, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &pagev1.Page{
			Id:     page.id.String(),
			Name:   page.name,
			Route:  page.route,
			Path:   conv.PathToInt32Slice(page.path),
			Groups: page.accessGroups,
		})
	}

	msg := &websocketv1.InitializeHost{
		ApiKey:     apiKey,
		SdkName:    "sourcetool-go",
		SdkVersion: "0.1.0",
		Pages:      pagesPayload,
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), msg)
	if err != nil {
		logger.Log.Fatal("failed to send initialize host message", zap.Error(err))
	}

	if e := resp.GetException(); e != nil {
		logger.Log.Fatal("initialize host message failed", zap.String("message", e.Message))
	}

	logger.Log.Info("initialize host message sent", zap.Any("response", resp))
}

func (r *runtime) handleInitializeClient(msg *websocketv1.InitializeClient) error {
	if msg.SessionId == nil {
		return errdefs.ErrInvalidParameter(errors.New("session id is required"))
	}
	sessionID, err := uuid.FromString(conv.SafeValue(msg.SessionId))
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	session := session.New(sessionID, pageID)
	r.sessionManager.SetSession(session)

	page := r.pageManager.getPage(pageID)
	if page == nil {
		return errdefs.ErrInternal(fmt.Errorf("page not found: %s", pageID))
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: session,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
			SessionId: sessionID.String(),
			Status:    websocketv1.ScriptFinished_STATUS_FAILURE,
		})

		return errdefs.ErrRunPage(err)
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
		SessionId: sessionID.String(),
		Status:    websocketv1.ScriptFinished_STATUS_SUCCESS,
	})

	return nil
}

func (r *runtime) handleRerunPage(msg *websocketv1.RerunPage) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return errdefs.ErrSessionNotFound(fmt.Errorf("session not found: %s", sessionID))
	}

	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}
	page := r.pageManager.getPage(pageID)
	if page == nil {
		return errdefs.ErrPageNotFound(fmt.Errorf("page not found: %s", pageID))
	}

	if sess.PageID != pageID {
		sess.State.ResetStates()
	}

	newWidgetStates := make(map[uuid.UUID]session.WidgetState)
	for _, widget := range msg.States {
		id, err := uuid.FromString(widget.Id)
		if err != nil {
			return errdefs.ErrInvalidParameter(err)
		}
		switch t := widget.Type.(type) {
		case *widgetv1.Widget_TextInput:
			newWidgetStates[id] = convertTextInputProtoToState(id, t.TextInput)
		case *widgetv1.Widget_NumberInput:
			newWidgetStates[id] = convertNumberInputProtoToState(id, t.NumberInput)
		case *widgetv1.Widget_DateInput:
			state, err := convertDateInputProtoToState(id, t.DateInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_DateTimeInput:
			state, err := convertDateTimeInputProtoToState(id, t.DateTimeInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_TimeInput:
			state, err := convertTimeInputProtoToState(id, t.TimeInput, time.Local)
			if err != nil {
				return errdefs.ErrInvalidParameter(err)
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_Form:
			newWidgetStates[id] = convertFormProtoToState(id, t.Form)
		case *widgetv1.Widget_Button:
			newWidgetStates[id] = convertButtonProtoToState(id, t.Button)
		case *widgetv1.Widget_Markdown:
			newWidgetStates[id] = convertMarkdownProtoToState(id, t.Markdown)
		case *widgetv1.Widget_Columns:
			newWidgetStates[id] = convertColumnsProtoToState(id, t.Columns)
		case *widgetv1.Widget_ColumnItem:
			newWidgetStates[id] = convertColumnItemProtoToState(id, t.ColumnItem)
		case *widgetv1.Widget_Table:
			newWidgetStates[id] = convertTableProtoToState(id, t.Table)
		case *widgetv1.Widget_Selectbox:
			newWidgetStates[id] = convertSelectboxProtoToState(id, t.Selectbox)
		case *widgetv1.Widget_MultiSelect:
			newWidgetStates[id] = convertMultiSelectProtoToState(id, t.MultiSelect)
		case *widgetv1.Widget_Checkbox:
			newWidgetStates[id] = convertCheckboxProtoToState(id, t.Checkbox)
		case *widgetv1.Widget_CheckboxGroup:
			newWidgetStates[id] = convertCheckboxGroupProtoToState(id, t.CheckboxGroup)
		case *widgetv1.Widget_Radio:
			newWidgetStates[id] = convertRadioProtoToState(id, t.Radio)
		case *widgetv1.Widget_TextArea:
			newWidgetStates[id] = convertTextAreaProtoToState(id, t.TextArea)
		default:
			return errdefs.ErrInvalidParameter(fmt.Errorf("unknown widget type: %T", t))
		}
	}

	sess.State.SetStates(newWidgetStates)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: sess,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
			SessionId: sessionID.String(),
			Status:    websocketv1.ScriptFinished_STATUS_FAILURE,
		})

		return errdefs.ErrRunPage(err)
	}

	r.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.ScriptFinished{
		SessionId: sessionID.String(),
		Status:    websocketv1.ScriptFinished_STATUS_SUCCESS,
	})

	sess.State.ResetButtons()

	return nil
}

func (r *runtime) handleCloseSession(msg *websocketv1.CloseSession) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return errdefs.ErrInvalidParameter(err)
	}

	r.sessionManager.DisconnectSession(sessionID)

	return nil
}

func (r *runtime) sendException(id, sessionID string, err error) {
	e, ok := err.(*errdefs.Error)
	if !ok {
		v := errdefs.ErrInternal(err)
		e = v.(*errdefs.Error)
	}

	exception := &exceptionv1.Exception{
		Title:      e.Title,
		Message:    e.Message,
		StackTrace: e.StackTrace(),
		SessionId:  sessionID,
	}

	r.wsClient.Enqueue(id, exception)
}

func (r *runtime) Close() error {
	if r.wsClient == nil {
		return nil
	}
	return r.wsClient.Close()
}
