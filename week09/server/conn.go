package server

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"

	"week09/conf"
	"week09/util"
)

type Connection struct {
	conn       *net.TCPConn
	ConnId     uint64
	s          *Server
	ctx        context.Context
	cancel     context.CancelFunc
	sendMsgBuf chan *util.Message
	recvMsgBuf chan *util.Message
}

func NewConnection(conn *net.TCPConn, connId uint64, s *Server) *Connection {
	return &Connection{
		conn:       conn,
		ConnId:     connId,
		s:          s,
		sendMsgBuf: make(chan *util.Message, conf.TCPConfig.SendBufferSize),
		recvMsgBuf: make(chan *util.Message, conf.TCPConfig.RecvBufferSize),
	}
}

func (c *Connection) MsgHandle() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.recvMsgBuf:
			// 业务操作判断
			switch msg.GetOperator() {
			case 1:
				sendMsg := util.NewMessage(1, 1, 1, []byte("pong..."))
				c.SendMessage(sendMsg)
			case 2:
				sendMsg := util.NewMessage(1, 2, 2, []byte("Golang..."))
				c.SendMessage(sendMsg)
			}
		}
	}
}

func (c *Connection) WriteHandle() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.sendMsgBuf:
			b, err := util.Pack(msg)
			if err != nil {
				log.Println(err)
				return
			}
			if _, err := c.conn.Write(b); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Connection) SendMessage(msg *util.Message) {
	c.sendMsgBuf <- msg
}

func (c *Connection) Serve() {
	defer c.Close()
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.MsgHandle()
	go c.WriteHandle()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			r := bufio.NewReader(c.conn)
			msg, err := util.Unpack(r)
			if err != nil {
				if err == io.EOF {
					c.s.RemoveConn(c)
					c.Close()
					return
				}
				log.Println(err)
				return
			}
			c.recvMsgBuf <- msg
		}
	}
}

func (c *Connection) Close() error {
	c.cancel()
	return c.conn.Close()
}
