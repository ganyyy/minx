package ziface

// IRequest 请求的接口.
type IRequest interface {
	GetConnection() IConnection // 获取客户端的连接
	GetData() []byte            // 获取客户端传入的数据
	GetMsgID() uint32           // 获取消息ID

	// TODO 后期可以添加自己需要的内容
}
