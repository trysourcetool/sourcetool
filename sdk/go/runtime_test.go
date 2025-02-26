package sourcetool

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestRuntime_HandleInitializeClient(t *testing.T) {
	pages := make(map[uuid.UUID]*page)
	pageID := uuid.Must(uuid.NewV4())

	// Test page handler
	handlerCalled := false
	testPage := &page{
		id:   pageID,
		name: "Test Page",
		handler: func(ui UIBuilder) error {
			handlerCalled = true
			return nil
		},
	}
	pages[pageID] = testPage

	mockClient := mock.NewClient()
	r := &runtime{
		wsClient:       mockClient,
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	// Create test message
	sessionID := uuid.Must(uuid.NewV4())
	initClient := &websocketv1.InitializeClient{
		SessionId: conv.NilValue(sessionID.String()),
		PageId:    pageID.String(),
	}

	// Register handler
	mockClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch m := msg.Type.(type) {
		case *websocketv1.Message_InitializeClient:
			return r.handleInitializeClient(m.InitializeClient)
		}
		return nil
	})

	// Send message
	mockClient.Enqueue(uuid.Must(uuid.NewV4()).String(), initClient)

	// Verify that session was created
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		t.Error("session was not created")
	}

	// Verify that page handler was called
	if !handlerCalled {
		t.Error("page handler was not called")
	}
}

func TestRuntime_HandleRerunPage(t *testing.T) {
	pages := make(map[uuid.UUID]*page)
	pageID := uuid.Must(uuid.NewV4())
	sessionID := uuid.Must(uuid.NewV4())

	// Test page handler
	handlerCallCount := 0
	testPage := &page{
		id:   pageID,
		name: "Test Page",
		handler: func(ui UIBuilder) error {
			handlerCallCount++
			return nil
		},
	}
	pages[pageID] = testPage

	mockClient := mock.NewClient()
	r := &runtime{
		wsClient:       mockClient,
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	// Initialize session
	sess := session.New(sessionID, pageID)
	r.sessionManager.SetSession(sess)

	// Create test message
	rerunPage := &websocketv1.RerunPage{
		SessionId: sessionID.String(),
		PageId:    pageID.String(),
	}

	// Register handler
	mockClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch m := msg.Type.(type) {
		case *websocketv1.Message_RerunPage:
			return r.handleRerunPage(m.RerunPage)
		}
		return nil
	})

	// Send message
	mockClient.Enqueue(uuid.Must(uuid.NewV4()).String(), rerunPage)

	// Verify that page handler was called
	if handlerCallCount != 1 {
		t.Errorf("page handler call count = %d, want 1", handlerCallCount)
	}

	// Send message again
	mockClient.Enqueue(uuid.Must(uuid.NewV4()).String(), rerunPage)

	// Verify that page handler was called again
	if handlerCallCount != 2 {
		t.Errorf("page handler call count after rerun = %d, want 2", handlerCallCount)
	}
}

func TestRuntime_HandleCloseSession(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())

	mockClient := mock.NewClient()
	r := &runtime{
		wsClient:       mockClient,
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(make(map[uuid.UUID]*page)),
	}

	// Initialize session
	sess := session.New(sessionID, pageID)
	r.sessionManager.SetSession(sess)

	// Create test message
	closeSession := &websocketv1.CloseSession{
		SessionId: sessionID.String(),
	}

	// Register handler
	mockClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch m := msg.Type.(type) {
		case *websocketv1.Message_CloseSession:
			return r.handleCloseSession(m.CloseSession)
		}
		return nil
	})

	// Send message
	mockClient.Enqueue(uuid.Must(uuid.NewV4()).String(), closeSession)

	// Verify that session was deleted
	if got := r.sessionManager.GetSession(sessionID); got != nil {
		t.Error("session was not deleted")
	}
}
