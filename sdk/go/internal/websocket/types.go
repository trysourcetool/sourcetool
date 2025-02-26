package websocket

import (
	"fmt"

	exceptionv1 "github.com/trysourcetool/sourcetool-proto/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	"google.golang.org/protobuf/proto"
)

type MessageHandlerFunc func(*websocketv1.Message) error

func unmarshalMessage(data []byte) (*websocketv1.Message, error) {
	var msg websocketv1.Message
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func marshalMessage(msg *websocketv1.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func NewMessage(id string, payload proto.Message) (*websocketv1.Message, error) {
	msg := &websocketv1.Message{
		Id: id,
	}

	switch p := payload.(type) {
	case *websocketv1.InitializeHost:
		msg.Type = &websocketv1.Message_InitializeHost{InitializeHost: p}
	case *websocketv1.InitializeClient:
		msg.Type = &websocketv1.Message_InitializeClient{InitializeClient: p}
	case *websocketv1.RenderWidget:
		msg.Type = &websocketv1.Message_RenderWidget{RenderWidget: p}
	case *websocketv1.RerunPage:
		msg.Type = &websocketv1.Message_RerunPage{RerunPage: p}
	case *websocketv1.CloseSession:
		msg.Type = &websocketv1.Message_CloseSession{CloseSession: p}
	case *websocketv1.ScriptFinished:
		msg.Type = &websocketv1.Message_ScriptFinished{ScriptFinished: p}
	case *exceptionv1.Exception:
		msg.Type = &websocketv1.Message_Exception{Exception: p}
	default:
		return nil, fmt.Errorf("unsupported message type: %T", payload)
	}

	return msg, nil
}
