package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"week09/conf"
	"week09/server"

	"golang.org/x/sync/errgroup"
)

type Program struct {
	tcpServer *server.Server
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewProgram() *Program {
	ctx, cancel := context.WithCancel(context.Background())
	return &Program{
		tcpServer: server.NewServer(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (p *Program) Run() error {
	g, ctx := errgroup.WithContext(p.ctx)
	g.Go(func() error {
		return p.tcpServer.ListenAndServer()
	})
	g.Go(func() error {
		<-ctx.Done()
		return p.tcpServer.Stop()
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				p.Stop()
			}
		}
	})
	return g.Wait()
}

func (p *Program) Stop() error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags)
	var confFilePath string
	flag.StringVar(&confFilePath, "conf", "../../conf.yaml", "conf file path.")
	flag.Parse()
	conf.InitConf(confFilePath)

	program := NewProgram()
	if err := program.Run(); err != nil {
		if !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}
}
