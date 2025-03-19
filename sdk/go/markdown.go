package sourcetool

import (
	"github.com/gofrs/uuid/v5"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func (b *uiBuilder) Markdown(body string) {
	markdownOpts := &options.MarkdownOptions{
		Body: body,
	}

	sess := b.session
	if sess == nil {
		return
	}
	page := b.page
	if page == nil {
		return
	}
	cursor := b.cursor
	if cursor == nil {
		return
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeMarkdown, path)
	markdownState := sess.State.GetMarkdown(widgetID)
	if markdownState == nil {
		markdownState = &state.MarkdownState{
			ID: widgetID,
		}
	}
	markdownState.Body = markdownOpts.Body
	sess.State.Set(widgetID, markdownState)

	markdown := convertStateToMarkdownProto(markdownState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Markdown{
				Markdown: markdown,
			},
		},
	})

	cursor.next()
}

func convertStateToMarkdownProto(state *state.MarkdownState) *widgetv1.Markdown {
	return &widgetv1.Markdown{
		Body: state.Body,
	}
}

func convertMarkdownProtoToState(id uuid.UUID, data *widgetv1.Markdown) *state.MarkdownState {
	return &state.MarkdownState{
		ID:   id,
		Body: data.Body,
	}
}
