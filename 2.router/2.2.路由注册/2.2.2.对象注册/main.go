package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
对象注册
1.BindObject（面试的时候以举例子的方式说明使用的方式）
	1.1
	被注册的对象方法必须为 func(*ghttp.Request)的形式
	Index是一个特殊的方法，当注册路由规则为/object时候，HTTP请求到/object/index /object都将映射到Index方法

	1.2 路由内置变量:
	{.struct} 表示路由注册时当前注册的对象名
	{.method} 表示路由注册时当前注册的方法名


2.BindObjectMethod
3.BindObjectRest
*/

// 对象注册
type Controller struct{}

func (c *Controller) Index(r *ghttp.Request) {
	r.Response.Writeln("index")
}

func (c *Controller) Show(r *ghttp.Request) {
	r.Response.Writeln("show")
}

// 路由内置变量
type Order struct{}

func (o *Order) List(r *ghttp.Request) {
	r.Response.Writeln("list")
}

func main() {
	s := g.Server()

	// 1.对象注册
	c := new(Controller)
	s.BindObject("/object", c)
	// 2.路由内置变量, 当HTTP请求/order-list时，则会调用Order结构体的List方法
	o := new(Order)
	s.BindObject("/{.struct}-{.method}", o)

	//

	s.SetPort(8199)
	s.Run()
}
