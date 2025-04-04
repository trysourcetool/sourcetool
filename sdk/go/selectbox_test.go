package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/selectbox"
)

func TestConvertStateToSelectboxProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := int32(1)
	defaultValue := int32(0)
	options := []string{"Option 1", "Option 2"}
	placeholder := "Select an option"

	selectboxState := &state.SelectboxState{
		ID:           id,
		Label:        "Test Selectbox",
		Value:        &value,
		Options:      options,
		Placeholder:  placeholder,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToSelectboxProto(selectboxState)

	if data == nil {
		t.Fatal("convertStateToSelectboxProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, selectboxState.Label},
		{"Value", *data.Value, *selectboxState.Value},
		{"Options length", len(data.Options), len(selectboxState.Options)},
		{"Placeholder", data.Placeholder, selectboxState.Placeholder},
		{"DefaultValue", *data.DefaultValue, *selectboxState.DefaultValue},
		{"Required", data.Required, selectboxState.Required},
		{"Disabled", data.Disabled, selectboxState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertSelectboxProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := int32(1)
	defaultValue := int32(0)
	options := []string{"Option 1", "Option 2"}
	placeholder := "Select an option"

	data := &widgetv1.Selectbox{
		Label:        "Test Selectbox",
		Value:        &value,
		Options:      options,
		Placeholder:  placeholder,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertSelectboxProtoToState(id, data)

	if state == nil {
		t.Fatal("convertSelectboxProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", *state.Value, *data.Value},
		{"Options length", len(state.Options), len(data.Options)},
		{"Placeholder", state.Placeholder, data.Placeholder},
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

func TestSelectbox(t *testing.T) {
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

	label := "Test Selectbox"
	options := []string{"Option 1", "Option 2"}
	defaultValue := "Option 1"
	placeholder := "Select an option"

	value := builder.Selectbox(label,
		selectbox.Options(options...),
		selectbox.DefaultValue(defaultValue),
		selectbox.Placeholder(placeholder),
		selectbox.Required(true),
	)

	if value == nil {
		t.Fatal("Selectbox returned nil")
	}
	if value.Value != defaultValue {
		t.Errorf("Selectbox value = %v, want %v", value.Value, defaultValue)
	}
	if value.Index != 0 {
		t.Errorf("Selectbox index = %v, want 0", value.Index)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeSelectbox, []int{0})
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		t.Fatal("Selectbox state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Options length", len(state.Options), len(options)},
		{"Placeholder", state.Placeholder, placeholder},
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

func TestSelectbox_WithFormatFunc(t *testing.T) {
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

	label := "Test Selectbox"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.Selectbox(label,
		selectbox.Options(options...),
		selectbox.FormatFunc(formatFunc),
	)

	widgetID := builder.generatePageID(state.WidgetTypeSelectbox, []int{0})
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		t.Fatal("Selectbox state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	for i, opt := range state.Options {
		if opt != expectedOptions[i] {
			t.Errorf("Formatted option[%d] = %v, want %v", i, opt, expectedOptions[i])
		}
	}
}
