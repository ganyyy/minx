package ziface

// IMsgHandle 消息管理
type IMsgHandle interface {
	DoMsgHandle(IRequest)        // 非阻塞处理请求
	AddRouter(uint32, IRouter)   // 添加一个新路由
	StartWorkPool()              // 启动WorkPool
	SendMsgToTaskQueue(IRequest) // 将请求加入到任务队列中, 等待工作池处理
}
