package ziface

// IConnManager 连接管理抽象接口
type IConnManager interface {
	Add(IConnection)                 // 添加一个连接
	Remove(IConnection)              // 移除一个连接
	Get(uint32) (IConnection, error) // 根据ID获取连接
	Len() int                        // 返回连接的数量
	RemoveAll()                      // 清空所有连接
}
