package sourcetool

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"

	"github.com/trysourcetool/sourcetool-go/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToDateInputProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)

	dateInputState := &state.DateInputState{
		ID:           id,
		Label:        "Test DateInput",
		Value:        &now,
		Placeholder:  "Select date",
		DefaultValue: &now,
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     &maxDate,
		MinValue:     &minDate,
		Location:     time.Local,
	}

	data := convertStateToDateInputProto(dateInputState)

	if data == nil {
		t.Fatal("convertStateToDateInputProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, dateInputState.Label},
		{"Value", conv.SafeValue(data.Value), dateInputState.Value.Format(time.DateOnly)},
		{"Placeholder", data.Placeholder, dateInputState.Placeholder},
		{"DefaultValue", conv.SafeValue(data.DefaultValue), dateInputState.DefaultValue.Format(time.DateOnly)},
		{"Required", data.Required, dateInputState.Required},
		{"Disabled", data.Disabled, dateInputState.Disabled},
		{"Format", data.Format, dateInputState.Format},
		{"MaxValue", data.MaxValue, dateInputState.MaxValue.Format(time.DateOnly)},
		{"MinValue", data.MinValue, dateInputState.MinValue.Format(time.DateOnly)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertDateInputProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	dateStr := now.Format(time.DateOnly)
	maxDateStr := now.AddDate(1, 0, 0).Format(time.DateOnly)
	minDateStr := now.AddDate(-1, 0, 0).Format(time.DateOnly)

	data := &widgetv1.DateInput{
		Label:        "Test DateInput",
		Value:        conv.NilValue(dateStr),
		Placeholder:  "Select date",
		DefaultValue: conv.NilValue(dateStr),
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     maxDateStr,
		MinValue:     minDateStr,
	}

	state, err := convertDateInputProtoToState(id, data, time.Local)
	if err != nil {
		t.Fatalf("convertDateInputProtoToState returned error: %v", err)
	}

	if state == nil {
		t.Fatal("convertDateInputProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value.Format(time.DateOnly), conv.SafeValue(data.Value)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue.Format(time.DateOnly), conv.SafeValue(data.DefaultValue)},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"Format", state.Format, data.Format},
		{"MaxValue", state.MaxValue.Format(time.DateOnly), data.MaxValue},
		{"MinValue", state.MinValue.Format(time.DateOnly), data.MinValue},
		{"Location", state.Location, time.Local},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertDateInputProtoToState_InvalidDate(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := &widgetv1.DateInput{
		Value: conv.NilValue("invalid-date"),
	}

	_, err := convertDateInputProtoToState(id, data, time.Local)
	if err == nil {
		t.Error("Expected error for invalid date, got nil")
	}
}

func TestDateInput(t *testing.T) {
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

	label := "Test DateInput"
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)
	placeholder := "Select date"
	format := "YYYY-MM-DD"
	location := *time.UTC

	value := builder.DateInput(label,
		dateinput.DefaultValue(now),
		dateinput.Placeholder(placeholder),
		dateinput.Required(true),
		dateinput.Disabled(true),
		dateinput.Format(format),
		dateinput.MaxValue(maxDate),
		dateinput.MinLength(minDate),
		dateinput.Location(location),
	)

	if value == nil {
		t.Fatal("DateInput returned nil")
	}
	if !value.Equal(now) {
		t.Errorf("DateInput value = %v, want %v", value, now)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeDateInput, []int{0})
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		t.Fatal("DateInput state not found")
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
		{"Format", state.Format, format},
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
	if !state.MaxValue.Equal(maxDate) {
		t.Errorf("MaxValue = %v, want %v", state.MaxValue, maxDate)
	}
	if !state.MinValue.Equal(minDate) {
		t.Errorf("MinValue = %v, want %v", state.MinValue, minDate)
	}
}

func TestDateInput_DefaultValues(t *testing.T) {
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

	label := "Test DateInput"

	builder.DateInput(label)

	widgetID := builder.generatePageID(state.WidgetTypeDateInput, []int{0})
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		t.Fatal("DateInput state not found")
	}

	if state.Format != "YYYY/MM/DD" {
		t.Errorf("Default Format = %v, want YYYY/MM/DD", state.Format)
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
