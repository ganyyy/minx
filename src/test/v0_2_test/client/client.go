package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test Start ...")

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp4", "127.0.0.1:7777")

	if err != nil {
		fmt.Println("client dial error:", err)
		return
	}

	for i := 0; i < 4; i++ {
		_, err := conn.Write([]byte("Hello World!"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}

		buf := make([]byte, 512)

		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		fmt.Printf("server callback: %s, count: %d\n", buf[:cnt], cnt)

		time.Sleep(time.Second)
	}
}
