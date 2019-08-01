package znet

import (
	"../utils"
	"../ziface"
	"errors"
	"fmt"
	"io"
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
	// 退出消息通知chan
	ExitBuffChan chan bool
	// 处理路由管理
	msgHandle ziface.IMsgHandle
	// 读写管道, 无缓冲
	msgChan chan []byte
	// 带缓冲的管道
	msgBufChan chan []byte
	// 服务器的引用
	TCPServer ziface.IServer
}

// NewConnection 返回一个新的客户端连接结构体
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, mh ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClose:      false,
		ExitBuffChan: make(chan bool, 1),
		msgHandle:    mh,
		msgChan:      make(chan []byte),
		msgBufChan:   make(chan []byte, utils.GlobalObject.MaxPacketSize),
	}
	// 将连接添加到管理器中
	server.GetConnMgr().Add(c)
	return c
}

// Start 启动
func (c *Connection) Start() {
	// 开启读业务
	go c.StartRead()

	// 开启写业务
	go c.StartWrite()

	// 连接开始的回调
	c.TCPServer.CallOnConnStart(c)
}

// Stop 停止
func (c *Connection) Stop() {
	if c.isClose {
		return
	}
	c.isClose = true

	// 连接结束的回调
	c.TCPServer.CallOnConnStop(c)

	// 清理服务器的连接管理器
	c.TCPServer.GetConnMgr().Remove(c)

	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Close Error:", err)
	}
	// 通知管道进行关闭
	c.ExitBuffChan <- true
	// 关闭管道
	close(c.ExitBuffChan)
	close(c.msgBufChan)
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

// SendBuffMsg 带缓冲的发送方式
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("Connection is close ")
	}
	dp := NewDataPack()

	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg err, id:", msgId)
		return errors.New("Pack message error ")
	}
	c.msgBufChan <- msg
	return nil
}

// StartRead 读线程
func (c *Connection) StartRead() {
	defer fmt.Println(c.RemoteAddr().String(), "[conn reader exit]")
	// 收尾工作
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("conn error:", err)
		}
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
		req := &Request{
			conn: c,
			data: msg,
		}

		// 如果存在work pool就交给其来处理, 否则就正常处理
		if utils.GlobalObject.WorkPoolSize > 0 {
			c.msgHandle.SendMsgToTaskQueue(req)
		} else {
			// 处理请求
			go c.msgHandle.DoMsgHandle(req)
		}

	}
}

// StartWrite 写线程分离
func (c *Connection) StartWrite() {
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit]")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("write error:", err)
		}
		c.Stop()
	}()

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("client write err:", err)
				return
			}
		case data, ok := <-c.msgBufChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("client write buf err:", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is close")
				break
			}
		case <-c.ExitBuffChan:
			// 已关闭
			return
		}
	}
}
