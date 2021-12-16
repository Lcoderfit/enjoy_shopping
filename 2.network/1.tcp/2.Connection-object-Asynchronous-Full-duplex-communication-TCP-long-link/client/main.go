package client

import (
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"learn-gf/2.network/1.tcp/2.Connection-object-Asynchronous-Full-duplex-communication-TCP-long-link/funcs"
	"learn-gf/2.network/1.tcp/2.Connection-object-Asynchronous-Full-duplex-communication-TCP-long-link/types"
	"time"
)

func main() {
	conn, err := gtcp.NewConn("127.0.0.1:8300")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 心跳消息
	gtimer.SetInterval(time.Second, func() {
		if err := funcs.SendPkg(conn, "heartbeat"); err != nil {
			panic(err)
		}
	})
	//
	gtimer.SetTimeout(3*time.Second, func() {
		if err := funcs.SendPkg(conn, "hello", "my name's john!"); err != nil {
			panic(err)
		}
	})
	for {
		msg, err := funcs.RecvPkg(conn)
		if err != nil {
			if err.Error() == "EOF" {
				glog.Println("server1 closed")
			}
			break
		}

		switch msg.Act {
		case "hello":
			onServerHello(conn, msg)
		case "heartbeat":
			onServerHeartBeat(conn, msg)
		case "doexit":
			onServerDoExit(conn, msg)
		default:
			glog.Errorf("invalid package: %v", msg)
			break
		}
	}
}

func onServerHello(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("hello response message from [%s]: %s", conn.RemoteAddr().String(), msg.Data)
}

func onServerHeartBeat(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("heartbeat from [%s]", conn.RemoteAddr().String())
}

func onServerDoExit(conn *gtcp.Conn, msg *types.Msg) {
	glog.Println("exit command from [%s]", conn.RemoteAddr().String())
	conn.Close()
}
