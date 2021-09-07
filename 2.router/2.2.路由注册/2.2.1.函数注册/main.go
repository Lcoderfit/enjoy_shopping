package main

import (
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.函数注册
2.包方法注册
3.对象方法注册(注意，不支持对象注册，只支持对象方法注册)
4.接口的并发安全性(设置难点)
*/

var (
	// 创建一个并发安全的int类型变量
	total = gtype.NewInt()
)

// 包方法：包内定义的方法
func Total(r *ghttp.Request) {
	// 每访问一次就自增1，统计访问总次数
	r.Response.Writeln("total: ", total.Add(1))
}

type Controller struct {
	total *gtype.Int
}

// 对象方法
func (c *Controller) Total(r *ghttp.Request) {
	r.Response.Writeln("total: ", c.total.Add(1))
}

func main() {
	s := g.Server()

	// 1.函数注册（直接使用匿名函数注册）
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})

	// 2.包方法注册（使用了包内定义的方法）
	s.BindHandler("/total", Total)

	// 3.对象方法注册
	c := &Controller{
		total: gtype.NewInt(),
	}
	s.BindHandler("/total1", c.Total)

	s.SetPort(8199)
	s.Start()

	// 4.不支持对象注册，只支持对象方法注册
	//s1 := g.Server("object")
	//s1.BindHandler("/total", c)
	//s1.SetPort(8200)
	//s1.Start()

	g.Wait()
}
