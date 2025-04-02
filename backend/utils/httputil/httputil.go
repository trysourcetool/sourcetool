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
	"strings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/net/html"

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
			zap.Stack("stack_trace"),
			zap.String("email", email),
			zap.String("cause", "application"),
		)

		WriteJSON(
			w,
			http.StatusInternalServerError,
			errdefs.ErrInternal(err),
		)
		return
	}

	fields := []zap.Field{
		zap.String("email", email),
		zap.String("frames", v.Frames[0].String()),
		zap.Stack("stack_trace"),
	}

	switch {
	case v.Status >= 500:
		fields = append(fields, zap.String("cause", "application"))
		logger.Logger.Error(err.Error(), fields...)
	case v.Status >= 402, v.Status == 400:
		fields = append(fields, zap.String("cause", "user"))
		logger.Logger.Error(err.Error(), fields...)
	}

	WriteJSON(w, v.Status, v)
}

func ValidateRequest(p any) error {
	v := validator.New()
	v.RegisterValidation("password", validatePassword)

	if err := v.Struct(p); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one letter
	hasLetter := false
	for _, c := range password {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return false
	}

	// Check for at least one digit
	hasDigit := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Check for valid characters only
	validChars := regexp.MustCompile(`^[a-zA-Z0-9!?_+*'"\` + "`" + `#$%&\-^\\@;:,./=~|[\](){}<>]+$`)
	return validChars.MatchString(password)
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
