package main

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.gtcp.NewServer和gtcp.NewConn()
	1.1 gtcp.NewServer使用
		[go] gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn){
			defer conn.Close()
			for{
				data, err := conn.Recv(-1)
				xxxx
			}
		}).Run

	1.2 gtcp.NewConn("127.0.0.1:8300")
		conn, err := gtcp.NewConn("127.0.0.1:8300")
		defer conn.Close()
		if err != nil {}
		conn.SendRecv([]byte(), -1)

2.conn.LocalAddr() conn.RemoteAddr
	本地地址+端口： conn.LocalAddr
	远程地址+端口: conn.RemoteAddr

	对于客户端来说，本地地址+端口即为客户端地址+端口（客户端端口一般是随机的，因为客户端是于服务端的固定端口号建立连接）
	对于客户端来说，远程地址+端口即服务端地址+端口（服务端地址一般固定）

3.gtcp.NewConn("127.0.0.1:8300")和gtcp.NewPoolConn("127.0.0.1:8300")
	3.1 如果是NewConn，则每次发送数据后都会关闭，连接，下次会再次建立新的连接，所以打印出来的conn.LocalAddr中的端口号会不同
	例如：
		>2021-10-01 15:37:33 127.0.0.1:53224 127.0.0.1:8300
		>2021-10-01 15:37:34 127.0.0.1:53225 127.0.0.1:8300
		>2021-10-01 15:37:35 127.0.0.1:53240 127.0.0.1:8300
		>2021-10-01 15:37:36 127.0.0.1:53241 127.0.0.1:8300

	3.2 如果是NewPoolConn("127.0.0.1:8300"),则会创建一个连接池，本示例中client每经过1s发送一次数据，每次从
		连接池中获取的都是同一个连接（因为连接池会实现创建一部分连接放入连接池中，默认存活时间10分钟，等需要使用连接时直接
		从连接池中拿），即conn.LocalAddr是同一个端口号，每次发送完数据后conn.Close()也不会真正的关闭连接，而是将
		之前客户端与服务端之间已创建的连接放入连接池中(所以端口号相同)
	例如：
		>2021-10-01 15:36:41 127.0.0.1:54104 127.0.0.1:8300
		>2021-10-01 15:36:42 127.0.0.1:54104 127.0.0.1:8300
		>2021-10-01 15:36:43 127.0.0.1:54104 127.0.0.1:8300
*/

func main() {
	go gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte(">"), data...)); err != nil {
					fmt.Println(err)
				}
			}
			if err != nil {
				break
			}
		}
	}).Run()

	time.Sleep(time.Second)

	for {
		if conn, err := gtcp.NewPoolConn("127.0.0.1:8300"); err == nil {
			// (*PoolConn).SendRecv(data []byte, receive int) receive为-1表示接收所有数据直到连接中无数据时发送
			if b, err := conn.SendRecv([]byte(gtime.Datetime()), -1); err == nil {
				// conn.LocalAddr(): 本地网络地址
				// conn.RemoteAddr(): 远程网络地址
				fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
			} else {
				fmt.Println(err)
			}
			conn.Close()
		} else {
			glog.Error(err)
		}
		time.Sleep(time.Second)
	}
}
