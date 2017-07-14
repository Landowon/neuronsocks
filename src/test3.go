// 这个里面暂时能实现sock5的握手

package main

import (
	// "bytes"

	"bytes"
	"errors"
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

type user_query struct {
	cmd    int
	atype  int
	target []byte
	port   int
}

func serve_client(C net.Conn) {
	fmt.Println("Accept from ", C.RemoteAddr().String())
	// var status = 0
	c := 0
	data := make([]byte, 2000)

	var remote_conn net.Conn = nil

	for {
		l, err := C.Read(data)
		// fmt.Println(c, " 客户端发来数据 HEX: ", data, "\n")
		if err != nil {
			fmt.Println("读取客户端数据错误")
			break
		}

		switch c {
		case 0: // 等待握手
			fmt.Println("step 1")
			handshake_step(C, data, l)
		case 1:
			fmt.Println("step 2")
			remote_conn, err = client_query(C, data, l)
			if err != nil {
				fmt.Println("远程连接失败", err.Error())
				break
			}
		case 2:
			// 在这儿链接新的地方
			fmt.Println("step 3")
			if remote_conn == nil {
				fmt.Println("无法连接远程主机")
				break
			}
			err = get_data(C, remote_conn, data, l)
			if err != nil {
				fmt.Println("传输数据错误", err.Error())
			}
			break
		default:
			c = 10
		}

		c += 1
		if c > 5 {
			break
		}
	}
	fmt.Println("step z")

	defer C.Close()
}

// 处理客户端handshake
func handshake_step(C net.Conn, buf []byte, l int) (err error) {
	// reader := bytes.NewReader(buf)
	proto_version := buf[0]
	if proto_version != '\x05' {
		fmt.Println("协议的版本不是05的，而是:", string(proto_version))
		return errors.New("proto error!")
	}
	fmt.Println("协议的版本是05版")
	num_methods := int(buf[1])
	fmt.Println("一共有这么多种方式:", num_methods)
	if num_methods+2 > l {
		// 断包什么的
		fmt.Println("包长度不对!")
		return errors.New("invalid packet!")
	}
	methods := buf[2 : 2+num_methods]
	met_0_inarray := false
	for i := range methods {
		if i == '\x00' {
			met_0_inarray = true
			break
		}
	}
	if len(methods) == 0 || !met_0_inarray {
		fmt.Println("没有合适的方法!")
		return errors.New("no valid method!")
	}
	C.Write([]byte{'\x05', '\x00'})
	return nil
}

//  处理客户端query
func client_query(C net.Conn, buf []byte, l int) (net.Conn, error) {
	proto_version := buf[0]
	if proto_version != '\x05' {
		fmt.Println("协议的版本不是05的，而是:", string(proto_version))
		return nil, errors.New("proto error!")
	}
	fmt.Println("协议的版本是05版")
	cmd := int(buf[1])

	switch cmd {
	case 1:
	case 2:
		fallthrough // 不知道怎么写
	case 3:
		fallthrough
	default:
		fmt.Println("未知命令:", cmd)
		return nil, errors.New("unknown cmd!")
	}

	// rsv := buf[2]
	// atyp := int(buf[3])
	dstlen := int(buf[4])
	dst := buf[5 : 5+dstlen]
	port := buf[5+dstlen : 5+dstlen+2]

	//  接下来是一堆处理奇葩情况的代码
	// fmt.Print(string(dst), port)

	// return &user_query{cmd, atyp, dst, 256*int(port[0]) + int(port[1])}, nil

	ns, err := net.LookupHost(string(dst))
	if err != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return nil, errors.New("no valid ns")
	}

	good_ns := ns[0]

	fmt.Println("ns is:", good_ns)

	RC, err := net.DialTCP("tcp", nil, &net.TCPAddr{net.ParseIP(good_ns), 256*int(port[0]) + int(port[1]), ""})

	if err != nil {
		return nil, err
	}
	// 欠一个回应
	C.Write(bytes.Join([][]byte{[]byte{'\x05', '\x00', '\x00'}, buf[3:l]}, []byte{}))
	return RC, nil
}

// func handle_connect_cmd()

//  处理客户端query

func get_data(C net.Conn, RC net.Conn, clidata []byte, l int) error {
	defer RC.Close()
	_, err := RC.Write(clidata[:l])

	// RC.SetReadDeadline(time.a)

	fmt.Println("step x", string(clidata[:l]))
	if err != nil {
		return err
	}
	remdata := make([]byte, 20000)
	for {
		rl, err := RC.Read(remdata)
		// fmt.Println("step y", rl, RC.RemoteAddr().String()) // string(remdata[:rl]))
		if err != nil {
			return err
		}

		_, err = C.Write(remdata[:rl])
		if err != nil {
			return err
		}
		// return nil
	}
	return nil
}

// func ipstring_to_bytearray(s string) ([]byte, error) {

// }
