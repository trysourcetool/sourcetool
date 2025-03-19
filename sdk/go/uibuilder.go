package sourcetool

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/checkbox"
	"github.com/trysourcetool/sourcetool-go/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/dateinput"
	"github.com/trysourcetool/sourcetool-go/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/multiselect"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/radio"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textarea"
	"github.com/trysourcetool/sourcetool-go/textinput"
	"github.com/trysourcetool/sourcetool-go/timeinput"
)

type UIBuilder interface {
	Context() context.Context
	Markdown(string)
	TextInput(string, ...textinput.Option) string
	NumberInput(string, ...numberinput.Option) *float64
	DateInput(string, ...dateinput.Option) *time.Time
	DateTimeInput(string, ...datetimeinput.Option) *time.Time
	TimeInput(string, ...timeinput.Option) *time.Time
	Selectbox(string, ...selectbox.Option) *selectbox.Value
	MultiSelect(string, ...multiselect.Option) *multiselect.Value
	Radio(string, ...radio.Option) *radio.Value
	Checkbox(string, ...checkbox.Option) bool
	CheckboxGroup(string, ...checkboxgroup.Option) *checkboxgroup.Value
	TextArea(string, ...textarea.Option) string
	Table(any, ...table.Option) table.Value
	Button(string, ...button.Option) bool
	Form(string, ...form.Option) (UIBuilder, bool)
	Columns(int, ...columns.Option) []UIBuilder
}

type uiBuilder struct {
	runtime *runtime
	context context.Context
	cursor  *cursor
	session *session.Session
	page    *page
}

func (b *uiBuilder) Context() context.Context {
	return b.context
}

func (b *uiBuilder) generatePageID(widgetType state.WidgetType, path []int) uuid.UUID {
	if b.page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, v := range path {
		strPath[i] = strconv.Itoa(v)
	}
	return uuid.NewV5(b.page.id, widgetType.String()+"-"+strings.Join(strPath, "_"))
}

type path []int

func (p path) String() string {
	strPath := make([]string, len(p))
	for i, num := range p {
		strPath[i] = fmt.Sprint(num)
	}
	return strings.Join(strPath, "")
}

type cursor struct {
	parentPath []int
	index      int
}

func newCursor() *cursor {
	return &cursor{
		parentPath: []int{},
		index:      0,
	}
}

func (c *cursor) getPath() path {
	path := make([]int, len(c.parentPath))
	copy(path, c.parentPath)
	path = append(path, c.index)
	return path
}

func (c *cursor) next() {
	c.index++
}
