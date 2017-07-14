package main

import (
	// "bytes"
	"fmt"
	"net"
)

func main() {
	// for{
	// conn, err := listen.Accept
	// }
	fmt.Println("服务器启动！")

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("0.0.0.0"), 8888, ""})
	if err != nil {
		fmt.Println("监听端口失败!")
		return
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接收客户端连接异常!")
			continue
		}
		// defer func(conn net.Conn){conn.Close(); fmt.Println("defer已经执行")}(conn)	// 不知道这个defer在崩溃的时候是否会执行
		go serve_client(conn)

	}

}

func serve_client(C net.Conn) {
	fmt.Println("Accept from ", C.RemoteAddr().String())
	c := 0
	data := make([]byte, 128)
	for {
		_, err := C.Read(data)
		fmt.Println(c, " 客户端发来数据 HEX: ", data, "\n")
		if err != nil {
			fmt.Println("读取客户端数据错误")
			break
		}
		C.Write([]byte{'\x05', '\x00'})
		c += 1
		if c > 5 {
			break
		}
	}
	C.Write([]byte{'\x05', '\x00'})
	defer C.Close()
}
