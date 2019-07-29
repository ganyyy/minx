package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 服务
	Server()
}

