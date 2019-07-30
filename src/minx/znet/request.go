package znet

import "../ziface"

// Request 请求的具体结构
type Request struct {
	conn ziface.IConnection
	data []byte
}

// GetConnection 获取请求的当前连接
func (r *Request)GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取当前连接的数据
func (r *Request)GetData() []byte {
	return r.data
}
