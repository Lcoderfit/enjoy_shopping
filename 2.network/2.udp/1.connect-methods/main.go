package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.创建udp客户端和服务端(除了名字不一样，用法跟tcp非常相似)
	1.1 gudp.NewServer()
	1.2 gudp.NewConn()

2.gupd与gtcp建立的连接的不同
	2.1 gtcp可以保持长连接，而且可以创建连接池，但是gudp不行（因为是面向无连接的协议），无法与服务端保持连接，
	每次通信后必须创建新的连接
	> 2021-10-01 17:20:58 127.0.0.1:60633 127.0.0.1:8300
	> 2021-10-01 17:20:59 127.0.0.1:60634 127.0.0.1:8300
	> 2021-10-01 17:21:00 127.0.0.1:54546 127.0.0.1:8300
*/

func main() {
	go gudp.NewServer("127.0.0.1:8300", func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					g.Log().Line(true).Error(err)
				}
			}
			if err != nil {
				g.Log().Line(true).Error(err)
			}
		}
	}).Run()

	time.Sleep(time.Second)

	// Client
	for {
		if conn, err := gudp.NewConn("127.0.0.1:8300"); err != nil {
			g.Log().Line(true).Error(err)
		} else {
			if b, err := conn.SendRecv([]byte(gtime.Datetime()), -1); err != nil {
				g.Log().Line(true).Error(err)
			} else {
				fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
			}
			conn.Close()
		}
		time.Sleep(time.Second)
	}
}
