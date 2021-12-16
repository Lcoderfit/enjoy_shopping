package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/util/gconv"
	"time"
)

/*
1.certificate signed by unknown authority
	客户端会对证书进行校验，通过
	tlsConfig, err := gtcp.LoadKeyCrt(crtFile, keyFile)
	xxxx
	// 跳过客户端对证书的校验
	tlsConfig.InSecureSkipVerify = true
	conn, err := gtcp.NewConnTLS(address, tlsConfig)

2.gtcp.NewServerKeyCrt(address, crtFile, keyFile)
	创建一个支持TLS安全加密通信的tcp服务

3.gtcp.LoadKeyCrt(crtFile, keyFile)
	根据指定的验证文件和密钥文件创建一个TLS配置对象

4.gtcp.NewConnKeyCrt(address, crtFile, keyFile)
	根据指定的crtFile和keyFile创建一个TLS连接

5.gtcp.NewConnTLS(address, tlsConfig)
	根据指定的TLS配置对象创建一个TLS连接

*/

func main() {
	// 创建一个支持TLS安全加密通信的TCP服务端
	address := "127.0.0.1:8300"
	crtFile := "server.crt"
	keyFile := "server.key"
	// TLS Server
	go gtcp.NewServerKeyCrt(address, crtFile, keyFile, func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			// 传入-1表示一次性接收tcp客户端发送的所有数据然后返回
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				fmt.Println(string(data))
			}
			if err != nil {
				g.Log().Error(err)
				break
			}
		}
	}).Run()

	time.Sleep(time.Second)

	// Client
	tlsConfig, err := gtcp.LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		panic(err)
	}
	// 跳过客户端对证书的校验
	tlsConfig.InsecureSkipVerify = true
	conn, err := gtcp.NewConnTLS(address, tlsConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//conn, err := gtcp.NewConnKeyCrt(address, crtFile, keyFile)
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()

	for i := 0; i < 10; i++ {
		if err := conn.Send([]byte(gconv.String(i))); err != nil {
			g.Log().Error(err)
		}
		time.Sleep(time.Second)
		if i == 5 {
			conn.Close()
			break
		}
	}

	time.Sleep(5 * time.Second)
}
