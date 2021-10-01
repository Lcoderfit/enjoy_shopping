package main

import (
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"learn-gf/2.network/1.tcp/2.连接对象-异步全双工通信-TCP长链接/funcs"
	"learn-gf/2.network/1.tcp/2.连接对象-异步全双工通信-TCP长链接/types"
	"time"
)

/*
1.gtimer.SetTimeout(delay, job)
	经过delay时间后执行一次job函数

2.gtimer.SetInterval(interval, job)
	没经过interval时间后执行一次job

3.异步全双工通信
*/

func main() {
	gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn) {
		defer conn.Close()
		// 测试消息，10s后让客户端主动退出
		// SetTimeout：传入两个参数，第一个是超时时间，第二个是func(){}函数，作用是经过超时时间后运行一次函数
		gtimer.SetTimeout(10*time.Second, func() {
			funcs.SendPkg(conn, "doexit")
		})
		for {
			msg, err := funcs.RecvPkg(conn)
			if err != nil {
				if err.Error() == "EOF" {
					glog.Println("client closed")
				}
				break
			}
			switch msg.Act {
			case "hello":
				onClientHello(conn, msg)
			case "heartbeat":
				onClientHeartBeat(conn, msg)
			default:
				glog.Errorf("invalid message: %v", msg)
				break
			}
		}
	})
}

func onClientHello(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("hello message from [%s]: %s", conn.RemoteAddr().String(), msg.Data)
	funcs.SendPkg(conn, msg.Act, "Nice to meet you")
}

func onClientHeartBeat(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("heartbeat from [%s]", conn.RemoteAddr().String())
}
