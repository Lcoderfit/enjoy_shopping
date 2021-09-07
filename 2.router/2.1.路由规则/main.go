package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 1.路由匹配模式
// 	支持三种路由匹配规则
// 	字段匹配： /{name} 对URI任意位置的参数进行截取匹配
// 	命名匹配： /:name  对URI指定层级的参数进行命名匹配
// 	模糊匹配： /*name  对uri指定层级之后的部分进行模糊匹配（如果/后面无内容，仍然可以匹配成功）
// 2.r.URL.Path和r.Router.Uri的区别
// 	r.Router当前匹配的路由规则信息
// 	r.URL.Path是浏览器输入的在域名之后的那部分URL路径
// 	r.Router.Uri是自己编写代码注册的 路由匹配规则
// 3.路由检索算法
// 	采用深度优先算法，层次越深的规则优先级越高；例如：
// 	/:name/*any 会覆盖 /:name规则
// 4.路由规则格式： [HTTPMETHOD:]路由规则[@域名]
// 5.精准匹配规则
//
// 6.获取路由参数(GetRouterMap不需要传入参数，其他方法都有两个参数，一个是路由参数的名称，另一个是默认值，当路由参数不存在时即返回默认值)
// r.GetRouterMap
// r.Get
// r.GetRouterString
// r.GetRouterVar
// r.GetRouterValue
//
// 7.r.Response.Writeln("hello world") // 写入一行内容
//
// 8.路由匹配优先级
// 		层级越深优先级越高
// 		同一层级下，精准匹配优先级>三种模糊匹配
// 		同一层级下，字段匹配>命名匹配>模糊匹配 （{} > : > *）

func main() {
	s := g.Server()
	// 支持三种路由匹配规则
	// 命名匹配： /:name	匹配/之后的一段路由
	// 模糊匹配： /*name	匹配/之后的任意部分，例如/a /a/b  /a.go..... 如果/后什么都没有也可以匹配
	// 字段匹配： /{name} 对路由字符串进行截取，/go-{target}/ 可以匹配/go-to-bytedance； target匹配到go-bytedance的值
	//                   URL                               结果
	// http://127.0.0.1:8199/user/list/2.html      /user/list/{field}.html
	// http://127.0.0.1:8199/user/update           /:name/update
	// http://127.0.0.1:8199/user/info             /:name/:action
	// http://127.0.0.1:8199/user                  /:name/*any
	s.BindHandler("/:name", func(r *ghttp.Request) {
		r.Response.Write(r.URL.Path + "\t" + r.Router.Uri)
	})
	s.BindHandler("/:name/update", func(r *ghttp.Request) {
		r.Response.Write(r.URL.Path + "\t" + r.Router.Uri)

	})
	s.BindHandler("/:name/:action", func(r *ghttp.Request) {
		r.Response.Write(r.URL.Path + "\t" + r.Router.Uri)

	})
	s.BindHandler("/:name/*any", func(r *ghttp.Request) {
		r.Response.Write(r.URL.Path + "\t" + r.Router.Uri)
	})
	s.BindHandler("/user/list/{field}.html", func(r *ghttp.Request) {
		r.Response.Write(r.URL.Path + "\t" + r.Router.Uri)
	})

	// 该路由仅在GET请求下有效，路由格式： [HTTPMETHOD:]路由规则[@域名]
	// r.Router当前匹配的路由规则信息
	s.BindHandler("GET:/{table}/list/{page}.html", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
	})

	// 该路由规则仅在GET请求及localhost域名下有效
	s.BindHandler("GET:/order/info/{order_id}@localhost", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
	})

	// 该路由仅在DELETE请求下有效
	// curl -XDELETE http://127.0.0.1:8199/comment/1000
	// {"Domain":"default","Method":"DELETE","Priority":2,"Uri":"/comment/{id}"}
	s.BindHandler("DELETE:/comment/{id}", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
		r.Response.Write("\n")
		// 以键值对的形式返回所有路由参数
		r.Response.Write(r.GetRouterMap())
		r.Response.Write("\n")
		// 返回id参数的值
		r.Response.Write(r.Get("id"))
		r.Response.Write("\n")
		r.Response.Write(r.GetRouterString("id"))
		r.Response.Write("\n")
		r.Response.Write(r.GetRouterValue("id"))
		r.Response.Write("\n")
		r.Response.Write(r.GetRouterVar("id"))
	})

	s.SetPort(8199)
	s.Start()

	s1 := g.Server("s1")
	s1.SetPort(8200)
	// Writeln写入一行内容
	// curl -XGET http://127.0.0.1:8200/user/list/1.html
	// 1
	s1.BindHandler("/user/list/{page}.html", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("page"))
	})
	// 字段参数规则 和 命名参数规则混合使用
	// curl -XGET http://127.0.0.1:8200/user/info/save.php
	// user
	// attr
	// act
	s1.BindHandler("/{object}/:attr/{act}.php", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("object"))
		r.Response.Writeln("attr")
		r.Response.Writeln("act")
	})
	// 多种模糊匹配规则混用
	// curl -XGET http://127.0.0.1:8200/class3-math/john/score
	// class3
	// math
	// name
	// act
	s1.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request) {
		r.Response.Writeln(r.Get("class"))
		r.Response.Writeln(r.Get("course"))
		r.Response.Writeln("name")
		r.Response.Writeln("act")
	})
	s1.Start()

	g.Wait()
}
