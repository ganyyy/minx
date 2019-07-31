package main

import (
	"../../minx/znet"
	"fmt"
	"io"
	"net"
)

func main()  {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server listen error:", err)
		return
	}

	// 创建服务器goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			return
		}

		go func(conn net.Conn) {
			dp := znet.NewDataPack()

			for {
				headData := make([]byte, dp.GetHeadLen())
				// 先读头, 8个字节
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("Read client head error,", err)
					return
				}

				// 将头解包
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("Unpack error,", err)
					return
				}

				// 读数据
				if msgHead.GetDataLen() > 0 {
					// 这是啥用法？
					// 类型断言, A.(T)
					// A只能是接口, T可以是接口也可以是类型
					msg, err := msgHead.(*znet.Message)
					if !err {
						fmt.Println("Error Message type")
						return
					}
					msg.Data = make([]byte, msgHead.GetDataLen())
				}

			}
		}(conn)
	}
}