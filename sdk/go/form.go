package sourcetool

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/form"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func (b *uiBuilder) Form(buttonLabel string, opts ...form.Option) (UIBuilder, bool) {
	formOpts := &options.FormOptions{
		ButtonLabel:    buttonLabel,
		ButtonDisabled: false,
		ClearOnSubmit:  false,
	}

	for _, o := range opts {
		o.Apply(formOpts)
	}

	sess := b.session
	if sess == nil {
		return b, false
	}
	page := b.page
	if page == nil {
		return b, false
	}
	cursor := b.cursor
	if cursor == nil {
		return b, false
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeForm, path)
	formState := sess.State.GetForm(widgetID)
	if formState == nil {
		formState = &state.FormState{
			ID:    widgetID,
			Value: false,
		}
	}
	formState.ButtonLabel = formOpts.ButtonLabel
	formState.ButtonDisabled = formOpts.ButtonDisabled
	formState.ClearOnSubmit = formOpts.ClearOnSubmit
	sess.State.Set(widgetID, formState)

	form := convertStateToFormProto(formState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Form{
				Form: form,
			},
		},
	})

	cursor.next()

	childCursor := newCursor()
	childCursor.parentPath = path

	childBuilder := &uiBuilder{
		runtime: b.runtime,
		session: sess,
		page:    page,
		cursor:  childCursor,
	}

	return childBuilder, formState.Value
}

func convertStateToFormProto(state *state.FormState) *widgetv1.Form {
	if state == nil {
		return nil
	}
	return &widgetv1.Form{
		Value:          state.Value,
		ButtonLabel:    state.ButtonLabel,
		ButtonDisabled: state.ButtonDisabled,
		ClearOnSubmit:  state.ClearOnSubmit,
	}
}

func convertFormProtoToState(id uuid.UUID, data *widgetv1.Form) *state.FormState {
	if data == nil {
		return nil
	}
	return &state.FormState{
		ID:             id,
		Value:          data.Value,
		ButtonLabel:    data.ButtonLabel,
		ButtonDisabled: data.ButtonDisabled,
		ClearOnSubmit:  data.ClearOnSubmit,
	}
}
