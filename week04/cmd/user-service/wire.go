
//+build wireinject

package main

import (
	"github.com/changzhijay2/geekbang_go/week04/internal/biz"
	"github.com/changzhijay2/geekbang_go/week04/internal/dao"
	"github.com/changzhijay2/geekbang_go/week04/internal/server"
	"github.com/changzhijay2/geekbang_go/week04/pkg/mysql"
	"github.com/google/wire"
)

func InitServer() (*server.Server, error) {
	wire.Build(mysql.NewDB, dao.NewDao, biz.NewBiz, server.NewServer)
	return &server.Server{}, nil
}
