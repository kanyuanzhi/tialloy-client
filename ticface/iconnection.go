package ticface

import (
	"context"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetConn() net.Conn

	SendMsg(msgID uint32, data []byte) error
	Context() context.Context
}
