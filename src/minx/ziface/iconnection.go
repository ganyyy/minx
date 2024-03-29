package ziface

import "net"

// IConnection 客户端连接类
type IConnection interface {
	// 启动连接时的处理
	Start()
	// 关闭连接时的处理
	Stop()
	// 获取连接套接字
	GetTCPConnection() *net.TCPConn
	// 获取连接Id
	GetConnID() uint32
	// 获取远端地址
	RemoteAddr() net.Addr
	// 将数据发送给客户端
	SendMsg(uint32, []byte) error
	// 带缓冲的发送方式
	SendBuffMsg(uint32, []byte) error
}

// HandFun 统一的客户端处理函数接口
type HandFun func(*net.TCPConn, []byte, int) error
