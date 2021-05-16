package server

import (
	"fmt"
	"log"
	"net"

	v1 "github.com/changzhijay2/geekbang_go/week04/api/user/v1"
	"github.com/changzhijay2/geekbang_go/week04/internal/biz"
	"github.com/changzhijay2/geekbang_go/week04/pkg/config"

	"google.golang.org/grpc"
)

type Server struct {
	s   *grpc.Server
	biz *biz.Biz
}

func NewServer(biz *biz.Biz) *Server {
	return &Server{
		s:   grpc.NewServer(),
		biz: biz,
	}
}

func (s *Server) RegisterService() {
	v1.RegisterUserServiceServer(s.s, s.biz)
}

func (s *Server) ListenAndServer() error {
	s.RegisterService()
	addr := fmt.Sprintf("%s:%d", config.GrpcConfig.Host, config.GrpcConfig.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("grpc server listening on %s.\n", addr)
	return s.s.Serve(l)
}

func (s *Server) Stop() error {
	log.Println("grpc server stop.")
	s.s.Stop()
	return nil
}
