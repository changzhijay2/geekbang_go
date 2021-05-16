package mysql

import (
	"database/sql"
	"fmt"

	"github.com/changzhijay2/geekbang_go/week04/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() *sql.DB {
	mysqldb := config.MysqlConfig.DBName
	mysqlUserName := config.MysqlConfig.Username
	mysqlPwd := config.MysqlConfig.Password
	mysqlHost := config.MysqlConfig.Host
	mysqlPort := config.MysqlConfig.Port
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlUserName, mysqlPwd, mysqlHost, mysqlPort, mysqldb)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
