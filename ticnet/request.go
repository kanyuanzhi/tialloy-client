package ticnet

import (
	"github.com/kanyuanzhi/tialloy-client/ticface"
)

type Request struct {
	conn    ticface.IConnection
	message ticface.IMessage
}

func (r *Request) GetConnection() ticface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.message.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgID()
}

func NewRequest(conn ticface.IConnection, message ticface.IMessage) ticface.IRequest {
	return &Request{
		conn:    conn,
		message: message,
	}
}
