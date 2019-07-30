package znet

import (
	"fmt"
	"net"

	"../ziface"
)

//Connection 客户端连接结构
type Connection struct {
	// 套接字
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 是否已关闭
	isClose bool

	// 退出消息通知chan
	ExitBuffChan chan bool

	// 处理路由
	Router ziface.IRouter
}

// NewConnection 返回一个新的客户端连接结构体
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClose:      false,
		ExitBuffChan: make(chan bool, 1),
		Router:       router,
	}
	return c
}

// Start 启动
func (c *Connection) Start() {
	// 开启读业务请求
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			// 等待管道输入, 然后结束处理
			return
		}
	}
}

// Stop 停止
func (c *Connection) Stop() {
	if c.isClose {
		return
	}
	c.isClose = true

	// TODO 用户关闭连接的回调处理

	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Close Error:", err)
	}
	// 通知管道进行关闭
	c.ExitBuffChan <- true
	// 关闭管道
	close(c.ExitBuffChan)
}

// GetTCPConnection 获取连接套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取连接地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// StartReader 读线程
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")

	// 收尾
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()

	// 死循环读线程
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		// 如果读出错
		if err != nil {
			fmt.Println("Reader Error:", err)
			return
		}

		// 将当前连接包装成一个Request
		req := Request{
			conn: c,
			data: buf,
		}

		// 使用路由处理请求
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}

}
