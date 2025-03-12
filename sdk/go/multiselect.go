package sourcetool

import (
	"github.com/gofrs/uuid/v5"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/multiselect"
)

func (b *uiBuilder) MultiSelect(label string, opts ...multiselect.Option) *multiselect.Value {
	multiSelectOpts := &options.MultiSelectOptions{
		Label:        label,
		DefaultValue: nil,
		Placeholder:  "",
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, o := range opts {
		o.Apply(multiSelectOpts)
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
	if len(multiSelectOpts.DefaultValue) != 0 {
		for _, o := range multiSelectOpts.DefaultValue {
			for i, opt := range multiSelectOpts.Options {
				if opt == o {
					defaultVal = append(defaultVal, int32(i))
					break
				}
			}
		}
	}

	widgetID := b.generateMultiSelectID(label, path)
	multiSelectState := sess.State.GetMultiSelect(widgetID)
	if multiSelectState == nil {
		multiSelectState = &state.MultiSelectState{
			ID:    widgetID,
			Value: defaultVal,
		}
	}
	if multiSelectOpts.FormatFunc == nil {
		multiSelectOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(multiSelectOpts.Options))
	for i, v := range multiSelectOpts.Options {
		displayVals[i] = multiSelectOpts.FormatFunc(v, i)
	}

	multiSelectState.Label = multiSelectOpts.Label
	multiSelectState.Options = displayVals
	multiSelectState.Placeholder = multiSelectOpts.Placeholder
	multiSelectState.DefaultValue = defaultVal
	multiSelectState.Required = multiSelectOpts.Required
	multiSelectState.Disabled = multiSelectOpts.Disabled
	sess.State.Set(widgetID, multiSelectState)

	multiSelectProto := convertStateToMultiSelectProto(multiSelectState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_MultiSelect{
				MultiSelect: multiSelectProto,
			},
		},
	})

	cursor.next()

	var value *multiselect.Value
	if multiSelectState.Value != nil {
		value = &multiselect.Value{
			Values:  make([]string, len(multiSelectState.Value)),
			Indexes: make([]int, len(multiSelectState.Value)),
		}
		for i, idx := range multiSelectState.Value {
			value.Values[i] = multiSelectOpts.Options[idx]
			value.Indexes[i] = int(idx)
		}
	}

	return value
}

func (b *uiBuilder) generateMultiSelectID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeMultiSelect.String()+"-"+label+"-"+path.String())
}

func convertStateToMultiSelectProto(state *state.MultiSelectState) *widgetv1.MultiSelect {
	if state == nil {
		return nil
	}
	return &widgetv1.MultiSelect{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertMultiSelectProtoToState(id uuid.UUID, data *widgetv1.MultiSelect) *state.MultiSelectState {
	if data == nil {
		return nil
	}
	return &state.MultiSelectState{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
