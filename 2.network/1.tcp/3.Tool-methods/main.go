package main

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
)

/*
1.gtcp.SendRecv第一个地址参数需要加上端口号，否则会报错
	panic: dial tcp: address www.baidu.com: missing port in address


2.安全加密通信相关
	2.1
	func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error)
	func NewNetConn(addr string, timeout ...int) (net.Conn, error)

	2.2 创建支持TLS安全加密通信的TCP客户端
	func NewNetConnKeyCrt(addr, crtFile, keyFile string) (net.Conn, error)
	func NewNetConnTLS(addr string, tlsConfig *tls.Config) (net.Conn, error)

3. gtcp.Sendxxx系列函数，可以直接通过给指定地址发送数据，获取响应结果（一般用于短链接请求）
	func Send(addr string, data []byte, retry ...Retry) error
	func SendPkg(addr string, data []byte, option ...PkgOption) error
	func SendPkgWithTimeout(addr string, data []byte, timeout time.Duration, option ...PkgOption) error
	func SendRecv(addr string, data []byte, receive int, retry ...Retry) ([]byte, error)
	func SendRecvPkg(addr string, data []byte, option ...PkgOption) ([]byte, error)
	func SendRecvPkgWithTimeout(addr string, data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error)
	func SendRecvWithTimeout(addr string, data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error)
	func SendWithTimeout(addr string, data []byte, timeout time.Duration, retry ...Retry) error

*/

func main() {
	data, err := gtcp.SendRecv("www.baidu.com:80", []byte("HEAD / HTTP/1.1\n\n"), -1)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
