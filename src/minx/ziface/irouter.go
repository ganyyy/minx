package ziface

// IRoute 路由的通用接口, 每一个业务处理都需要实现该接口
type IRouter interface {
	PreHandle(request IRequest)  // 处理之前
	Handle(request IRequest)     // 实际处理
	PostHandle(request IRequest) // 处理之后
}
