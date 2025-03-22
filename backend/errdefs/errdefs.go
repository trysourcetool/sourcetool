package errdefs

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	ErrInternal                           = Status("internal_server_error", 500)
	ErrDatabase                           = Status("database_error", 500)
	ErrPermissionDenied                   = Status("permission_denied", 403)
	ErrInvalidArgument                    = Status("invalid_argument", 400)
	ErrAlreadyExists                      = Status("already_exists", 409)
	ErrUnauthenticated                    = Status("unauthenticated", 401)
	ErrAPIKeyNotFound                     = Status("api_key_not_found", 404)
	ErrEnvironmentNotFound                = Status("environment_not_found", 404)
	ErrEnvironmentSlugAlreadyExists       = Status("environment_slug_already_exists", 409)
	ErrEnvironmentDeletionNotAllowed      = Status("environment_deletion_not_allowed", 409)
	ErrGroupNotFound                      = Status("group_not_found", 404)
	ErrGroupSlugAlreadyExists             = Status("group_slug_already_exists", 409)
	ErrHostInstanceNotFound               = Status("host_instance_not_found", 404)
	ErrHostInstanceStatusNotOnline        = Status("host_instance_status_not_online", 404)
	ErrHostConnectionNotFound             = Status("host_connection_not_found", 404)
	ErrHostConnectionClosed               = Status("host_connection_closed", 502)
	ErrOrganizationNotFound               = Status("organization_not_found", 404)
	ErrOrganizationSubdomainAlreadyExists = Status("organization_subdomain_already_exists", 409)
	ErrPageNotFound                       = Status("page_not_found", 404)
	ErrSessionNotFound                    = Status("session_not_found", 404)
	ErrUserNotFound                       = Status("user_not_found", 404)
	ErrUserEmailAlreadyExists             = Status("user_email_already_exists", 409)
	ErrUserRegistrationRequestNotFound    = Status("user_registration_request_not_found", 404)
	ErrUserInvitationNotFound             = Status("user_invitation_not_found", 404)
	ErrUserGoogleAuthRequestNotFound      = Status("user_google_auth_request_not_found", 404)
	ErrUserOrganizationAccessNotFound     = Status("user_organization_access_not_found", 404)
	ErrUserGroupNotFound                  = Status("user_group_not_found", 404)
	ErrUserMultipleOrganizations          = Status("user_multiple_organizations", 422)
)

type Meta []any

type Error struct {
	ID     string         `json:"id"`
	Status int            `json:"status"`
	Title  string         `json:"title"`
	Detail string         `json:"detail"`
	Meta   map[string]any `json:"meta"`
	Frames stackTrace     `json:"-"`
}

type StatusFunc func(error, ...any) error

func Status(title string, status int) StatusFunc {
	return func(err error, vals ...any) error {
		e := &Error{
			ID:     errID(),
			Status: status,
			Title:  title,
			Detail: err.Error(),
			Meta:   make(map[string]any),
			Frames: newFrame(callers()),
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

func (e *Error) Error() string {
	if e.Detail == "" {
		return e.Title
	}

	return e.Detail
}

func (e *Error) StackTrace() []string {
	if len(e.Frames) == 0 {
		return nil
	}
	var stack []string
	for _, frame := range e.Frames {
		stack = append(stack, frame.String())
	}
	return stack
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

func errID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)

	return base64.StdEncoding.EncodeToString(b)
}

func IsEnvironmentNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "environment_not_found"
}

func IsGroupNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "group_not_found"
}

func IsOrganizationNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "organization_not_found"
}

func IsPageNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "page_not_found"
}

func IsSessionNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "session_not_found"
}

func IsHostInstanceNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "host_instance_not_found"
}

func IsUserNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_not_found"
}

func IsUserRegistrationRequestNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_registration_request_not_found"
}

func IsUserInvitationNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_invitation_not_found"
}

func IsUserOrganizationAccessNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_organization_access_not_found"
}

func IsUserGoogleAuthRequestNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_google_auth_request_not_found"
}
