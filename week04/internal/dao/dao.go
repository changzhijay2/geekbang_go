package dao

import (
	"database/sql"

	"github.com/changzhijay2/geekbang_go/week04/internal/model"
)

type Dao struct {
	db *sql.DB
}

func NewDao(db *sql.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (d *Dao) GetUserInfo(id int) (*model.User, error) {
	user := new(model.User)
	row := d.db.QueryRow("select * from user where user_id = ?", id)
	err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.Age)
	if err != nil {
		return nil, err
	}
	return user, nil
}
