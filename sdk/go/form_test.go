package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func TestConvertStateToFormProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	formState := &state.FormState{
		ID:             id,
		Value:          true,
		ButtonLabel:    "Submit",
		ButtonDisabled: true,
		ClearOnSubmit:  true,
	}

	data := convertStateToFormProto(formState)

	if data == nil {
		t.Fatal("convertStateToFormProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Value", data.Value, formState.Value},
		{"ButtonLabel", data.ButtonLabel, formState.ButtonLabel},
		{"ButtonDisabled", data.ButtonDisabled, formState.ButtonDisabled},
		{"ClearOnSubmit", data.ClearOnSubmit, formState.ClearOnSubmit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertFormProtoToState(t *testing.T) {
	data := &widgetv1.Form{
		Value:          true,
		ButtonLabel:    "Submit",
		ButtonDisabled: true,
		ClearOnSubmit:  true,
	}

	state := convertFormProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertFormProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Value", state.Value, data.Value},
		{"ButtonLabel", state.ButtonLabel, data.ButtonLabel},
		{"ButtonDisabled", state.ButtonDisabled, data.ButtonDisabled},
		{"ClearOnSubmit", state.ClearOnSubmit, data.ClearOnSubmit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestForm(t *testing.T) {
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

	buttonLabel := "Submit"
	childBuilder, submitted := builder.Form(buttonLabel)

	if childBuilder == nil {
		t.Fatal("Form returned nil builder")
	}
	if submitted {
		t.Error("Form returned true for submitted, want false")
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generateFormID([]int{0})
	state := sess.State.GetForm(widgetID)
	if state == nil {
		t.Fatal("Form state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ButtonLabel", state.ButtonLabel, buttonLabel},
		{"ButtonDisabled", state.ButtonDisabled, false}, // default value
		{"ClearOnSubmit", state.ClearOnSubmit, false},   // default value
		{"Value", state.Value, false},                   // default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestForm_WithOptions(t *testing.T) {
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

	buttonLabel := "Submit"
	childBuilder, submitted := builder.Form(buttonLabel,
		form.ButtonDisabled(true),
		form.ClearOnSubmit(true),
	)

	if childBuilder == nil {
		t.Fatal("Form returned nil builder")
	}
	if submitted {
		t.Error("Form returned true for submitted, want false")
	}

	widgetID := builder.generateFormID([]int{0})
	state := sess.State.GetForm(widgetID)
	if state == nil {
		t.Fatal("Form state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ButtonLabel", state.ButtonLabel, buttonLabel},
		{"ButtonDisabled", state.ButtonDisabled, true},
		{"ClearOnSubmit", state.ClearOnSubmit, true},
		{"Value", state.Value, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
