package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gproc"
	"time"
)

/*
1.gproc.Pid()
	1.1 返回当前进程的PID
*/

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"graceful":        true,
		"gracefulTimeout": 3,
	})
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello")
	})
	s.BindHandler("/pid", func(r *ghttp.Request) {
		r.Response.Writeln(gproc.Pid())
	})
	s.BindHandler("/sleep", func(r *ghttp.Request) {
		r.Response.Writeln(gproc.Pid())
		time.Sleep(10 * time.Second)
		r.Response.Writeln(gproc.Pid())
	})
	s.EnableAdmin()
	s.SetPort(8300)
	s.Start()

	g.Wait()
}
