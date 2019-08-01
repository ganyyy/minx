package znet

import (
	"../utils"
	"../ziface"
	"fmt"
)

// MsgHandle 消息管理结构
type MsgHandle struct {
	Apis map[uint32]ziface.IRouter		// 每一个业务对应的处理
	WorkPoolSize uint32					// 工作线程的数量
	TaskQueue []chan ziface.IRequest	// 任务队列管道, 每一个工作线程配一个
}

// StartOneWork 启动一个工作池
func (mh *MsgHandle) StartOneWork(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Worker ID:%d is start\n", workerID)
	// 暂定work不会退出, 会一直等待并处理请求
	for {
		select {
		case request := <- taskQueue:
			// 等待接收请求并处理
			mh.DoMsgHandle(request)
		}
	}
}

// StartWorkPool 启动工作池
func (mh *MsgHandle) StartWorkPool() {
	maxTaskLen := utils.GlobalObject.MaxWorkTaskLen
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		// 给指定work创建消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, maxTaskLen)
		// 启动消息处理
		go mh.StartOneWork(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 添加请求到工作队列
func (mh *MsgHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	// 根据连接ID来确定交给哪个work处理
	connID := req.GetConnection().GetConnID()
	workID := connID % mh.WorkPoolSize
	fmt.Printf("Add ConnID:%d to Work %d, request msgID:%d\n", connID, workID, req.GetMsgID())
	// 将请求放入到指定队列中
	mh.TaskQueue[workID] <- req
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
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}
