package ziface

// IMsgHandle 消息管理
type IMsgHandle interface {
	DoMsgHandle(IRequest)      // 非阻塞处理请求
	AddRouter(uint32, IRouter) // 添加一个新路由
}
