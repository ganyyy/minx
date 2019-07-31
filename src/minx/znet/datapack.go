package znet

import (
	"../utils"
	"../ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// DataPack 拆包封包结构
type DataPack struct{}

// GetHeadLen 获取头部长度, 固定为8字节
func (db *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

// Pack 包装, 将Message包装成二进制数据流
func (db *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建缓冲区
	dataBuf := bytes.NewBuffer([]byte{})

	// 以小端格式进行包装

	// 长度
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// ID
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 数据
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

// Unpack 解包, 将二进制数据流转为Message
func (db *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(data)

	msg := &Message{}

	// 开始解包

	// 长度
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// ID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断数据长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data length. ")
	}

	// 将解包分成两个过程:
	// 1.先解出头部来确定数据的长度
	// 2.根据头部的信息在从buf中解出数据
	return msg, nil
}

// NewDatePack 获取一个新的拆封包结构
func NewDataPack() *DataPack {
	return &DataPack{}
}
