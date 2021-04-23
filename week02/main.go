package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

/*
1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

答：应该抛给上层，Dao层只负责数据访问，具体错误应交给Service层的调用者处理。
*/

type Dao struct {
	db *sql.DB
}

func NewDao() *Dao {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/school")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return &Dao{
		db: db,
	}
}

func (d *Dao) GetNameByID(id int) (string, error) {
	var name string
	err := d.db.QueryRow("select name from student where id = ?", id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.Wrap(err, "")
		} else {
			return "", errors.Wrap(err, fmt.Sprintf("sql: You have an error in your SQL syntax, id = %d", id))
		}
	}
	return name, nil
}

func (d *Dao) Close() {
	d.db.Close()
}

func main() {
	dao := NewDao()
	defer dao.Close()
	name, err := dao.GetNameByID(9527)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	fmt.Println(name)
}
