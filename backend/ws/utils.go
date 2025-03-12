package ws

import (
	"context"
	"runtime"

	"github.com/blendle/zapdriver"
	"github.com/gorilla/websocket"
	exceptionv1 "github.com/trysourcetool/sourcetool/proto/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/logger"
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
	currentUser := ctxutils.CurrentUser(ctx)
	var email string
	if currentUser != nil {
		email = currentUser.Email
	}

	e, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "application")),
		)

		v := errdefs.ErrInternal(err)
		e, _ = v.(*errdefs.Error)
	}

	switch {
	case e.Status >= 500:
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "application")),
			zapdriver.Labels(zapdriver.Label("frames", e.Frames[0].String())),
		)
	case e.Status >= 402, e.Status == 400:
		logger.Logger.Error(
			err.Error(),
			zapdriver.ErrorReport(runtime.Caller(0)),
			zapdriver.Labels(zapdriver.Label("email", email)),
			zapdriver.Labels(zapdriver.Label("cause", "user")),
			zapdriver.Labels(zapdriver.Label("frames", e.Frames[0].String())),
		)
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
		logger.Logger.Sugar().Errorf("Failed to marshal error message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.Logger.Sugar().Errorf("Failed to write error message: %v", err)
		return
	}
}
