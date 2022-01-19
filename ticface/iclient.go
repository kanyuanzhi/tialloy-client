package ticface

import "net"

type IClient interface {
	Start()
	Stop()
	Serve()
	Dial() net.Conn

	AddRouter(msgID uint32, router IRouter)

	SetOnConnStart(func(connection IConnection))
	CallOnConnStart(connection IConnection)
}
