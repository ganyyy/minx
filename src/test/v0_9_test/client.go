package main

import (
	"../../minx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client dial error:", err)
		return
	}

	dp := znet.NewDataPack()
	headData := make([]byte, dp.GetHeadLen())

	for {
		sendData, err := dp.Pack(znet.NewMessage(1, []byte("Zinx 0.9 Client Test Message")))
		if err != nil {
			fmt.Println("Client pack message err:", err)
			return
		}
		cnt, err := conn.Write(sendData)
		if err != nil {
			fmt.Println("Client write err:", err)
			return
		}
		fmt.Println("Send bytes count:", cnt)

		// 准备解包
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("Client read head err:", err)
			return
		}

		recvMsg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Client unpack err:", err)
			return
		}

		if recvMsg.GetDataLen() > 0 {
			// 直接SetData试试行不行
			recvMsg.SetData(make([]byte, recvMsg.GetDataLen()))

			_, err = io.ReadFull(conn, recvMsg.GetData())
			if err != nil {
				fmt.Println("Client read data err:", err)
				return
			}
			fmt.Printf("===> Recv Msg: ID = %d, len = %d, data = %s\n", recvMsg.GetMsgId(), recvMsg.GetDataLen(), string(recvMsg.GetData()))
		}

		time.Sleep(time.Second)
	}

}
