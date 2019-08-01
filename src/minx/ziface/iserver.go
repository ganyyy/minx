package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 服务
	Serve()
	// 添加路由
	AddRouter(uint32, IRouter)
	// 获取连接管理器
	GetConnMgr() IConnManager

	// 设置连接创建时的回调
	SetOnConnStart(func (IConnection))
	// 设置连接关闭时的回调
	SetOnConnStop(func (IConnection))
	// 调用连接开始时的回调
	CallOnConnStart(IConnection)
	// 调用连接结束时的回调
	CallOnConnStop(IConnection)
}

