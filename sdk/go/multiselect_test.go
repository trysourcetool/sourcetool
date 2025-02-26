package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/multiselect"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"
)

func TestConvertStateToMultiSelectProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int32{0, 2}
	defaultValue := []int32{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	multiSelectState := &state.MultiSelectState{
		ID:           id,
		Label:        "Test MultiSelect",
		Value:        value,
		Options:      options,
		Placeholder:  "Select options",
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToMultiSelectProto(multiSelectState)

	if data == nil {
		t.Fatal("convertStateToMultiSelectProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, multiSelectState.Label},
		{"Value", data.Value, multiSelectState.Value},
		{"Options length", len(data.Options), len(multiSelectState.Options)},
		{"Placeholder", data.Placeholder, multiSelectState.Placeholder},
		{"DefaultValue", data.DefaultValue, multiSelectState.DefaultValue},
		{"Required", data.Required, multiSelectState.Required},
		{"Disabled", data.Disabled, multiSelectState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertMultiSelectProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int32{0, 2}
	defaultValue := []int32{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	data := &widgetv1.MultiSelect{
		Label:        "Test MultiSelect",
		Value:        value,
		Options:      options,
		Placeholder:  "Select options",
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertMultiSelectProtoToState(id, data)

	if state == nil {
		t.Fatal("convertMultiSelectProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
		{"Options length", len(state.Options), len(data.Options)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue, data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestMultiSelect(t *testing.T) {
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

	label := "Test MultiSelect"
	options := []string{"Option 1", "Option 2", "Option 3"}
	defaultValue := []string{"Option 1", "Option 3"}
	placeholder := "Select options"

	value := builder.MultiSelect(label,
		multiselect.Options(options...),
		multiselect.DefaultValue(defaultValue...),
		multiselect.Placeholder(placeholder),
		multiselect.Required(true),
		multiselect.Disabled(true),
	)

	if value == nil {
		t.Fatal("MultiSelect returned nil")
	}
	if !reflect.DeepEqual(value.Values, defaultValue) {
		t.Errorf("MultiSelect values = %v, want %v", value.Values, defaultValue)
	}
	expectedIndexes := []int{0, 2}
	if !reflect.DeepEqual(value.Indexes, expectedIndexes) {
		t.Errorf("MultiSelect indexes = %v, want %v", value.Indexes, expectedIndexes)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generateMultiSelectID(label, []int{0})
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
		t.Fatal("MultiSelect state not found")
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
		{"Disabled", state.Disabled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestMultiSelect_WithFormatFunc(t *testing.T) {
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

	label := "Test MultiSelect"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.MultiSelect(label,
		multiselect.Options(options...),
		multiselect.FormatFunc(formatFunc),
	)

	widgetID := builder.generateMultiSelectID(label, []int{0})
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
		t.Fatal("MultiSelect state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	if !reflect.DeepEqual(state.Options, expectedOptions) {
		t.Errorf("Formatted options = %v, want %v", state.Options, expectedOptions)
	}
}
