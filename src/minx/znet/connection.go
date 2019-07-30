package znet

import (
	"../ziface"
	"fmt"
	"net"
)

//Connection 客户端连接结构
type Connection struct {
	// 套接字
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 是否已关闭
	isClose bool

	// 连接处理方法
	handleAPI ziface.HandFun

	// 退出消息通知chan
	ExitBuffChan chan bool
}

// NewConnection 返回一个新的客户端连接结构体
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandFun) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		handleAPI:    callbackApi,
		isClose:      false,
		ExitBuffChan: make(chan bool, 1),
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
func (c *Connection) GetTCPConnection() net.Conn {
	return c.Conn
}

// GetConnID 获取连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAttr 获取连接地址
func (c *Connection) RemoteAttr() net.Addr {
	return c.Conn.RemoteAddr()
}

// StartReader 读线程
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")

	// 收尾
	defer fmt.Println(c.RemoteAttr().String(), " conn reader exit")
	defer c.Stop()

	// 死循环读线程
	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		// 如果读出错
		if err != nil {
			fmt.Println("Reader Error:", err)
			return
		}

		// 如果处理出错
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("Handle Error:", err)
			return
		}
	}

}
