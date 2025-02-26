package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeMarkdown WidgetType = "markdown"

type MarkdownState struct {
	ID   uuid.UUID
	Body string
}

func (s *MarkdownState) IsWidgetState()      {}
func (s *MarkdownState) GetType() WidgetType { return WidgetTypeMarkdown }
