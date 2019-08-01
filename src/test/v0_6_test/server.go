package main

import (
	"../../minx/ziface"
	"../../minx/znet"
	"fmt"
	"os"
)

// PingRouter Router1
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Printf("===> Recv Msg: ID = %d, data = %s\n", request.GetMsgID(), string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call handle error")
		panic("handle error")
	}
}

// HelloRouter Router2
type HelloRouter struct {
	znet.BaseRouter
}

func (p *HelloRouter) Handle(request ziface.IRequest)  {
	fmt.Printf("===> Recv Msg: ID = %d, data = %s\n", request.GetMsgID(), string(request.GetData()))
	err := request.GetConnection().SendMsg(2, []byte("hello world\n"))
	if err != nil {
		fmt.Println("call handle error")
		panic("handle error")
	}
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Init error:", err)
		}
	}()

	err := os.Setenv("ZINX_ROOT_PATH", "E:/Code/zinx/src")
	if err != nil {
		fmt.Println("Error env set")
		return
	}

	// 创建一个服务器
	s := znet.NewServer("[zinx v0.5]")

	// 添加路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 开启服务
	s.Serve()
}
