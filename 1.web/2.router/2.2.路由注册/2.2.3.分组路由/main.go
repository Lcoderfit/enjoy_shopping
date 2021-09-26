package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

/*
1.分组路由
	1.1 s.Group("/api")或者s.Domain("localhost").Group("api") 创建一个*ghttp.RouterGroup分组路由对象（Server对象和Domain对象均可调用Group）
	1.2 然后	将HTTP请求绑定到与请求方法同名的 分组路由对象方法，
	group.ALL() 所有请求方法均会绑定到该 分组路由方法, 然后根据路由匹配，匹配到则调用 处理器/函数/对象
	九种HTTP方法对应的group方法
	1.3 统一分组路由下注册的其他路由将拥有统一uri前缀
	1.4 分组路由的10中对象方法（GET POST....ALL），其功能类似于BindObject, 也是有三个参数；
		第一个参数为路由参数，第二个参数可以传入函数/对象/对象方法， 但是BindObject第二个参数只能传入对象;
		第三个参数与BindObject一样，都是指定注册的对象方法，如果不传入，则对象的所有可导出方法都会被注册（对象的方法名构成最后一段路由）
		传入空字符会报错：invalid method name
		指定多个对象方法时，用英文逗号分隔： group.GET("/test", c, "Hello,F1")

2.层级注册
	2.1 中间件的执行顺序(多中间件注册执行顺序)
		分前置中间件和后置中间件，前置中间件：先执行中间件处理逻辑，然后r.Middleware.Next()执行下一个流程
		后置中间件：先执行r.Middleware.Next()执行下一流程，等处理完路由服务函数后再返回来执行 中间件处理逻辑
		注意： 无论前置还是后置，都是先进行路由匹配，匹配完之后，先执行前置中间件，然后执行路由服务函数，之后执行后置中间件;
		因为未匹配的路由会返回Not Found，而匹配上的路由，在执行前置中间件时，例如执行MiddlewareAuth，则会返回Forbidden

	2.2 s.Use  group.Middleware
		2.2.1 中间件都必须是 func(r *ghttp.Request)类型的函数
		2.2.2 s.Use定义全局中间件: 对动态请求拦截有效，无法拦截静态文件请求，可以传入多个中间件，路由匹配规则是模糊匹配，相当于/\*
		2.2.3 group.Middleware定义分组路由中间件: 绑定到当前分组路由下的所有路由，无法独立使用（与全局中间件的区别），必须在分组路由中使用
		2.2.4 分组路由中间件是绑定到分组路由上的服务方法，所以不存在路由规则匹配，分组路由中间件按照注册的先后顺序执行
		2.2.5 全局中间件遵循路由匹配规则？？？？？

	2.3 Hook函数(HookBeforeServe)？？？？？？

	2.4 r.Get()当请求参数与路由参数重复时，后一个参数会覆盖前一个参数，例如路由设置为/:name，HTTP请求 http://localhost/lcoder?name=1
	时，:name参数原本可以匹配到lcoder，但是会被后面的url参数name=1覆盖，最后r.Get("name")的值为1

3.批量注册
	group.ALLMap接收map[string]interface{}类型作为参数，interface{}可以为对象（如果为对象则会注册该对象所有的可导出方法），也可以为对象方法
	group.ALLMap(g.Map{})一般配置g.Map使用（g.Map为内置的map[string]interface{}类型）, 但是ALLMap中注册的路由不能指定 接受特定HTTP请求方法
*/

func MiddlewareAuth(r *ghttp.Request) {
	g.Log().Line(true).Println("log-auth")
	token := r.Get("token")
	if token == "123456" {
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareCORS(r *ghttp.Request) {
	g.Log().Line(true).Println("log-cors")
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareLog(r *ghttp.Request) {
	g.Log().Line(true).Println("log-log")
	r.Middleware.Next()
	g.Log().Println(r.Response.Status, r.URL.Path)
}

type Controller struct{}

func (c *Controller) Hello(r *ghttp.Request) {
	r.Response.Writeln("Hello")
}

func (c *Controller) F1(r *ghttp.Request) {
	r.Response.Writeln("F1")
}

func (c *Controller) F2(r *ghttp.Request) {
	r.Response.Writeln("F1")
}

func main() {
	s := g.Server()
	// 创建分组路由对象
	group := s.Group("/api")
	group.ALL("/all", func(r *ghttp.Request) {
		r.Response.Writeln("all")
	})
	group.GET("/get", func(r *ghttp.Request) {
		r.Response.Writeln("get")
	})
	group.POST("/post", func(r *ghttp.Request) {
		r.Response.Writeln("post")
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("middleware")
	// 注册中间件
	s1.Use(MiddlewareLog)

	// 分组路由
	s1.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth, MiddlewareCORS)
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		group.Group("/order", func(group *ghttp.RouterGroup) {
			group.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
			group.PUT("/update", func(r *ghttp.Request) {
				r.Response.Writeln("update")
			})
		})
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/info", func(r *ghttp.Request) {
				r.Response.Writeln("info")
			})
			group.POST("/edit", func(r *ghttp.Request) {
				r.Response.Writeln("drop")
			})
			group.DELETE("/drop", func(r *ghttp.Request) {
				r.Response.Writeln("drop")
			})
		})
		group.Group("/hook", func(group *ghttp.RouterGroup) {
			group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				r.Response.Write("hook any")
			})
			group.Hook("/:name", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				r.Response.Writeln("hook name")
			})
		})
	})
	s1.SetPort(8201)
	s1.Start()

	c := new(Controller)
	s2 := g.Server("ALLMap")
	s2.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", c, "Hello,F1")
		group.ALLMap(g.Map{
			"/t1": c.F1,
			"/t2": c.F2,
			"/t3": c,
		})
	})
	s2.SetPort(8202)
	s2.Start()

	g.Wait()
}
