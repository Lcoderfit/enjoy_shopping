package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"time"
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
		如果在ghttp.HookAfterOutput中存在r.Response.Writexxx语句，则将失效，因为HookAfterOutput是在响应内容输出到客户端之后执行的

2.注册相同路由的Hook(相同事件注册)
	如果多个Hook函数注册的路由不同，则会按照路由匹配优先级顺序执行相应的Hook函数
	如果多个Hook函数注册的路由相同，则会按照注册的优先顺序执行，

3.ghttp.HandlerFunc
	是func(r *ghttp.Request)的匿名函数变量(HandlerFunc = func(r *ghttp.Request))

4.BindHookHandler和BindHookHandlerByMap group.Hook
	4.1 BindHookHandler(pattern, hook string, ghttp.HandlerFunc) 例如:
		s.BindHookHandler("/", ghttp.HookBeforeServer, func(r *ghttp.Request){xxx})
	4.2 BindHookHandlerByMap(pattern, map[string]ghttp.HandlerFunc)
		注意：BindHookHandlerByMap由于第二个参数是一个map，所以key必须唯一，即每种类型的hook只能注册一个
		s.BindHookHandlerByMap("/", map[string]ghttp.HandlerFunc{
			ghttp.HookBeforeServer: func(r *ghttp.Request){xxx},
			ghttp.HookAfterServer: func(r *ghttp.Request){xxx},
			ghttp.HookBeforeOutput: func(r *ghttp.Request){xxx},
			ghttp.HookAfterOutput: func(r *ghttp.Request){xxx},
		})
	4.3 group也有自带的Hook函数
		group.Hook("/", ghttp.HookBeforeServer, func(r *ghttp.Request){xxx})

5.
r.SetParam是设置路由参数，具有最高优先级
r.Response.SetBuffer会先将输出缓冲区清空，然后写入新的内容，参数为[]byte类型

6.路由匹配优先级
	精准匹配>三种模糊匹配
	字段匹配 > 命名匹配 > 模糊匹配

7.中间件和Hook(HookBeforeServer)均能处理跨域，一般使用中间件居多

8.ExitHook
	8.1 r.ExitHook() 只会退出当前种类Hook后续的其他Hook流程，对其他种类的Hook不影响
		例如HookBeforeServer有3个，hbs1， hbs2， hbs3
		HookAfterServer有2个，has1, has2
		假如在hbs2和has1中添加r.ExitHook()，则执行顺序hbs1 -> hbs2 -> has1，
		r.ExitHook()所在的Hook函数属于那种类型，则同属于该类型的后续的其他Hook均不会执行，但不同类的hook不受影响

	8.2 r.ExitAll()
		在Hook中调用r.ExitAll(),后续所有流程均不会执行(注意，前提是在Hook中不调用r.Middleware.Next())
		(不仅是在HookBeforeServer中是这样，在其他Hook中也是这样，例如HookAfterServer，服务函数和中间件都执行完了才会执行
		HookAfterServer，而在者后续的流程都是线性的（中间件中因为有r.Middleware.Next(),所以是栈式的）),所以r.ExitAll()
		在Hook中调用会使后续所有流程退出

9.Hook事件回调应用场景
	9.1 一般用于接口权限控制，在权限校验的事件回调函数中执行r.Redirect*方法，
		又没有调用r.ExitAll()方法退出业务执行，往往会产生http multiple response writeheader calls错误提示。
		因为r.Redirect*方法会往返回的header中写入Location头；而随后的业务服务接口往往会往header
		写入Content-Type/Content-Length头，这两者有冲突造成的。

10.r.Router.Uri和r.URL.Path的区别
	10.1 r.Router.Uri是注册的服务回调方法的路由 r.URL.Path是HTTP请求的URL
	10.2 注意：由于HookBeforeServer是在服务回调方法之前执行的，所以此时还没有匹配到服务回调函数路由，request.Router对象是nil的
			只有在匹配到服务回调方法之后，在对应的服务回调方法内r.Router对象才是有值的

11.中间件与Hook的区别
	11.1 Hook作用的范围更大，例如HookBeforeServer是在整个服务流程之前执行;而中间件可仅作用于部分服务(也可用全局)
	11.2 Hook采用的是对特定事件的钩子触发机制，中间件采用洋葱模型(处理过程像洋葱外皮一样，从外到内，在从内到外)
	11.3 Hook灵活性较差，只能对整个服务进行处理，中间件可以作用于部分服务(灵活性更高)

*/

func beforeServerHook1(r *ghttp.Request) {
	r.SetParam("name", "GoFrame")
	r.Response.Writeln("set name")
}

func beforeServerHook2(r *ghttp.Request) {
	r.SetParam("site", "https://goframe.org")
	r.Response.Writeln("set site")
}

