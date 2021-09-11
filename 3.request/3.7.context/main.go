package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.r.SetCtxVar(key, value)和r.SetCtx(ctx context.Context)的区别
	1.1 这两个方法效果上其实是等效的，前者在内部使用：r.context = context.WithValue(r.Context(), key, value)
	为r.context设置上下文变量, 而后者需要先创建一个ctx := context.WithValue(r.Context(), key, value)
	然后再将ctx赋值给r.context: r.context = ctx

2.r.Context()
	该方法返回r.context的值，所以 r.context = context.WithValue(r.Context(), key, value)其实就相当于对当前r.context的
	值拷贝一份，然后拷贝后的上下文设置key和value，再将新的上下文重新赋值给r.context (新的r.context比老的r.context多了一对key value的映射关系)

3.获取设置的上下文中的变量值
r.Context().Value(key) 直接返回接口类型的值
r.GetCtxVar(key) 其实内部调用了 value := r.Context().Value(key)，然后根据value的值创建了一个*gvar.Var类型的值返回
r.GetCtx().Value(key) r.GetCtx()内部返回r.Context(),然后调用Value就相当于r.Context().Value(key)

r.GetCtx()返回当前context.Context对象，同r.Context()
r.GetCtxVar(key, def) 获取当前上下文变量，可以设置当变量不存在时的默认值

r.SetCtx(ctx) 设置自定义context.Context对象
r.SetCtxVar(key, value) 设置上下文变量
*/

const (
	TraceIdName = "trace-id"
)

func MiddlewareCtxVar(r *ghttp.Request) {
	r.SetCtxVar(TraceIdName, "HBm876TFCde435Tgf")
	r.Middleware.Next()
}

func MiddlewareCtx(r *ghttp.Request) {
	ctx := context.WithValue(r.Context(), TraceIdName, "HBm876TFCde435Tgf")
	r.SetCtx(ctx)
	r.Middleware.Next()
}

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCtxVar)
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Writeln(r.GetCtxVar(TraceIdName))
			r.Response.Writeln(r.Context().Value(TraceIdName))
			r.Response.Writeln(r.GetCtx().Value(TraceIdName))
		})
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("ctx")
	s1.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCtx)
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Writeln(r.GetCtx().Value(TraceIdName))
			r.Response.Writeln(r.Context().Value(TraceIdName))
			r.Response.Writeln(r.GetCtxVar(TraceIdName))
		})
	})
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
