package biz

import (
	"context"

	v1 "github.com/changzhijay2/geekbang_go/week04/api/user/v1"
	"github.com/changzhijay2/geekbang_go/week04/internal/dao"
)

type Biz struct {
	dao *dao.Dao
	v1.UnimplementedUserServiceServer
}

func NewBiz(dao *dao.Dao) *Biz {
	return &Biz{dao: dao}
}

func (b *Biz) GetUserInfo (ctx context.Context, req *v1.UserRequest) (*v1.UserResponse, error) {
	resp := new(v1.UserResponse)
	id := req.UserId
	user, err := b.dao.GetUserInfo(int(id))
	if err != nil {
		return nil, err
	}
	resp.UserId = int32(user.UserId)
	resp.UserName = user.Username
	resp.Password = user.Password
	resp.Age = int32(user.Age)
	return resp, nil
}
