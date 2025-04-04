package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkboxgroup"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func (b *uiBuilder) CheckboxGroup(label string, opts ...checkboxgroup.Option) *checkboxgroup.Value {
	checkboxGroupOpts := &options.CheckboxGroupOptions{
		Label:        label,
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, o := range opts {
		o.Apply(checkboxGroupOpts)
	}

	sess := b.session
	if sess == nil {
		return nil
	}
	page := b.page
	if page == nil {
		return nil
	}
	cursor := b.cursor
	if cursor == nil {
		return nil
	}
	path := cursor.getPath()

	var defaultVal []int32
	if len(checkboxGroupOpts.DefaultValue) != 0 {
		for _, o := range checkboxGroupOpts.DefaultValue {
			for i, opt := range checkboxGroupOpts.Options {
				if opt == o {
					defaultVal = append(defaultVal, int32(i))
					break
				}
			}
		}
	}

	widgetID := b.generatePageID(state.WidgetTypeCheckboxGroup, path)
	checkboxGroupState := sess.State.GetCheckboxGroup(widgetID)
	if checkboxGroupState == nil {
		checkboxGroupState = &state.CheckboxGroupState{
			ID:    widgetID,
			Value: defaultVal,
		}
	}
	if checkboxGroupOpts.FormatFunc == nil {
		checkboxGroupOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(checkboxGroupOpts.Options))
	for i, v := range checkboxGroupOpts.Options {
		displayVals[i] = checkboxGroupOpts.FormatFunc(v, i)
	}

	checkboxGroupState.Label = checkboxGroupOpts.Label
	checkboxGroupState.Options = displayVals
	checkboxGroupState.DefaultValue = defaultVal
	checkboxGroupState.Required = checkboxGroupOpts.Required
	checkboxGroupState.Disabled = checkboxGroupOpts.Disabled
	sess.State.Set(widgetID, checkboxGroupState)

	checkboxGroupProto := convertStateToCheckboxGroupProto(checkboxGroupState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_CheckboxGroup{
				CheckboxGroup: checkboxGroupProto,
			},
		},
	})

	cursor.next()

	var value *checkboxgroup.Value
	if checkboxGroupState.Value != nil {
		value = &checkboxgroup.Value{
			Values:  make([]string, len(checkboxGroupState.Value)),
			Indexes: make([]int, len(checkboxGroupState.Value)),
		}
		for i, idx := range checkboxGroupState.Value {
			value.Values[i] = checkboxGroupOpts.Options[idx]
			value.Indexes[i] = int(idx)
		}
	}

	return value
}

func convertStateToCheckboxGroupProto(state *state.CheckboxGroupState) *widgetv1.CheckboxGroup {
	if state == nil {
		return nil
	}
	return &widgetv1.CheckboxGroup{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertCheckboxGroupProtoToState(id uuid.UUID, data *widgetv1.CheckboxGroup) *state.CheckboxGroupState {
	if data == nil {
		return nil
	}
	return &state.CheckboxGroupState{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
