package ticnet

import (
	"context"
	"errors"
	"fmt"
	"github.com/kanyuanzhi/tialloy-client/global"
	"github.com/kanyuanzhi/tialloy-client/ticface"
	"github.com/kanyuanzhi/tialloy-client/ticlog"
	"io"
	"net"
	"sync"
)

type Connection struct {
	Client ticface.IClient
	Conn   net.Conn

	sync.RWMutex

	MsgChan    chan []byte
	MsgHandler ticface.IMsgHandler
	ctx        context.Context
	cancel     context.CancelFunc
	IsClosed   bool
}

func NewConnection(client ticface.IClient, conn net.Conn, handler ticface.IMsgHandler) ticface.IConnection {
	return &Connection{
		Client:     client,
		Conn:       conn,
		IsClosed:   false,
		MsgChan:    make(chan []byte),
		MsgHandler: handler,
	}
}

func (c *Connection) StartReader() {
	// TODO unfinished
	ticlog.Log.Info("tcp reader goroutine is running")
	defer c.Reconnect()
	defer ticlog.Log.Warn("tcp reader goroutine exited")
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			dp := NewDataPack()

			dataHeadBuf := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.GetConn(), dataHeadBuf); err != nil {
				ticlog.Log.Error(err)
				return
			}

			message, err := dp.Unpack(dataHeadBuf)
			if err != nil {
				ticlog.Log.Error(err)
				return
			}

			var dataBuf []byte
			if message.GetDataLen() > 0 {
				dataBuf = make([]byte, message.GetDataLen())
				if _, err := io.ReadFull(c.GetConn(), dataBuf); err != nil {
					ticlog.Log.Error(err)
					return
				}
			}

			message.SetData(dataBuf)
			request := NewRequest(c, message)

			if global.Object.TcpWorkerPoolSize > 0 {
				go c.MsgHandler.SendMsgToTaskQueue(request)
			} else {
				go c.MsgHandler.DoMsgHandler(request)
			}
		}

	}
}

func (c *Connection) StartWriter() {
	ticlog.Log.Info("tcp writer goroutine is running")
	defer ticlog.Log.Warn("tcp writer goroutine exited")
	for {
		select {
		case msg := <-c.MsgChan:
			if _, err := c.Conn.Write(msg); err != nil {
				ticlog.Log.Error(err)
				break
			}
		case <-c.ctx.Done():
			ticlog.Log.Trace("exit")
			return
		}
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.StartReader()
	go c.StartWriter()

	c.Client.CallOnConnStart(c)
}

func (c *Connection) Reconnect() {
	c.cancel()
	c.Conn = c.Client.Dial()
	go c.Start()
}

func (c *Connection) Stop() {
	c.Lock()
	defer c.Unlock()

	if c.IsClosed == true {
		return
	}

	c.Conn.Close()
	c.cancel()

	close(c.MsgChan)

	c.IsClosed = true
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()
	binaryMessage, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack tcp msgID=%d err", msgID))
	}
	c.MsgChan <- binaryMessage
	return nil
}

func (c *Connection) GetConn() net.Conn {
	return c.Conn
}

func (c *Connection) Context() context.Context {
	return c.ctx
}
