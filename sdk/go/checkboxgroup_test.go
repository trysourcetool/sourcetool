package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkboxgroup"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToCheckboxGroupProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int32{0, 2}
	defaultValue := []int32{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	checkboxGroupState := &state.CheckboxGroupState{
		ID:           id,
		Label:        "Test CheckboxGroup",
		Value:        value,
		Options:      options,
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToCheckboxGroupProto(checkboxGroupState)

	if data == nil {
		t.Fatal("convertStateToCheckboxGroupProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, checkboxGroupState.Label},
		{"Value", data.Value, checkboxGroupState.Value},
		{"Options length", len(data.Options), len(checkboxGroupState.Options)},
		{"DefaultValue", data.DefaultValue, checkboxGroupState.DefaultValue},
		{"Required", data.Required, checkboxGroupState.Required},
		{"Disabled", data.Disabled, checkboxGroupState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertCheckboxGroupProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int32{0, 2}
	defaultValue := []int32{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	data := &widgetv1.CheckboxGroup{
		Label:        "Test CheckboxGroup",
		Value:        value,
		Options:      options,
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertCheckboxGroupProtoToState(id, data)

	if state == nil {
		t.Fatal("convertCheckboxGroupProtoToState returned nil")
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

func TestCheckboxGroup(t *testing.T) {
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

	label := "Test CheckboxGroup"
	options := []string{"Option 1", "Option 2", "Option 3"}
	defaultValue := []string{"Option 1", "Option 3"}

	value := builder.CheckboxGroup(label,
		checkboxgroup.Options(options...),
		checkboxgroup.DefaultValue(defaultValue...),
		checkboxgroup.Required(true),
		checkboxgroup.Disabled(true),
	)

	if value == nil {
		t.Fatal("CheckboxGroup returned nil")
	}
	if !reflect.DeepEqual(value.Values, defaultValue) {
		t.Errorf("CheckboxGroup values = %v, want %v", value.Values, defaultValue)
	}
	expectedIndexes := []int{0, 2}
	if !reflect.DeepEqual(value.Indexes, expectedIndexes) {
		t.Errorf("CheckboxGroup indexes = %v, want %v", value.Indexes, expectedIndexes)
	}

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeCheckboxGroup, []int{0})
	state := sess.State.GetCheckboxGroup(widgetID)
	if state == nil {
		t.Fatal("CheckboxGroup state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Options length", len(state.Options), len(options)},
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

func TestCheckboxGroup_WithFormatFunc(t *testing.T) {
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

	label := "Test CheckboxGroup"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.CheckboxGroup(label,
		checkboxgroup.Options(options...),
		checkboxgroup.FormatFunc(formatFunc),
	)

	widgetID := builder.generatePageID(state.WidgetTypeCheckboxGroup, []int{0})
	state := sess.State.GetCheckboxGroup(widgetID)
	if state == nil {
		t.Fatal("CheckboxGroup state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	if !reflect.DeepEqual(state.Options, expectedOptions) {
		t.Errorf("Formatted options = %v, want %v", state.Options, expectedOptions)
	}
}
