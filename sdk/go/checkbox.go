package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkbox"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func (b *uiBuilder) Checkbox(label string, opts ...checkbox.Option) bool {
	checkboxOpts := &options.CheckboxOptions{
		Label:        label,
		DefaultValue: false,
		Required:     false,
		Disabled:     false,
	}

	for _, o := range opts {
		o.Apply(checkboxOpts)
	}

	sess := b.session
	if sess == nil {
		return false
	}
	page := b.page
	if page == nil {
		return false
	}
	cursor := b.cursor
	if cursor == nil {
		return false
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeCheckbox, path)
	checkboxState := sess.State.GetCheckbox(widgetID)
	if checkboxState == nil {
		checkboxState = &state.CheckboxState{
			ID:    widgetID,
			Value: checkboxOpts.DefaultValue,
		}
	}
	checkboxState.Label = checkboxOpts.Label
	checkboxState.DefaultValue = checkboxOpts.DefaultValue
	checkboxState.Required = checkboxOpts.Required
	checkboxState.Disabled = checkboxOpts.Disabled
	sess.State.Set(widgetID, checkboxState)

	checkboxProto := convertStateToCheckboxProto(checkboxState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Checkbox{
				Checkbox: checkboxProto,
			},
		},
	})

	cursor.next()

	return checkboxState.Value
}

func convertStateToCheckboxProto(state *state.CheckboxState) *widgetv1.Checkbox {
	if state == nil {
		return nil
	}
	return &widgetv1.Checkbox{
		Value:        state.Value,
		Label:        state.Label,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertCheckboxProtoToState(id uuid.UUID, data *widgetv1.Checkbox) *state.CheckboxState {
	if data == nil {
		return nil
	}
	return &state.CheckboxState{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
