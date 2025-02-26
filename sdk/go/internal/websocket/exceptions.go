package websocket

import (
	"github.com/trysourcetool/sourcetool-go/internal/errdefs"
	exceptionv1 "github.com/trysourcetool/sourcetool-proto/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
)

func (c *client) sendException(id string, err error) {
	e, ok := err.(*errdefs.Error)
	if !ok {
		v := errdefs.ErrInternal(err)
		e = v.(*errdefs.Error)
	}

	exception := &exceptionv1.Exception{
		Title:      e.Title,
		Message:    e.Message,
		StackTrace: e.StackTrace(),
	}

	c.Enqueue(id, &websocketv1.Message{
		Id: id,
		Type: &websocketv1.Message_Exception{
			Exception: exception,
		},
	})
}
