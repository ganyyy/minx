package main

import (
	"../../../minx/ziface"
	"../../../minx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRouter
}

//func (p *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
//	if err != nil {
//		fmt.Println("call pre error")
//	}
//}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call handle error")
	}
}

//func (p *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call router PostHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping\n"))
//	if err != nil {
//		fmt.Println("call post error")
//	}
//}

func main() {

	//err := os.Setenv("ZINX_ROOT_PATH", )

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Init error:", err)
		}
	}()

	// 创建一个服务器
	s := znet.NewServer("[zinx v0.2]")

	// 添加一个路由
	s.AddRouter(&PingRouter{})

	// 开启服务
	s.Serve()
}


