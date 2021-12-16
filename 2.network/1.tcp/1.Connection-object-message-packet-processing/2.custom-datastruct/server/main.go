package main

import (
	"encoding/json"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"learn-gf/2.network/1.tcp/1.Connection-object-message-packet-processing/2.custom-datastruct/types"
)

/*
1.运行报错：
	如果不使用package main就用gf run main.go运行，则会报错：
		build running error: fork/exec bin\main.exe : This version of %1 is not compatible with the version of Windows you're runni
	ng. Check your computer's system information and then contact the software publisher.

	如果不使用package main 就直接 go run main.go，则会报错：
		go run: cannot run non-main package

2.内置方法
type Conn
    func (c *Conn) SendPkg(data []byte, option ...PkgOption) error
    func (c *Conn) SendPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) error
    func (c *Conn) SendRecvPkg(data []byte, option ...PkgOption) ([]byte, error)
    func (c *Conn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error)
    func (c *Conn) RecvPkg(option ...PkgOption) (result []byte, err error)
    func (c *Conn) RecvPkgWithTimeout(timeout time.Duration, option ...PkgOption) ([]byte, error)

注意：有data参数的，均为[]byte类型

SendPkg: 发送包
SendPkgWithTimeout: 发送包，并对数据包设置过期时间
SendRecvPkg 发送包并接收服务端响应
SendRecvPkgWithTimeout: 发送报并接收响应，设置接收响应时的超时时间
SendPkg: 发送数据包并且设置超时时间
RecvPkgWithTimeout: 接收数据包并设置超时时间

3.服务端
	gtcp.NewServer("127.0.0.1:8300", func(conn *gtcp.Conn{
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				if err.Error() == "break" {
					glog.Println("client closed")
				}
			}
		}
	})): 服务端监听端口，返回*Server实例(需要调用.Run方法启动监听程序)
	gtcp.NewConn("127.0.0.1:8300") 客户端建立连接

4.客户端
	conn, err := gtcp.NewConn("127.0.0.1:8300")
	if err xxx
	defer conn.Close()
	if result, err := conn.SendRecvPkg(info); err != nil {
		if err.Error() == "EOF" {
			xxx
		}
	}
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
