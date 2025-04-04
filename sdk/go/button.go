package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/button"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func (b *uiBuilder) Button(label string, opts ...button.Option) bool {
	buttonOpts := &options.ButtonOptions{
		Label:    label,
		Disabled: false,
	}

	for _, o := range opts {
		o.Apply(buttonOpts)
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

	widgetID := b.generatePageID(state.WidgetTypeButton, path)
	buttonState := sess.State.GetButton(widgetID)
	if buttonState == nil {
		buttonState = &state.ButtonState{
			ID:    widgetID,
			Value: false,
		}
	}
	buttonState.Label = buttonOpts.Label
	buttonState.Disabled = buttonOpts.Disabled
	sess.State.Set(widgetID, buttonState)

	button := convertStateToButtonProto(buttonState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Button{
				Button: button,
			},
		},
	})

	cursor.next()

	return buttonState.Value
}

func convertStateToButtonProto(state *state.ButtonState) *widgetv1.Button {
	return &widgetv1.Button{
		Value:    state.Value,
		Label:    state.Label,
		Disabled: state.Disabled,
	}
}

func convertButtonProtoToState(id uuid.UUID, data *widgetv1.Button) *state.ButtonState {
	if data == nil {
		return nil
	}
	return &state.ButtonState{
		ID:       id,
		Value:    data.Value,
		Label:    data.Label,
		Disabled: data.Disabled,
	}
}
