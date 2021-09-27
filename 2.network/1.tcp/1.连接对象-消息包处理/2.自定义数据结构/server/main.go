package main

import (
	"encoding/json"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"learn-gf/2.network/1.tcp/1.连接对象-消息包处理/2.自定义数据结构/types"
)

/*
如果不使用package main就用gf run main.go运行，则会报错：
	build running error: fork/exec bin\main.exe : This version of %1 is not compatible with the version of Windows you're runni
ng. Check your computer's system information and then contact the software publisher.

如果不使用package main 就直接 go run main.go，则会报错：
	go run: cannot run non-main package
*/

func main() {
	// 服务端，接收客户端数据并格式化为指定数据结构，打印
	// 注意，gtcp.NewServer返回的是一个*Server对象，所以需要调用Run方法监听端口
	gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn) {
		defer conn.Close()
		// 循环读取
		for {
			// 接收数据
			data, err := conn.RecvPkg()
			if err != nil {
				if err.Error() == "EOF" {
					glog.Println("client closed")
				}
				break
			}
			info := &types.NodeInfo{}
			if err := json.Unmarshal(data, info); err != nil {
				glog.Println("invalid package structure: %s", err.Error())
			} else {
				glog.Println(info)
				conn.SendPkg([]byte("ok"))
			}
		}
	}).Run()
}
