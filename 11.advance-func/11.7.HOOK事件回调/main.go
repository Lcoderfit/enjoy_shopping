package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

/*
1.Hook的四种类型
	1.1 ghttp.HookBeforeServer
		在中间件和服务函数之前执行
	1.2 ghttp.HookAfterServer
		在中间件和服务函数之后执行，但是在ghttp.HookBeforeOutput之前执行
	1.3 ghttp.HookBeforeOutput
		在ghttp.HookAfterServer之后执行，但是在响应内容输出到客户端之前执行
	1.4 ghttp.HookAfterOutput
		在响应输出到客户端之后执行

	注意：ghttp.HookBeforeOutput的原理，如果在HookBeforeOutput之前的服务函数中有r.Response.Writexx语句，则会先将内容写入到输出缓存中，
	且如果ghttp.HookBeforeOutput中存在r.Response.Writexx语句，则也会存入输出缓存，等HookBeforeOutput运行完之后才会将输出缓存中内容
	输出到客户端

 */

func main() {
	p := "/:name/info/{uid}"
	s := g.Server()
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) { glog.Println(ghttp.HookBeforeServe) },
		ghttp.HookAfterServe:  func(r *ghttp.Request) { glog.Println(ghttp.HookAfterServe) },
		ghttp.HookBeforeOutput: func(r *ghttp.Request) {
			glog.Println(ghttp.HookBeforeOutput)
			r.Response.Writeln(ghttp.HookBeforeOutput)
		},
		ghttp.HookAfterOutput: func(r *ghttp.Request) {
			glog.Println(ghttp.HookAfterOutput)
			r.Response.Writeln(ghttp.HookAfterOutput)
		},
	})
	s.BindHandler(p, func(r *ghttp.Request) {
		g.Log().Line(false).Println("用户:", r.Get("name"), ", uid", r.Get("uid"))
		r.Response.Writeln("用户:", r.Get("name"), ", uid", r.Get("uid"))
	})
	s.SetPort(8300)
	s.Start()

	g.Wait()
}
