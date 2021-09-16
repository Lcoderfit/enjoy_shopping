package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("可以同时通过HTTPS和HTTP访问")
	})
	// 证书和密钥的路径
	s.EnableHTTPS(
		"/etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem",
		"/etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem",
	)
	s.SetPort(8200)
	s.SetHTTPSPort(8201)
	s.Start()

	g.Wait()
}
