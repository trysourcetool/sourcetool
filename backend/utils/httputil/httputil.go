package httputil

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/blendle/zapdriver"
	"github.com/go-playground/validator/v10"
	"golang.org/x/net/html"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func WriteBytes(w http.ResponseWriter, status int, b []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func WriteErrJSON(ctx context.Context, w http.ResponseWriter, err error) {
	currentUser := ctxutil.CurrentUser(ctx)
	var email string
	if currentUser != nil {
		email = currentUser.Email
	}

	v, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "application")),
		)

		WriteJSON(
			w,
			http.StatusInternalServerError,
			errdefs.ErrInternal(err),
		)
		return
	}

	switch {
	case v.Status >= 500:
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "application")),
			zapdriver.Labels(zapdriver.Label("frames", v.Frames[0].String())),
		)
	case v.Status >= 402, v.Status == 400:
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "user")),
			zapdriver.Labels(zapdriver.Label("frames", v.Frames[0].String())),
		)
	}

	WriteJSON(w, v.Status, v)
}

func ValidateRequest(p any) error {
	v := validator.New()
	v.RegisterValidation("password", validatePassword)

	if err := v.Struct(p); err != nil {
		return err
	}

	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	r := regexp.MustCompile(`^[!-~]{8,32}$`)
	return r.MatchString(fl.Field().String())
}

func ValidateJSONString(s string) error {
	var in any
	return json.Unmarshal([]byte(s), &in)
}

func ValidateHTMLString(s string) error {
	if _, err := html.Parse(strings.NewReader(s)); err != nil {
		return err
	}
	return nil
}

func ValidateXMLString(s string) error {
	var in any
	return xml.Unmarshal([]byte(s), &in)
}

func GetIP(r *http.Request) (string, error) {
	addr := r.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = r.RemoteAddr
	}

	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", addr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", addr)
	}

	return userIP.String(), nil
}

func HTTPScheme() string {
	if config.Config.Env == config.EnvLocal {
		return "http"
	}
	return "https"
}

func GetSubdomainFromHost(host string) (string, error) {
	if host == "" {
		return "", errors.New("empty host")
	}
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid host format")
	}
	return parts[0], nil
}
