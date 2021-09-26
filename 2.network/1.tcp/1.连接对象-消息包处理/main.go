package main

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"time"
)

func main() {
	go gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("receive: ", data)
		}
	}).Run()

	time.Sleep(time.Second)
	// Client
	conn, err := gtcp.NewConn("127.0.0.1:8300")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 0; i < 1000; i++ {
		if err := conn.SendPkg([]byte(gconv.String(i))); err != nil {
			glog.Error(err)
		}
		time.Sleep(time.Second)
	}
}
