package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/radio"
)

func TestConvertStateToRadioProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := int32(1)
	defaultValue := int32(0)
	options := []string{"Option 1", "Option 2"}

	radioState := &state.RadioState{
		ID:           id,
		Label:        "Test Radio",
		Value:        &value,
		Options:      options,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToRadioProto(radioState)

	if data == nil {
		t.Fatal("convertStateToRadioProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, radioState.Label},
		{"Value", *data.Value, *radioState.Value},
		{"Options length", len(data.Options), len(radioState.Options)},
		{"DefaultValue", *data.DefaultValue, *radioState.DefaultValue},
		{"Required", data.Required, radioState.Required},
		{"Disabled", data.Disabled, radioState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertRadioProtoToState(t *testing.T) {
	value := int32(1)
	defaultValue := int32(0)
	options := []string{"Option 1", "Option 2"}

	data := &widgetv1.Radio{
		Label:        "Test Radio",
		Value:        &value,
		Options:      options,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertRadioProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertRadioProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, state.ID},
		{"Label", state.Label, data.Label},
		{"Value", *state.Value, *data.Value},
		{"Options length", len(state.Options), len(data.Options)},
		{"DefaultValue", *state.DefaultValue, *data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestRadio(t *testing.T) {
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

	label := "Test Radio"
	options := []string{"Option 1", "Option 2"}
	defaultValue := "Option 1"

	value := builder.Radio(label,
		radio.WithOptions(options...),
		radio.WithDefaultValue(defaultValue),
		radio.WithRequired(true),
	)

	if value == nil {
		t.Fatal("Radio returned nil")
	}
	if value.Value != defaultValue {
		t.Errorf("Radio value = %v, want %v", value.Value, defaultValue)
	}
	if value.Index != 0 {
		t.Errorf("Radio index = %v, want 0", value.Index)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeRadio, []int{0})
	state := sess.State.GetRadio(widgetID)
	if state == nil {
		t.Fatal("Radio state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Options length", len(state.Options), len(options)},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestRadio_WithFormatFunc(t *testing.T) {
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

	label := "Test Radio"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.Radio(label,
		radio.WithOptions(options...),
		radio.WithFormatFunc(formatFunc),
	)

	widgetID := builder.generatePageID(state.WidgetTypeRadio, []int{0})
	state := sess.State.GetRadio(widgetID)
	if state == nil {
		t.Fatal("Radio state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	for i, opt := range state.Options {
		if opt != expectedOptions[i] {
			t.Errorf("Formatted option[%d] = %v, want %v", i, opt, expectedOptions[i])
		}
	}
}
