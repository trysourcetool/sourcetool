package markdown

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.MarkdownOptions)
}
