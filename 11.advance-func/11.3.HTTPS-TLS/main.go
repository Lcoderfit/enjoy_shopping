package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s.EnableHTTPS("server.crt", "server.key")
	s.SetPort(8300)
	s.Start()

	// 证书路径在服务器上，下面的代码需要在服务器上运行
	s1 := g.Server("s1")
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("可以同时通过HTTPS和HTTP访问")
	})
	// 证书和密钥的路径
	s1.EnableHTTPS(
		"/etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem",
		"/etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem",
	)
	s1.SetPort(8200)
	s1.SetHTTPSPort(8201)
	s1.Start()

	g.Wait()
}
