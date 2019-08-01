package znet

import (
	"../ziface"
	"fmt"
)

// MsgHandle 消息管理结构
type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

// DoMsgHandle 处理
func (mh *MsgHandle) DoMsgHandle(req ziface.IRequest) {
	router, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		fmt.Println("Error msg id:", req.GetMsgID())
		panic("DoMsgHandle error")
	}

	// 具体处理
	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}

// AddRouter 添加路由
func (mh *MsgHandle) AddRouter(id uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[id]; ok {
		fmt.Println("Repeat id:", id)
		panic("AddRouter error")
	}

	mh.Apis[id] = router
	fmt.Printf("Add new router, id:%d\n", id)
}

// NewMsgHandle 获取一个新的消息管理
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
