package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
Hook事件回调

1.r.Exit()和r.Response.WritexxxExit
	WritexxExit的内部其实是Writexx之后调用了r.Exit退出当前请求流程

2.r.Exit()和r.ExitAll()
	2.1 正常情况下的请求流程
		request->静态文件检索->HookBeforeServer->中间件->服务函数->中间件->HookAfterServer
		->HookBeforeOutput->响应输出->HookAfterOutput

		其中如果有多个HookBeforeServer，则按照路由匹配优先级执行（例如/:name优先于/\*执行，
		在/:name的hook中调用r.Exit()，仍会执行/\*对应的hook）

		如果有多个中间件，例如：s.Use(m1, m2), 则中间件按照注册的顺序执行（m1比m2先执行）
		且中间件执行类似于栈，先注册(m1)的先执行r.Middleware.Next()前面的部分，然后再执行后注册的(m2)
		中r.Middleware.Next()的前半部分，之后执行完服务函数后，再执行m2中的r.Middleware.Next()的后半部分，
		最后执行m1中的r.Middleware.Next()的后半部分

	2.2
	r.Exit()只是退出当前请求流程（退出当前hook函数/中间件函数/服务函数），
	但是后续如果还有其他 hook函数/中间件/服务函数 流程仍会执行

	2.3 r.ExitAll()
		2.3.1 如果在HookBeforeServer中执行，则后续所有流程均不会执行（
			包括中间件也不会执行，因为HookBeforeServer的执行流程在中间件之前）
		2.3.2 如果在中间件中执行，则分两种情况：
			2.3.2.1 在r.Middleware.Next()之前执行r.ExitAll():
						则所有后续流程均不执行（包括后面注册的中间件，服务函数，HookAfterxxx）
			2.3.2.2 在r.Middleware.Next()之后执行r.ExitAll():
						则退出当前中间件流程，都之前注册的中间件的r.Middleware.Next()之后的流程不产生影响，
						但是HookAfterxxx流程将不再执行
		2.3.3 在服务函数中执行
			对中间件流程不产生影响，但是会退出HookAfterxx流程

*/

func MiddlewareTest(r *ghttp.Request) {
	g.Log().Line(true).Println("begin middleware test")
	r.Middleware.Next()
	g.Log().Line(true).Println("before middleware ExitAll")
	//r.ExitAll()
	g.Log().Line(true).Println("after middleware ExitAll")
}

func MiddlewareTest1(r *ghttp.Request) {
	g.Log().Line(true).Println("begin middleware t1")
	r.Middleware.Next()
	g.Log().Line(true).Println("before middleware-t1 ExitAll")
	//r.ExitAll()
	g.Log().Line(true).Println("after middleware-t1 ExitAll")
}

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.Use(MiddlewareTest, MiddlewareTest1)
	s.BindHandler("/a", func(r *ghttp.Request) {
		if r.GetInt("type") == 1 {
			r.Response.Writeln("smith")
			r.Exit()
		}
		//r.Response.Writeln("john")
		r.Response.Writeln("john")
		r.ExitAll()
	})
	s.BindHookHandler("/:name", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Writeln("before-server:hook name")
		r.Exit()
		g.Log().Line(true).Println("before-server:hook name")
	})
	s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Writeln("before-server:hook any")
		//r.ExitAll()
		g.Log().Line(true).Println("before-server:hook any")
	})

	s.BindHookHandler("/:name", ghttp.HookAfterServe, func(r *ghttp.Request) {
		r.Response.Writeln("after-server:hook name")
		g.Log().Line(true).Println("after-server: hook name")
	})
	s.BindHookHandler("/*", ghttp.HookAfterServe, func(r *ghttp.Request) {
		r.Response.Writeln("after-server: hook any")
		g.Log().Line(true).Println("after-server:hook any")
	})
	s.Start()

	g.Wait()
}
