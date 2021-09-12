package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.异常堆栈信息和错误堆栈信息
	1.1 如果抛出的异常信息并不包含堆栈内容，则WebServer会自动以panic位置为基础创建一个包含堆栈信息的错误对象
	1.2 如果抛出的异常是一个gerror组件的错误对象，获取实现了堆栈打印接口的错误对象，WebServer会直接打印错误对象，不会自动创建
		gerror.Wrap(err, "UpdateData error"): 将UpdateData error作为err的上层堆栈信息

*/

func MiddlewareErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		g.Log("exception").Error(err)
		// 返回固定的友好信息
		r.Response.ClearBuffer()
		r.Response.Writeln("服务器居然开小差了，请稍后再试吧")
	}
}

func MiddlewareErrorHandler1(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		r.Response.ClearBuffer()
		r.Response.Writef("%+v", err)
	}
}

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.Use(MiddlewareErrorHandler)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("db error: sql is xxxxx")
		})
	})
	s.Start()

	s1 := g.Server("v")
	s1.Use(MiddlewareErrorHandler1)
	s1.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("db error: sql is xxxx")
		})
	})
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
