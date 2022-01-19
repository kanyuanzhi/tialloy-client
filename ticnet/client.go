package ticnet

import (
	"fmt"
	"github.com/kanyuanzhi/tialloy-client/global"
	"github.com/kanyuanzhi/tialloy-client/ticface"
	"github.com/kanyuanzhi/tialloy-client/ticlog"
	"net"
	"time"
)

type Client struct {
	Name string

	ServerHost string
	ServerPort int

	MsgHandler ticface.IMsgHandler

	OnConnStart func(connection ticface.IConnection)
}

func NewClient() ticface.IClient {
	return &Client{
		Name:       global.Object.Name,
		ServerHost: global.Object.ServerHost,
		ServerPort: global.Object.ServerPort,
		MsgHandler: NewMsgHandler(),
	}
}

func (c *Client) Dial() net.Conn {
	ticker := time.NewTicker(global.Object.ReconnectInterval * time.Second)
	var err error
	var conn net.Conn
	for {
		select {
		case <-ticker.C:
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort))
			if err != nil {
				ticlog.Log.Errorf("touch server %s:%d failed, trying retouch every %d second(s)", c.ServerHost, c.ServerPort, global.Object.ReconnectInterval)
				break
			}
			return conn
		}
	}
}

func (c *Client) Start() {
	ticlog.Log.Info("client is starting")
	go func() {
		conn := c.Dial()
		ticlog.Log.Infof("touch server %s:%d successfully", c.ServerHost, c.ServerPort)
		dealConn := NewConnection(c, conn, c.MsgHandler)
		go dealConn.Start()
	}()
}

func (c *Client) Stop() {
	panic("implement me")
}

func (c *Client) Serve() {
	c.Start()

	select {}
}

func (c *Client) AddRouter(msgID uint32, router ticface.IRouter) {
	c.MsgHandler.AddRouter(msgID, router)
}

func (c *Client) SetOnConnStart(hookFunc func(connection ticface.IConnection)) {
	c.OnConnStart = hookFunc
}

func (c *Client) CallOnConnStart(connection ticface.IConnection) {
	if c.OnConnStart != nil {
		ticlog.Log.Tracef("call DoOnConnStartHook")
		c.OnConnStart(connection)
	} else {
		ticlog.Log.Tracef("there is no DoOnConnStartHook")
	}
}
