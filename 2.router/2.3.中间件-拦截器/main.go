package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

/*
1.通过r.Middleware属性控制请求流程
	r.Middleware.Next()跳出当前流程，执行下一流程
2.中间件类型
	前置中间件：先执行中间件逻辑，然后r.Middleware.Next()跳到下一个流程
	后置中间件：先执行r.Middleware.Next()跳到下一流程，然后再执行中间件逻辑
	中间件都必须是func(r *ghttp.Request)类型的函数

3.全局中间件和分组路由中间件
	全局中间件：通过s.Use()注册，遵循路由匹配规则？？？
		全局中间件，即使s.Use()夹在两个s.BindHandler之间，s.Use也是最先执行的
	分组路由中间件：绑定到分组路由，按照注册的顺序先后执行；执行分组路由下的所有服务接口(服务接口或服务函数)前都会先执行该中间件
		分组路由注册时，同级的路由，分组路由中间件只对中间件注册之后 注册的路由才会生效:
		group.ALL()
		// 在这之后注册的路由，MiddlewareAuth才会生效，对上面的group.ALL()是不起作用的
		group.Middleware(MiddlewareAuth)
		group.Group(xxx)


4.r.Response.CORSDefault() 设置允许所有类型的跨域请求
  r.Response.WriteStatus(http.StatusForbidden): 设置响应状态码
  r.Response.ClearBuffer() 清空响应内容
*/

func MiddlewareCORS(r *ghttp.Request) {
	// CORSDefault设置允许任意的跨域请求
	r.Response.Writeln("cors")
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token")
	if token == "123456" {
		r.Response.Writeln("auth")
		r.Middleware.Next()
	} else {
		// 禁止访问
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.Status >= http.StatusInternalServerError {
		// 清空响应buffer
		r.Response.ClearBuffer()
		r.Response.Writeln("哎哟我去，服务器居然开了小差了，请稍候再试吧")
	}
}

// 拦截处理所有请求
func MiddlewareLog(r *ghttp.Request) {
	r.Middleware.Next()
	errStr := ""
	if err := r.GetError(); err != nil {
		errStr = err.Error()
	}
	// 将响应状态码，url，错误信息输出到终端
	g.Log().Line(true).Println(r.Response.Status, r.URL.Path, errStr)
}

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS, MiddlewareAuth)
		group.ALL("/user/list", func(r *ghttp.Request) {
			r.Response.Writeln("list")
		})
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("s1")
	s1.Group("/admin", func(group *ghttp.RouterGroup) {
		group.ALL("/login", func(r *ghttp.Request) {
			r.Response.Writeln("login")
		})
		// group.Middleware(MiddlewareA) 如果在这里注册中间件A，则对上面的group.ALL不生效，对下面注册的所有路由才生效
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareAuth)
			group.ALL("/dashboard", func(r *ghttp.Request) {
				r.Response.Writeln("dashboard")
			})
		})
	})
	s1.SetPort(8201)
	s1.Start()

	s2 := g.Server("CORS")
	s2.Use(MiddlewareCORS)
	s2.Group("/api.v2", func(group *ghttp.RouterGroup) {
		// 1.如果使用 group.Middleware(MiddlewareAuth)会显示：
		// cors
		// auth
		// db error: sql is xxxxxx
		//
		// 2.如果注册Auth和ErrorHandler两个中间件，group.Middleware(MiddlewareAuth, MiddlewareErrorHandler)，显示如下：
		// 哎哟我去，服务器居然开了小差了，请稍候再试吧
		group.Middleware(MiddlewareAuth, MiddlewareErrorHandler)
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("db error: sql is xxxxxx")
		})
	})
	s2.SetPort(8202)
	s2.Start()

	s3 := g.Server("log")
	s3.Use(MiddlewareLog, MiddlewareCORS)
	s3.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth)
		group.ALL("/user/list", func(r *ghttp.Request) {
			panic("啊！我出错了！")
		})
	})
	s3.SetPort(8203)
	s3.Start()

	g.Wait()
}
