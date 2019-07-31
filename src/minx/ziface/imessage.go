package ziface

// IMessage 通用消息接口
type IMessage interface {
	GetDataLen() uint32  // 获取消息长度
	GetMsgId() uint32 // 获取消息ID
	GetData() []byte  // 获取消息二进制内容
	SetDataLen(uint32)   // 设置消息长度
	SetMsgId(uint32)  // 设置消息ID
	SetData([]byte)   // 设置消息二进制内容
}
