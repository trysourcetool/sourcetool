package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/textarea"
)

func TestConvertStateToTextAreaProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := int32(1000)
	minLength := int32(10)
	maxLines := int32(10)
	minLines := int32(3)

	textAreaState := &state.TextAreaState{
		ID:           id,
		Label:        "Test TextArea",
		Value:        conv.NilValue("test value"),
		Placeholder:  "Enter text",
		DefaultValue: conv.NilValue("default"),
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
		MaxLines:     &maxLines,
		MinLines:     &minLines,
		AutoResize:   true,
	}

	data := convertStateToTextAreaProto(textAreaState)

	if data == nil {
		t.Fatal("convertStateToTextAreaProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, textAreaState.Label},
		{"Value", data.Value, textAreaState.Value},
		{"Placeholder", data.Placeholder, textAreaState.Placeholder},
		{"DefaultValue", data.DefaultValue, textAreaState.DefaultValue},
		{"Required", data.Required, textAreaState.Required},
		{"Disabled", data.Disabled, textAreaState.Disabled},
		{"MaxLength", *data.MaxLength, *textAreaState.MaxLength},
		{"MinLength", *data.MinLength, *textAreaState.MinLength},
		{"MaxLines", *data.MaxLines, *textAreaState.MaxLines},
		{"MinLines", *data.MinLines, *textAreaState.MinLines},
		{"AutoResize", data.AutoResize, textAreaState.AutoResize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTextAreaProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := int32(1000)
	minLength := int32(10)
	maxLines := int32(10)
	minLines := int32(3)

	data := &widgetv1.TextArea{
		Label:        "Test TextArea",
		Value:        conv.NilValue("test value"),
		Placeholder:  "Enter text",
		DefaultValue: conv.NilValue("default"),
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
		MaxLines:     &maxLines,
		MinLines:     &minLines,
		AutoResize:   true,
	}

	state := convertTextAreaProtoToState(id, data)

	if state == nil {
		t.Fatal("convertTextAreaProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue, data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"MaxLength", *state.MaxLength, *data.MaxLength},
		{"MinLength", *state.MinLength, *data.MinLength},
		{"MaxLines", *state.MaxLines, *data.MaxLines},
		{"MinLines", *state.MinLines, *data.MinLines},
		{"AutoResize", state.AutoResize, data.AutoResize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextArea(t *testing.T) {
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

	label := "Test TextArea"
	defaultValue := "default value"
	placeholder := "Enter text"
	maxLength := int32(1000)
	minLength := int32(10)
	maxLines := int32(10)
	minLines := int32(3)

	value := builder.TextArea(label,
		textarea.DefaultValue(defaultValue),
		textarea.Placeholder(placeholder),
		textarea.Required(true),
		textarea.Disabled(true),
		textarea.MaxLength(maxLength),
		textarea.MinLength(minLength),
		textarea.MaxLines(maxLines),
		textarea.MinLines(minLines),
		textarea.AutoResize(false),
	)

	if value != defaultValue {
		t.Errorf("TextArea value = %v, want %v", value, defaultValue)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeTextArea, []int{0})
	state := sess.State.GetTextArea(widgetID)
	if state == nil {
		t.Fatal("TextArea state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", conv.SafeValue(state.Value), defaultValue},
		{"Placeholder", state.Placeholder, placeholder},
		{"DefaultValue", conv.SafeValue(state.DefaultValue), defaultValue},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
		{"MaxLength", *state.MaxLength, maxLength},
		{"MinLength", *state.MinLength, minLength},
		{"MaxLines", *state.MaxLines, maxLines},
		{"MinLines", *state.MinLines, minLines},
		{"AutoResize", state.AutoResize, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextArea_DefaultMinLines(t *testing.T) {
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

	label := "Test TextArea"

	builder.TextArea(label)

	widgetID := builder.generatePageID(state.WidgetTypeTextArea, []int{0})
	state := sess.State.GetTextArea(widgetID)
	if state == nil {
		t.Fatal("TextArea state not found")
	}

	if state.MinLines == nil {
		t.Fatal("MinLines is nil, want 2")
	}
	if *state.MinLines != 2 {
		t.Errorf("Default MinLines = %v, want 2", *state.MinLines)
	}
}
