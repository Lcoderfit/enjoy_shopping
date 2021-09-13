package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		//r.Response.Writeln("start")
		//time.Sleep(3 * time.Second)
		//r.Response.Writeln("end")
		r.Response.Writeln(r.GetMap())
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
