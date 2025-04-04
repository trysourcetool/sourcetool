package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/textarea"
)

func (b *uiBuilder) TextArea(label string, opts ...textarea.Option) string {
	defaultMinLines := int32(2)
	textAreaOpts := &options.TextAreaOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		MaxLength:    nil,
		MinLength:    nil,
		MinLines:     &defaultMinLines,
		MaxLines:     nil,
		AutoResize:   true,
	}

	for _, o := range opts {
		o.Apply(textAreaOpts)
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

	widgetID := b.generatePageID(state.WidgetTypeTextArea, path)
	textAreaState := sess.State.GetTextArea(widgetID)
	if textAreaState == nil {
		textAreaState = &state.TextAreaState{
			ID:    widgetID,
			Value: textAreaOpts.DefaultValue,
		}
	}
	textAreaState.Label = textAreaOpts.Label
	textAreaState.Placeholder = textAreaOpts.Placeholder
	textAreaState.DefaultValue = textAreaOpts.DefaultValue
	textAreaState.Required = textAreaOpts.Required
	textAreaState.Disabled = textAreaOpts.Disabled
	textAreaState.MaxLength = textAreaOpts.MaxLength
	textAreaState.MinLength = textAreaOpts.MinLength
	textAreaState.MaxLines = textAreaOpts.MaxLines
	textAreaState.MinLines = textAreaOpts.MinLines
	textAreaState.AutoResize = textAreaOpts.AutoResize
	sess.State.Set(widgetID, textAreaState)

	textAreaProto := convertStateToTextAreaProto(textAreaState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_TextArea{
				TextArea: textAreaProto,
			},
		},
	})

	cursor.next()

	return conv.SafeValue(textAreaState.Value)
}

func convertStateToTextAreaProto(state *state.TextAreaState) *widgetv1.TextArea {
	if state == nil {
		return nil
	}
	return &widgetv1.TextArea{
		Value:        state.Value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
		MaxLength:    state.MaxLength,
		MinLength:    state.MinLength,
		MaxLines:     state.MaxLines,
		MinLines:     state.MinLines,
		AutoResize:   state.AutoResize,
	}
}

func convertTextAreaProtoToState(id uuid.UUID, data *widgetv1.TextArea) *state.TextAreaState {
	if data == nil {
		return nil
	}
	return &state.TextAreaState{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
		MaxLength:    data.MaxLength,
		MinLength:    data.MinLength,
		MaxLines:     data.MaxLines,
		MinLines:     data.MinLines,
		AutoResize:   data.AutoResize,
	}
}
