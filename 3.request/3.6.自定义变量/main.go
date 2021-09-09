package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.自定义变量
	1.1
	r.Get()返回的是一个接口类型
	r.GetString() 返回一个string类型
	r.GetMap() 返回map[string]interface{}类型
	r.GetVar() 返回*gvar.Var类型，该类型又包含各种内置方法，例如r.GetVar().String()可以将值转换为string类型

	注意：除了GetMap第一个参数为设置的默认值外，其他三个方法均有第二个参数，第二个参数为设置的默认值

	1.2 自定义变量
	r.SetParam(key, value): 设置参数key的值为value, 注意：该方法设置的变量具有最高优先级，即它会覆盖所有客户端的输入参数
	设置的变量在整个请求流程中都是可以共享的

	r.GetParam(key): 获取自定义的变量参数，返回一个interface{}类型，第二个参数可以设置默认值，当key不存在时即返回设置的默认值，
	如果key不存在，且未设置默认值，则返回nil

	r.GetParamVar(key): 获取自定义的变量，返回*gvar.Var类型，第二个参数可以设置默认值，如果key不存在则返回默认值，
	如果key不存在，且未设置默认值，则返回nil，*gvar.Var类型的好处在于其包含各种与内置类型同名的方法，可以将值转换为各种内置类型

	1.3 注意：对于通过r.SetParam()设置的变量，r.Get/r.GetString/r.GetVar/r.GetMap 和 r.GetParam/r.GetParams均可获取变量值
 */

// 前置中间件1
func MiddlewareBefore1(r *ghttp.Request) {
	r.SetParam("name", "GoFrame")
	r.Response.Writeln("set name")
	r.Middleware.Next()
}

// 前置中间件2
func MiddlewareBefore2(r *ghttp.Request) {
	r.SetParam("site", "https://goframe.org")
	r.Response.Writeln("set site")
	r.Middleware.Next()
}

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareBefore1, MiddlewareBefore2)
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Writefln(
				"%s: %s",
				r.GetParamVar("name").String(),
				r.GetParamVar("site").String(),
			)
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
