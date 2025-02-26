package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/textinput"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"
)

func TestConvertStateToTextInputProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := int32(100)
	minLength := int32(10)

	textInputState := &state.TextInputState{
		ID:           id,
		Label:        "Test TextInput",
		Value:        conv.NilValue("test value"),
		Placeholder:  "Enter text",
		DefaultValue: conv.NilValue("default"),
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
	}

	data := convertStateToTextInputProto(textInputState)

	if data == nil {
		t.Fatal("convertStateToTextInputProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, textInputState.Label},
		{"Value", conv.SafeValue(data.Value), conv.SafeValue(textInputState.Value)},
		{"Placeholder", data.Placeholder, textInputState.Placeholder},
		{"DefaultValue", conv.SafeValue(data.DefaultValue), conv.SafeValue(textInputState.DefaultValue)},
		{"Required", data.Required, textInputState.Required},
		{"Disabled", data.Disabled, textInputState.Disabled},
		{"MaxLength", *data.MaxLength, *textInputState.MaxLength},
		{"MinLength", *data.MinLength, *textInputState.MinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTextInputProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := int32(100)
	minLength := int32(10)

	data := &widgetv1.TextInput{
		Label:        "Test TextInput",
		Value:        conv.NilValue("test value"),
		Placeholder:  "Enter text",
		DefaultValue: conv.NilValue("default"),
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
	}

	state := convertTextInputProtoToState(id, data)

	if state == nil {
		t.Fatal("convertTextInputProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", conv.SafeValue(state.Value), conv.SafeValue(data.Value)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", conv.SafeValue(state.DefaultValue), conv.SafeValue(data.DefaultValue)},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"MaxLength", *state.MaxLength, *data.MaxLength},
		{"MinLength", *state.MinLength, *data.MinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextInput(t *testing.T) {
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

	label := "Test TextInput"
	defaultValue := "default value"
	placeholder := "Enter text"
	maxLength := int32(100)
	minLength := int32(10)

	value := builder.TextInput(label,
		textinput.DefaultValue(defaultValue),
		textinput.Placeholder(placeholder),
		textinput.Required(true),
		textinput.Disabled(true),
		textinput.MaxLength(maxLength),
		textinput.MinLength(minLength),
	)

	if value != defaultValue {
		t.Errorf("TextInput value = %v, want %v", value, defaultValue)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generateTextInputID(label, []int{0})
	state := sess.State.GetTextInput(widgetID)
	if state == nil {
		t.Fatal("TextInput state not found")
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
