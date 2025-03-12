package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToMarkdownProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	markdownState := &state.MarkdownState{
		ID:   id,
		Body: "# Test Markdown",
	}

	data := convertStateToMarkdownProto(markdownState)

	if data == nil {
		t.Fatal("convertStateToMarkdownProto returned nil")
	}

	if data.Body != markdownState.Body {
		t.Errorf("Body = %v, want %v", data.Body, markdownState.Body)
	}
}

func TestConvertMarkdownProtoToState(t *testing.T) {
	data := &widgetv1.Markdown{
		Body: "# Test Markdown",
	}

	state := convertMarkdownProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertMarkdownProtoToState returned nil")
	}

	if state.Body != data.Body {
		t.Errorf("Body = %v, want %v", state.Body, data.Body)
	}
}

func TestMarkdown(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	sess := session.New(sessionID, pageID)

	mockWS := mock.NewClient()

	builder := &uiBuilder{
		context: context.Background(),
		session: sess,
		cursor:  newCursor(),
		page: &page{
			id: pageID,
		},
		runtime: &runtime{
			wsClient: mockWS,
		},
	}

	body := "# Test Markdown"
	builder.Markdown(body)

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generateMarkdownID(body, []int{0})
	state := sess.State.GetMarkdown(widgetID)
	if state == nil {
		t.Fatal("Markdown state not found")
	}

	if state.Body != body {
		t.Errorf("Body = %v, want %v", state.Body, body)
	}
}
