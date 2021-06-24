package main

import (
	"bufio"
	"fmt"
	"net"
	"week09/util"
)

func main() {
	pingMsg := util.NewMessage(1, 1, 1, []byte("ping."))
	b1, err := util.Pack(pingMsg)
	if err != nil {
		panic(err)
	}
	helloMsg := util.NewMessage(1, 2, 2, []byte("hello world."))
	b2, err := util.Pack(helloMsg)
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("tcp4", "127.0.0.1:5576")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if _, err := conn.Write(b1); err != nil {
		panic(err)
	}
	r := bufio.NewReader(conn)
	recvMsg1, err := util.Unpack(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(recvMsg1.GetBody()))

	if _, err := conn.Write(b2); err != nil {
		panic(err)
	}
	recvMsg2, err := util.Unpack(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(recvMsg2.GetBody()))
}
