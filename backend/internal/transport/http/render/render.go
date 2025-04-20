package render

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

func JSON(w http.ResponseWriter, status int, v any) error {
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

func Bytes(w http.ResponseWriter, status int, b []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func Error(ctx context.Context, w http.ResponseWriter, err error) {
	currentUser := internal.CurrentUser(ctx)
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

		JSON(
			w,
			http.StatusInternalServerError,
			errdefs.ErrInternal(err),
		)
		return
	}

	fields := []zap.Field{
		zap.String("email", email),
		zap.String("error_stacktrace", strings.Join(v.StackTrace(), "\n")),
	}

	switch {
	case v.Status >= 500:
		fields = append(fields, zap.String("cause", "application"))
		logger.Logger.Error(err.Error(), fields...)
	case v.Status >= 402, v.Status == 400:
		fields = append(fields, zap.String("cause", "user"))
		logger.Logger.Error(err.Error(), fields...)
	default:
		fields = append(fields, zap.String("cause", "internal_info"))
		logger.Logger.Warn(err.Error(), fields...)
	}

	JSON(w, v.Status, v)
}
