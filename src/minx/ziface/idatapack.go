package ziface

// IDataPack 数据包和解
type IDataPack interface {
	GetHeadLen() uint32              // 获取消息长度
	Pack(IMessage) ([]byte, error)   // 消息转二进制流
	Unpack([]byte) (IMessage, error) // 二进制流转消息
}