func main() {
	p := "/:name/info/{uid}"
	s := g.Server()
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) { glog.Println(ghttp.HookBeforeServe) },
		ghttp.HookAfterServe:  func(r *ghttp.Request) { glog.Println(ghttp.HookAfterServe) },
		ghttp.HookBeforeOutput: func(r *ghttp.Request) {
			// 这里可以证明HookBeforeOutput是在响应输出到客户端之前执行的
			// 执行流程：1.调用HookBeforeOutput之后，先执行r.Response.Writexx将响应内容写入输出缓存
			// 2.休眠10秒，（如果HookBeforeOutput不是在响应输出到客户端之前执行，则这里响应应该会先显示到浏览器页面，
			// 过10s后才在终端输出日志信息，但是实际上 响应输出到客户端和日志信息输出到终端几乎是同时进行的）
			// 3.glog.Println将日志信息输出到终端
			// 4.之后将所有在输出缓存中的内容输出到客户端
			// 5.执行HookAfterOutput
			r.Response.Writeln(ghttp.HookBeforeOutput)
			time.Sleep(10 * time.Second)
			glog.Println(ghttp.HookBeforeOutput)
		},
		ghttp.HookAfterOutput: func(r *ghttp.Request) {
			glog.Println(ghttp.HookAfterOutput)
			// 由于HookAfterOutput是在响应内容输出到客户端之后执行的，所以这条语句不会再输出内容到客户端了
			r.Response.Writeln(ghttp.HookAfterOutput)
		},
	})
	s.BindHandler(p, func(r *ghttp.Request) {
		g.Log().Line(false).Println("用户:", r.Get("name"), ", uid", r.Get("uid"))
		r.Response.Writeln("用户:", r.Get("name"), ", uid", r.Get("uid"))
	})
	s.SetPort(8300)
	s.Start()

	s1 := g.Server("s1")
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("name"))
		r.Response.Writeln(r.Get("site"))
	})
	s1.BindHookHandler("/", ghttp.HookBeforeServe, beforeServerHook1)
	s1.BindHookHandler("/", ghttp.HookBeforeServe, beforeServerHook2)
	s1.SetPort(8301)
	s1.Start()

	s2 := g.Server("s2")
	pattern1 := "/:name/info"
	s2.BindHookHandlerByMap(pattern1, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.SetParam("uid", 1000)
		},
	})
	s2.BindHandler(pattern1, func(r *ghttp.Request) {
		r.Response.Writeln("用户：", r.Get("name"), ", uid:", r.Get("uid"))
	})

	pattern2 := "/{object}/list/{page}.go"
	s2.BindHookHandlerByMap(pattern2, map[string]ghttp.HandlerFunc{
		// 1.首先会运行r.Response.Writeln(r.Router.Uri)将注册的路由写入输出缓冲区
		// 2.之后运行HookBeforeOutput，r.Response.SetBuffer会将缓冲区中内容清空并写入新的内容
		// 3.将fmt.Sprintf的内容转换为字节数组然后写入缓冲区，最终返回到客户端
		ghttp.HookBeforeOutput: func(r *ghttp.Request) {
			// SetBuffer先将输出缓存区中的内容清空，然后根据传入的data []byte参数重新写入缓冲区
			r.Response.SetBuffer([]byte(
				fmt.Sprintf("通过事件修改输出内容, object:%s, page:%s",
					r.Get("object"), r.GetRouterString("page"))),
			)
		},
	})
	s2.BindHandler(pattern2, func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s2.SetPort(8302)
	s2.Start()

	// 先执行HookBeforeServer，优先级：精准匹配>:name>*any
	// 然后执行服务函数
	s3 := g.Server("s3")
	s3.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Writeln("prority service")
	})
	s3.BindHookHandlerByMap("/priority/:name", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Writeln("/priority/:name")
		},
	})
	s3.BindHookHandlerByMap("/priority/*any", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Writeln("/priority/*any")
		},
	})
	s3.BindHookHandlerByMap("/priority/show", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Writeln("/priority/show")
		},
	})
	s3.SetPort(8303)
	s3.Start()

	// 如果不设置跨域，默认是不允许跨域的请求的
	s4 := g.Server("s4")
	s4.BindHandler("/api.v1", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s4.SetPort(8304)
	s4.Start()

	// 在百度页面按f12，发送ajax请求
	// $.get("http://localhost:8305/api.v1/order", function(result){
	//     console.log(result)
	// })
	// r.Response.CORSDefault()设置允许所有跨域请求
	s5 := g.Server("s5")
	s5.Group("/api.v1", func(group *ghttp.RouterGroup) {
		// group.Hook
		group.Hook("/*any", ghttp.HookBeforeServe, func(r *ghttp.Request) {
			// 使用默认的跨域设置处理跨域请求
			r.Response.CORSDefault()
		})
		group.GET("/order", Order)
	})
	s5.SetPort(8305)
	s5.Start()

	s6 := g.Server("s6")
	s6.BindHandler("/name", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	// 在BindHookHandlerByMap中每种类型的hook只能注册一次
	s6.BindHookHandlerByMap("/name", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Writeln(ghttp.HookBeforeServe, "1")
			r.ExitHook()
		},
		ghttp.HookAfterServe: func(r *ghttp.Request) {
			r.Response.Writeln(ghttp.HookAfterServe, "1")
			r.ExitHook()
		},
		ghttp.HookBeforeOutput: func(r *ghttp.Request) {
			r.ExitHook()
			r.Response.Writeln(ghttp.HookBeforeOutput, "1")
		},
		ghttp.HookAfterOutput: func(r *ghttp.Request) {
			r.Response.Writeln(ghttp.HookAfterOutput, "1")
		},
	})
	s6.BindHookHandler("/:name", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Writeln(ghttp.HookBeforeServe, "2")
	})
	s6.BindHookHandler("/:name", ghttp.HookAfterServe, func(r *ghttp.Request) {
		r.Response.Writeln(ghttp.HookAfterServe, "2")
	})
	s6.BindHookHandler("/:name", ghttp.HookBeforeOutput, func(r *ghttp.Request) {
		r.Response.Writeln(ghttp.HookBeforeOutput, "2")
	})
	s6.BindHookHandler("/:name", ghttp.HookAfterOutput, func(r *ghttp.Request) {
		r.Response.Writeln(ghttp.HookAfterOutput, "2")
	})
	s6.SetPort(8306)
	s6.Start()

	g.Wait()
}

func Order(r *ghttp.Request) {
	r.Response.Writeln("GET")
}
