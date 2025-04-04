package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func convertPathToInt32Slice(p path) []int32 {
	result := make([]int32, len(p))
	for i, v := range p {
		result[i] = int32(v)
	}
	return result
}

func (b *uiBuilder) TextInput(label string, opts ...textinput.Option) string {
	textInputOpts := &options.TextInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		MaxLength:    nil,
		MinLength:    nil,
	}

	for _, o := range opts {
		o.Apply(textInputOpts)
	}

	sess := b.session
	if sess == nil {
		return ""
	}
	page := b.page
	if page == nil {
		return ""
	}
	cursor := b.cursor
	if cursor == nil {
		return ""
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeTextInput, path)
	textInputState := sess.State.GetTextInput(widgetID)
	if textInputState == nil {
		textInputState = &state.TextInputState{
			ID:    widgetID,
			Value: textInputOpts.DefaultValue,
		}
	}
	textInputState.Label = textInputOpts.Label
	textInputState.Placeholder = textInputOpts.Placeholder
	textInputState.DefaultValue = textInputOpts.DefaultValue
	textInputState.Required = textInputOpts.Required
	textInputState.Disabled = textInputOpts.Disabled
	textInputState.MaxLength = textInputOpts.MaxLength
	textInputState.MinLength = textInputOpts.MinLength
	sess.State.Set(widgetID, textInputState)

	textInput := convertStateToTextInputProto(textInputState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_TextInput{
				TextInput: textInput,
			},
		},
	})

	cursor.next()

	return conv.SafeValue(textInputState.Value)
}

func convertStateToTextInputProto(state *state.TextInputState) *widgetv1.TextInput {
	if state == nil {
		return nil
	}
	return &widgetv1.TextInput{
		Value:        state.Value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
		MaxLength:    state.MaxLength,
		MinLength:    state.MinLength,
	}
}

func convertTextInputProtoToState(id uuid.UUID, data *widgetv1.TextInput) *state.TextInputState {
	if data == nil {
		return nil
	}
	return &state.TextInputState{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
		MaxLength:    data.MaxLength,
		MinLength:    data.MinLength,
	}
}
