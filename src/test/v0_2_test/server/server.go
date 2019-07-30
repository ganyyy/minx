package main

import "../../../minx/znet"

func main() {
	// 创建一个服务器
	s := znet.NewServer("[zinx v0.2]")

	// 开启服务
	s.Serve()
}


