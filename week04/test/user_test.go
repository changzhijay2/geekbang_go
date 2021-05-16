package test

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/changzhijay2/geekbang_go/week04/api/user/v1"
	"google.golang.org/grpc"
)

func TestGetRouter(t *testing.T) {
	serviceAddress := "127.0.0.1:5567"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	userClient := v1.NewUserServiceClient(conn)
	userReq := &v1.UserRequest{UserId: 1}
	userResp, err := userClient.GetUserInfo(context.Background(), userReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UserService GetUserInfo : Name %s, Password %s, Age %d\n", userResp.UserName, userResp.Password, userResp.Age)
}
