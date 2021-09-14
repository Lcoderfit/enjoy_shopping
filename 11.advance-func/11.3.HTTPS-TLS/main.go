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
	s.EnableHTTPS("/server.crt", "/server.key")
	s.SetPort(80)
	s.SetHTTPSPort(443)
	s.Start()

	g.Wait()
}

