package znet

import (
	"../ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection  // 所有连接
	connLock sync.RWMutex // 读写锁
}

// NewConnManager 获取一个连接管理器
func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加一个新的连接
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 添加一个连接
	c.connections[conn.GetConnID()] = conn
	fmt.Printf("conn add to connManager success, conn id: %d, total len:%d\n", conn.GetConnID(), c.Len())
}

// Remove 删除连接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 加锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除连接
	delete(c.connections, conn.GetConnID())
	fmt.Printf("remove conn:%d from connManager, num:%d\n", conn.GetConnID(), c.Len())
}

// Get 获取连接
func (c *ConnManager) Get(id uint32) (ziface.IConnection, error) {
	// 不允许读
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[id]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn id not found")
	}
}

// Len 获取当前连接的数量
func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) RemoveAll() {
	// 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for id, conn := range c.connections {
		// 先停了
		conn.Stop()
		// 在删了
		delete(c.connections, id)
	}

	fmt.Printf("remove all conn success, conn num = %d\n", c.Len())
}

