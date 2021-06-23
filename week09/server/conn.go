package server

import (
	"context"
	"log"
	"net"
)

type Connection struct {
	conn       *net.TCPConn
	ConnId     int
	ctx        context.Context
	sendMsgBuf chan *Message
}

func NewConnection(conn *net.TCPConn, connId int) *Connection {
	return &Connection{
		conn:       conn,
		ConnId:     connId,
		sendMsgBuf: make(chan *Message, 512),
	}
}

func (c *Connection) readHandle() error {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:

		}
	}
}

func (c *Connection) writeHandle() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.sendMsgBuf:
			if _, err := c.conn.Write(msg.Data); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Connection) SendMessage(msgId int, msg []byte) error {
	// c.sendMsgBuf <- msg

}

func (c *Connection) Serve() {

}

func (c *Connection) Close() {
	c.conn.Close()
}
