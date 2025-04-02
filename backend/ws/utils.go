package ws

import (
	"context"
	"strings"

	"github.com/gorilla/websocket"
	exceptionv1 "github.com/trysourcetool/sourcetool/proto/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

func SendResponse(conn *websocket.Conn, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}

	return nil
}

func SendErrResponse(ctx context.Context, conn *websocket.Conn, id string, err error) {
	currentUser := ctxutil.CurrentUser(ctx)
	var email string
	if currentUser != nil {
		email = currentUser.Email
	}

	e, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zap.Stack("stack_trace"),
			zap.String("email", email),
			zap.String("cause", "application"),
		)

		v := errdefs.ErrInternal(err)
		e, _ = v.(*errdefs.Error)
	} else {
		fields := []zap.Field{
			zap.String("email", email),
			zap.String("error_stacktrace", strings.Join(e.StackTrace(), "\n")),
		}

		switch {
		case e.Status >= 500:
			fields = append(fields, zap.String("cause", "application"))
			logger.Logger.Error(err.Error(), fields...)
		case e.Status >= 402, e.Status == 400:
			fields = append(fields, zap.String("cause", "user"))
			logger.Logger.Error(err.Error(), fields...)
		default:
			fields = append(fields, zap.String("cause", "internal_info"))
			logger.Logger.Warn(err.Error(), fields...)
		}
	}

	msg := &websocketv1.Message{
		Id: id,
		Type: &websocketv1.Message_Exception{
			Exception: &exceptionv1.Exception{
				Title:      e.Title,
				Message:    e.Detail,
				StackTrace: e.StackTrace(),
			},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Failed to marshal WS error message", zap.Error(err))
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.Logger.Error("Failed to write WS error message", zap.Error(err))
		return
	}
}
