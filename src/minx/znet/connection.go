package znet

import (
	"errors"
	"fmt"
	"io"
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

	// 处理路由管理
	msgHandle ziface.IMsgHandle

	// 读写管道, 无缓冲
	msgChan chan []byte
}

// NewConnection 返回一个新的客户端连接结构体
func NewConnection(conn *net.TCPConn, connID uint32, mh ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClose:      false,
		ExitBuffChan: make(chan bool, 1),
		msgHandle:    mh,
		msgChan:      make(chan []byte),
	}
	return c
}

// Start 启动
func (c *Connection) Start() {
	// 开启读业务
	go c.StartRead()

	// 开启写业务
	go c.StartWrite()

	fmt.Println("Conn Start")
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

// SendMsg 发送给客户端消息
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("Connection is close ")
	}

	dp := NewDataPack()

	// 包装
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg err, id:", msgId)
		return errors.New("Pack message error ")
	}

	// 写回客户端, 通过管道的形式通知写
	c.msgChan <- msg

	// 没错
	return nil
}

// StartRead 读线程
func (c *Connection) StartRead() {
	fmt.Println("Reader Goroutine is running")

	// 收尾工作
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("conn error:", err)
		}
		fmt.Println(c.RemoteAddr().String(), " conn reader exit")
		c.Stop()
	}()

	// 声明一个解包对象
	dp := NewDataPack()

	// 头部数据
	headData := make([]byte, dp.GetHeadLen())

	// 死循环读线程
	for {
		if c.isClose {
			// 如果出现问题并关闭, 则直接返回
			return
		}
		// 读头
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("client read data head err:", err)
			return
		}

		// 解头
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack head data err:", err)
			return
		}

		// 获取实际数据
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("client read data err:", err)
			}
		}

		// 放入数据
		msg.SetData(data)

		// 将当前连接包装成一个Request
		req := Request{
			conn: c,
			data: msg,
		}

		// 处理请求
		go c.msgHandle.DoMsgHandle(&req)
	}
}

// StartWrite 写线程分离
func (c *Connection) StartWrite() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("write error:", err)
		}
		fmt.Println("conn write exit")
		c.Stop()
	}()

	for {
		select {
		case data := <- c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("client write err:", err)
				return
			}
		case <-c.ExitBuffChan:
			// 已关闭
			fmt.Println("conn recv write exit sig")
			return
		}
	}
}
