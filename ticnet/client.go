package ticnet

import (
	"fmt"
	"github.com/kanyuanzhi/tialloy-client/ticface"
	"github.com/kanyuanzhi/tialloy-client/utils"
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
		Name:       utils.GlobalObject.Name,
		ServerHost: utils.GlobalObject.ServerHost,
		ServerPort: utils.GlobalObject.ServerPort,
		MsgHandler: NewMsgHandler(),
	}
}

func (c *Client) Dial() net.Conn {
	ticker := time.NewTicker(utils.GlobalObject.ReconnectionInterval)
	var err error
	var conn net.Conn
	for {
		select {
		case <-ticker.C:
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort))
			if err != nil {
				utils.GlobalLog.Errorf("touch server %s:%d failed", c.ServerHost, c.ServerPort)
				break
			}
			return conn
		}
	}
}

func (c *Client) Start() {
	utils.GlobalLog.Info("client starts")
	go func() {
		conn := c.Dial()
		utils.GlobalLog.Info("touch server successfully")
		dealConn := NewConnection(c, conn, c.MsgHandler)
		go dealConn.Start()
	}()
}

func (c *Client) Stop() {
	panic("implement me")
}

func (c *Client) Serve() {
	c.Start()

	select{}
}

func (c *Client) SetOnConnStart(hookFunc func(connection ticface.IConnection)) {
	c.OnConnStart = hookFunc
}

func (c *Client) CallOnConnStart(connection ticface.IConnection) {
	if c.OnConnStart != nil {
		utils.GlobalLog.Tracef("call DoOnConnStartHook")
		c.OnConnStart(connection)
	}else {
		utils.GlobalLog.Tracef("there is no DoOnConnStartHook")
	}
}