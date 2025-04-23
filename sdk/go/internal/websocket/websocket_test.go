package websocket

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool-go/internal/logger"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	"github.com/trysourcetool/sourcetool-go/internal/ptrconv"
)

func TestMain(m *testing.M) {
	if err := logger.Init(); err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type testServer struct {
	*httptest.Server
	connCh chan *websocket.Conn
}

func newTestServer() *testServer {
	connCh := make(chan *websocket.Conn, 1)
	s := &testServer{
		connCh: connCh,
	}

	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate API key
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		s.connCh <- conn
	}))

	return s
}

func (s *testServer) WaitForConnection(t *testing.T) *websocket.Conn {
	t.Helper()
	select {
	case conn := <-s.connCh:
		return conn
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for connection")
		return nil
	}
}

func TestNewClient(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:    wsURL,
		APIKey: "test_api_key",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}
}

func TestClient_MessageHandling(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:    wsURL,
		APIKey: "test_api_key",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}

	// Test message handler
	receivedCh := make(chan *websocketv1.Message, 1)
	client.RegisterHandler(func(msg *websocketv1.Message) error {
		receivedCh <- msg
		return nil
	})

	// Send test message
	testMsg := &websocketv1.Message{
		Id: "test_id",
		Type: &websocketv1.Message_InitializeClient{
			InitializeClient: &websocketv1.InitializeClient{
				SessionId: ptrconv.StringPtr("test_session"),
				PageId:    "test_page",
			},
		},
	}

	data, err := marshalMessage(testMsg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	// Wait for message reception
	select {
	case msg := <-receivedCh:
		if msg.Id != testMsg.Id {
			t.Errorf("message ID = %v, want %v", msg.Id, testMsg.Id)
		}
		if _, ok := msg.Type.(*websocketv1.Message_InitializeClient); !ok {
			t.Errorf("unexpected message type: %T", msg.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestClient_EnqueueWithResponse(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	instanceID := uuid.Must(uuid.NewV4())
	client, err := NewClient(Config{
		URL:        wsURL,
		APIKey:     "test_api_key",
		InstanceID: instanceID,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}

	// Start goroutine to wait for response
	go func() {
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("failed to read message: %v", err)
			return
		}

		msg, err := unmarshalMessage(data)
		if err != nil {
			t.Errorf("failed to unmarshal message: %v", err)
			return
		}

		// Send response
		resp := &websocketv1.Message{
			Id: msg.Id,
			Type: &websocketv1.Message_InitializeHostCompleted{
				InitializeHostCompleted: &websocketv1.InitializeHostCompleted{
					HostInstanceId: "test_host_instance_id",
				},
			},
		}

		data, err = marshalMessage(resp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}()

	// Send message and wait for response
	initHost := &websocketv1.InitializeHost{}
	resp, err := client.EnqueueWithResponse("test_id", initHost)
	if err != nil {
		t.Fatalf("failed to get response: %v", err)
	}

	if completed, ok := resp.Type.(*websocketv1.Message_InitializeHostCompleted); !ok {
		t.Errorf("unexpected response type: %T", resp.Type)
	} else if completed.InitializeHostCompleted.HostInstanceId != "test_host_instance_id" {
		t.Errorf("unexpected host instance id: got %s, want test_host_instance_id", completed.InitializeHostCompleted.HostInstanceId)
	}
}

func TestClient_Reconnection(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	reconnected := make(chan struct{})
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	instanceID := uuid.Must(uuid.NewV4())
	client, err := NewClient(Config{
		URL:            wsURL,
		APIKey:         "test_api_key",
		InstanceID:     instanceID,
		PingInterval:   100 * time.Millisecond,
		ReconnectDelay: 100 * time.Millisecond,
		OnReconnected: func() {
			reconnected <- struct{}{}
		},
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Verify initial connection
	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("initial connection not established")
	}

	// Disconnect connection
	conn.Close()

	// Wait for reconnection
	select {
	case <-reconnected:
		// Reconnection successful
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for reconnection")
	}

	// Verify new connection
	conn = s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("reconnection not established")
	}
}
