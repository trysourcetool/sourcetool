package sourcetool

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/timeinput"
)

func TestConvertStateToTimeInputProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()

	timeInputState := &state.TimeInputState{
		ID:           id,
		Label:        "Test TimeInput",
		Value:        &now,
		Placeholder:  "Select time",
		DefaultValue: &now,
		Required:     true,
		Disabled:     false,
		Location:     time.Local,
	}

	data := convertStateToTimeInputProto(timeInputState)

	if data == nil {
		t.Fatal("convertStateToTimeInputProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, timeInputState.Label},
		{"Value", conv.SafeValue(data.Value), timeInputState.Value.Format(time.TimeOnly)},
		{"Placeholder", data.Placeholder, timeInputState.Placeholder},
		{"DefaultValue", conv.SafeValue(data.DefaultValue), timeInputState.DefaultValue.Format(time.TimeOnly)},
		{"Required", data.Required, timeInputState.Required},
		{"Disabled", data.Disabled, timeInputState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTimeInputProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	timeStr := now.Format(time.TimeOnly)

	data := &widgetv1.TimeInput{
		Label:        "Test TimeInput",
		Value:        conv.NilValue(timeStr),
		Placeholder:  "Select time",
		DefaultValue: conv.NilValue(timeStr),
		Required:     true,
		Disabled:     false,
	}

	state, err := convertTimeInputProtoToState(id, data, time.Local)
	if err != nil {
		t.Fatalf("convertTimeInputProtoToState returned error: %v", err)
	}

	if state == nil {
		t.Fatal("convertTimeInputProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value.Format(time.TimeOnly), conv.SafeValue(data.Value)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue.Format(time.TimeOnly), conv.SafeValue(data.DefaultValue)},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"Location", state.Location.String(), time.Local.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTimeInputProtoToState_InvalidTime(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := &widgetv1.TimeInput{
		Value: conv.NilValue("invalid-time"),
	}

	_, err := convertTimeInputProtoToState(id, data, time.Local)
	if err == nil {
		t.Error("Expected error for invalid time, got nil")
	}
}

func TestTimeInput(t *testing.T) {
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

	label := "Test TimeInput"
	now := time.Now()
	placeholder := "Select time"
	location := *time.UTC

	value := builder.TimeInput(label,
		timeinput.DefaultValue(now),
		timeinput.Placeholder(placeholder),
		timeinput.Required(true),
		timeinput.Disabled(true),
		timeinput.Location(location),
	)

	if value == nil {
		t.Fatal("TimeInput returned nil")
	}
	if !value.Equal(now) {
		t.Errorf("TimeInput value = %v, want %v", value, now)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeTimeInput, []int{0})
	state := sess.State.GetTimeInput(widgetID)
	if state == nil {
		t.Fatal("TimeInput state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Placeholder", state.Placeholder, placeholder},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
		{"Location", state.Location.String(), location.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	if !state.Value.Equal(*value) {
		t.Errorf("Value = %v, want %v", state.Value, value)
	}
	if !state.DefaultValue.Equal(now) {
		t.Errorf("DefaultValue = %v, want %v", state.DefaultValue, now)
	}
}

func TestTimeInput_DefaultValues(t *testing.T) {
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

	label := "Test TimeInput"

	builder.TimeInput(label)

	widgetID := builder.generatePageID(state.WidgetTypeTimeInput, []int{0})
	state := sess.State.GetTimeInput(widgetID)
	if state == nil {
		t.Fatal("TimeInput state not found")
	}

	if state.Location != time.Local {
		t.Errorf("Default Location = %v, want %v", state.Location, time.Local)
	}
	if state.Required {
		t.Error("Default Required = true, want false")
	}
	if state.Disabled {
		t.Error("Default Disabled = true, want false")
	}
}
