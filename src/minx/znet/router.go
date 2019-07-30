package znet

import "../ziface"

// BaseRouter 所有路由的基类
type BaseRouter struct {}

/*
	直接实现这三种方法的原因是：
	有些路由可能不需要  任何处理 或者只需要 普通 处理.
*/
func (br *BaseRouter)PreHandle(request ziface.IRequest) {}
func (br *BaseRouter)Handle(request ziface.IRequest) {}
func (br *BaseRouter)PostHandle(request ziface.IRequest) {}

