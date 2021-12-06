package ticface

import "net"

type IClient interface {
	Start()
	Stop()
	Serve()
	Dial() net.Conn

	SetOnConnStart(func(connection IConnection))
	CallOnConnStart(connection IConnection)
}
