package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var ErrTcpServerClosed error = errors.New("tcp: Server closed")

type Server struct {
	connections map[int]*Connection
	mutex       sync.Mutex
	listen      *net.TCPListener
	isClose     bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[int]*Connection),
		isClose:     false,
	}
}

func (s *Server) ListenAndServer() error {
	host := "0.0.0.0"
	tcpPort := 12345
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, tcpPort))
	if err != nil {
		return err
	}
	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return err
	}
	log.Printf("tcp server listening on %s\n", addr)
	s.listen = listen
	return s.Serve(listen)
}

func (s *Server) Serve(l *net.TCPListener) error {
	var connId int = 0
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		c := NewConnection(conn, connId)
		s.AddConn(c)
		connId++

		go c.Serve()
	}
}

func (s *Server) AddConn(conn *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[conn.ConnId] = conn
}

func (s *Server) RemoveConn(conn *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connections, conn.ConnId)
}

func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Lock()
	if s.isClose {
		return nil
	}
	if err := s.listen.Close(); err != nil {
		return err
	}

	for connId, conn := range s.connections {
		delete(s.connections, connId)
		conn.Close()
	}
	s.isClose = true
	log.Println("tcp server stop.")
	return nil
}
