package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"week09/conf"
	"week09/util"
)

type Connection struct {
	conn       *net.TCPConn
	connId     uint64
	s          *Server
	mutex      sync.RWMutex
	isClose    bool
	sendMsgBuf chan *util.Message
	recvMsgBuf chan *util.Message
}

func NewConnection(conn *net.TCPConn, connId uint64, s *Server) *Connection {
	c := &Connection{
		conn:       conn,
		connId:     connId,
		s:          s,
		sendMsgBuf: make(chan *util.Message, conf.TCPConfig.SendBufferSize),
		recvMsgBuf: make(chan *util.Message, conf.TCPConfig.RecvBufferSize),
	}
	go c.MsgHandle()
	go c.WriteHandle()
	return c
}

func (c *Connection) MsgHandle() {
	for msg := range c.recvMsgBuf {
		// 业务操作判断
		switch msg.GetOperator() {
		case 1:
			sendMsg := util.NewMessage(1, 1, 1, []byte("pong..."))
			if err := c.MessageToSendChan(sendMsg); err != nil {
				log.Fatalln(err)
			}
		case 2:
			sendMsg := util.NewMessage(1, 2, 2, []byte("Golang..."))
			if err := c.MessageToSendChan(sendMsg); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (c *Connection) WriteHandle() {
	for msg := range c.sendMsgBuf {
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

func (c *Connection) MessageToSendChan(msg *util.Message) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClose {
		return errors.New("connection closed when send reply")
	}
	c.sendMsgBuf <- msg
	return nil
}

func (c *Connection) MessageToRecvChan(msg *util.Message) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClose {
		return errors.New("connection closed when send reply")
	}
	c.recvMsgBuf <- msg
	return nil
}

func (c *Connection) Serve() {
	defer c.Close()
	for {
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
		if err := c.MessageToRecvChan(msg); err != nil {
			return
		}
	}
}

func (c *Connection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClose {
		return nil
	}
	log.Println("connection closed, connID = ", c.connId)

	err := c.conn.Close()
	c.s.RemoveConn(c)
	close(c.sendMsgBuf)
	close(c.recvMsgBuf)
	c.isClose = true
	return err
}
