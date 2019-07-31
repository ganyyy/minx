package main

import (
	"../../../minx/ziface"
	"../../../minx/znet"
	"fmt"
	"os"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call handle error")
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
	s := znet.NewServer("[zinx v0.2]")

	// 添加一个路由
	s.AddRouter(&PingRouter{})

	// 开启服务
	s.Serve()
}


