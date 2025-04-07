package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/selectbox"
)

func (b *uiBuilder) Selectbox(label string, opts ...selectbox.Option) *selectbox.Value {
	selectboxOpts := &options.SelectboxOptions{
		Label:        label,
		DefaultValue: nil,
		Placeholder:  "",
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, o := range opts {
		o.Apply(selectboxOpts)
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

	var defaultVal *int32
	if selectboxOpts.DefaultValue != nil {
		for i, o := range selectboxOpts.Options {
			if conv.SafeValue(selectboxOpts.DefaultValue) == o {
				v := int32(i)
				defaultVal = &v
				break
			}
		}
	}

	widgetID := b.generatePageID(state.WidgetTypeSelectbox, path)
	selectboxState := sess.State.GetSelectbox(widgetID)
	if selectboxState == nil {
		selectboxState = &state.SelectboxState{
			ID:           widgetID,
			Value:        defaultVal,
			DefaultValue: defaultVal,
		}
	}

	if selectboxOpts.FormatFunc == nil {
		selectboxOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(selectboxOpts.Options))
	for i, v := range selectboxOpts.Options {
		displayVals[i] = selectboxOpts.FormatFunc(v, i)
	}

	selectboxState.Label = selectboxOpts.Label
	selectboxState.Options = displayVals
	selectboxState.Placeholder = selectboxOpts.Placeholder
	selectboxState.DefaultValue = defaultVal
	selectboxState.Required = selectboxOpts.Required
	selectboxState.Disabled = selectboxOpts.Disabled
	sess.State.Set(widgetID, selectboxState)

	selectboxProto := convertStateToSelectboxProto(selectboxState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Selectbox{
				Selectbox: selectboxProto,
			},
		},
	})

	cursor.next()

	var value *selectbox.Value
	if selectboxState.Value != nil {
		value = &selectbox.Value{
			Value: selectboxOpts.Options[*selectboxState.Value],
			Index: int(conv.SafeValue(selectboxState.Value)),
		}
	}

	return value
}

func convertStateToSelectboxProto(state *state.SelectboxState) *widgetv1.Selectbox {
	if state == nil {
		return nil
	}
	return &widgetv1.Selectbox{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertSelectboxProtoToState(id uuid.UUID, data *widgetv1.Selectbox) *state.SelectboxState {
	if data == nil {
		return nil
	}
	return &state.SelectboxState{
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
