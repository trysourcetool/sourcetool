package errdefs

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrInternal         = Exception("internal_server_error")
	ErrInvalidParameter = Exception("invalid_parameter")
	ErrSessionNotFound  = Exception("session_not_found")
	ErrPageNotFound     = Exception("page_not_found")
	ErrRunPage          = Exception("run_page_error")
)

type Meta []any

type Error struct {
	Title   string
	Message string
	Meta    map[string]any
	Frames  stackTrace
}

type ExceptionFunc func(error, ...any) error

func Exception(title string) ExceptionFunc {
	return func(err error, vals ...any) error {
		e := &Error{
			Title:   title,
			Message: err.Error(),
			Frames:  newFrame(callers()),
		}

		for _, any := range vals {
			switch any := any.(type) {
			case Meta:
				e.Meta = appendMeta(e.Meta, any...)
			}
		}

		x, ok := err.(*Error)
		if ok {
			e.Frames = x.Frames
		}

		return e
	}
}

func appendMeta(meta map[string]any, keyvals ...any) map[string]any {
	if meta == nil {
		meta = make(map[string]any)
	}
	var k string
	for n, v := range keyvals {
		if n%2 == 0 {
			k = fmt.Sprint(v)
		} else {
			meta[k] = v
		}
	}
	return meta
}

func (e *Error) Error() string {
	if e.Message == "" {
		return e.Title
	}

	return e.Message
}

func (e *Error) StackTrace() []string {
	if len(e.Frames) == 0 {
		return []string{}
	}
	strs := make([]string, len(e.Frames))
	for i, f := range e.Frames {
		strs[i] = f.String()
	}
	return strs
}

type frame struct {
	file           string
	lineNumber     int
	name           string
	programCounter uintptr
}

type stackTrace []*frame

func newFrame(pcs []uintptr) stackTrace {
	frames := []*frame{}

	for _, pc := range pcs {
		frame := &frame{programCounter: pc}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			return frames
		}
		frame.name = trimPkgName(fn)

		frame.file, frame.lineNumber = fn.FileLine(pc - 1)
		frames = append(frames, frame)
	}

	return frames
}

func (f *frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.file, f.lineNumber, f.name)
}

func trimPkgName(fn *runtime.Func) string {
	name := fn.Name()
	if ld := strings.LastIndex(name, "."); ld >= 0 {
		name = name[ld+1:]
	}

	return name
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])

	return pcs[0 : n-2]
}
