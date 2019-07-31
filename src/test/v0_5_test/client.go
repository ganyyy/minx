package main

import (
	"../../minx/znet"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial error,", err)
		return
	}

	dp := znet.NewDataPack()
	
	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &znet.Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'w', 'o', 'r', 'l', 'd'},
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}

	// 粘包
	sendData1 = append(sendData1, sendData2...)

	cnt, err := conn.Write(sendData1)
	if err != nil {
		fmt.Println("client send error:", err)
		return
	}
	fmt.Printf("send %d data", cnt)

	select {}
}