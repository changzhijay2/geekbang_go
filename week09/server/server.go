package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"week09/conf"
)

type Server struct {
	connections map[uint64]*Connection
	mutex       sync.Mutex
	listen      *net.TCPListener
}

func NewServer() *Server {
	return &Server{
		connections: make(map[uint64]*Connection),
	}
}

func (s *Server) ListenAndServer() error {
	host := conf.TCPConfig.Host
	port := conf.TCPConfig.Port
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return err
	}
	log.Printf("tcp server listening on %s\n", addr)
	s.mutex.Lock()
	s.listen = listen
	s.mutex.Unlock()
	return s.Serve(listen)
}

func (s *Server) Serve(l *net.TCPListener) error {
	var connId uint64 = 0
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			log.Println("accept failed, err:", err)
			return err
		}
		if len(s.connections) > conf.TCPConfig.MaxConn {
			log.Println("Maximum number of connections reached, ignoring this request, addr:", conn.RemoteAddr().String())
			continue
		}

		c := NewConnection(conn, connId, s)
		s.AddConn(c)
		connId++
		log.Printf("accept connection, addr:%s, connID = %d.\n", conn.RemoteAddr().String(), connId)

		go c.Serve()
	}
}

func (s *Server) AddConn(conn *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[conn.connId] = conn
}

func (s *Server) RemoveConn(conn *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connections, conn.connId)
}

func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.listen.Close(); err != nil {
		return err
	}
	for _, conn := range s.connections {
		conn.Close()
	}
	log.Println("tcp server stop.")
	return nil
}
