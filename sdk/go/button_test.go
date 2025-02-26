package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func TestConvertStateToButtonProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())

	buttonState := &state.ButtonState{
		ID:       id,
		Label:    "Test Button",
		Value:    true,
		Disabled: true,
	}

	data := convertStateToButtonProto(buttonState)

	if data == nil {
		t.Fatal("convertStateToButtonProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, buttonState.Label},
		{"Value", data.Value, buttonState.Value},
		{"Disabled", data.Disabled, buttonState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertButtonProtoToState(t *testing.T) {
	data := &widgetv1.Button{
		Label:    "Test Button",
		Value:    true,
		Disabled: true,
	}

	state := convertButtonProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertButtonProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
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

func TestButton(t *testing.T) {
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

	label := "Test Button"

	value := builder.Button(label,
		button.Disabled(true),
	)

	if value {
		t.Error("Button value = true, want false")
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generateButtonInputID(label, []int{0})
	state := sess.State.GetButton(widgetID)
	if state == nil {
		t.Fatal("Button state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", state.Value, false},
		{"Disabled", state.Disabled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestButton_DefaultState(t *testing.T) {
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

	label := "Test Button"

	builder.Button(label)

	widgetID := builder.generateButtonInputID(label, []int{0})
	state := sess.State.GetButton(widgetID)
	if state == nil {
		t.Fatal("Button state not found")
	}

	if state.Value {
		t.Error("Default Value = true, want false")
	}
	if state.Disabled {
		t.Error("Default Disabled = true, want false")
	}
}
