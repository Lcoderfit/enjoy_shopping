package main

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.gtime.Datetime()
	返回string类型的当前时间，格式为: "2006-01-02 15:04:05"

2.gtcp.NewPoolConn原理
	2.1 首先客户端会发送gtime.Datetime到服务端
	2.2 服务端收到后会发送: "> 2021-10-01 16:23:50"格式的数据到客户端，之后关闭连接
	2.3 客户端收到服务端返回的数据后，会将该数据打印，同时会在同一行将conn.LocalAddr()和conn.RemoteAddr()也打印出来
	2.4 客户端打印完之后调用conn.Close()将当前连接放回原连接池,
	2.5 客户端第二次for循环，NewPoolConn由于连接池的IO复用特性会取到上一次的连接，但是由于服务端在第一次请求后就关闭连接了，
	所以第二请求写入成功（实际上数据未发送到Server，需要通过接下来的读取操作才能检测到链接错误），但客户端的读取操作失败了，这一次
	调用conn.Close()则会将该连接对象销毁而不会再重新放回连接池
	2.6 客户端第三次for循环，调用NewPoolConn重新创建或从链接池中获取新的链接，重复上诉步骤，
	2.7 打印情况如下：
		// 第一次请求成功
		> 2021-10-01 16:23:47 127.0.0.1:53779 127.0.0.1:8300
		// 第二次客户端请求时读取失败
		read tcp 127.0.0.1:53779->127.0.0.1:8300: wsarecv: An established connection was aborted by the software in your host machine.
		// 第三次请求成功
		> 2021-10-01 16:23:50 127.0.0.1:53780 127.0.0.1:8300
		read tcp 127.0.0.1:53780->127.0.0.1:8300: wsarecv: An established connection was aborted by the software in your host machine

3.连接对象重建
	由于连接池的IO复用特性（客户端每次conn.Close()时不会销毁当前连接而是放回链接池，下次使用时仍会使用上一次的连接）;
	所以需要格外注意当通信发生错误时，需要立即丢弃当前连接，重新创建新连接

4.conn.Recv和conn.SendRecv()
	4.1 conn.Send(n)
		4.1.1 n > 0: 从连接中收到长度为n字节的数据才返回
		4.1.2 n = 0: 从当前缓冲区接收数据并立即返回
		4.1.3 n <0: 接收连接中的所有数据直到没有数据时才返回
*/

func main() {
	// server
	go gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					fmt.Println(err)
				}
			}
			if err != nil {
				break
			}
			// 每次接收数据后会先发送一次数据，然后关闭当前连接
			return
		}
	}).Run()

	time.Sleep(time.Second)

	for {
		conn, err := gtcp.NewPoolConn("127.0.0.1:8300")
		if err != nil {
			glog.Error(err)

		} else {
			if b, err := conn.SendRecv([]byte(gtime.Datetime()), -1); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
			}
			conn.Close()
		}
		time.Sleep(time.Second)
	}
}
