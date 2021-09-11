package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.r.Response.RedirectTo(): 会设置location响应头
	可以是一个相对路径，也可以是一个http地址

2.r.Response.RedirectBack(): 可以设置响应状态码
	RedirectBack会获取请求的Referer头部，引导客户端跳转到请求的上一页面
	好像没啥用？？？？因为当前页面点击返回链接的话还是返回当前页面
*/

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.RedirectTo("/login")
	})
	s.BindHandler("/login", func(r *ghttp.Request) {
		r.Response.Writeln("Login First")
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("back")
	s1.BindHandler("/page", func(r *ghttp.Request) {
		r.Response.Writeln(`<a href="/back">back</a>`)
	})
	s1.BindHandler("/back", func(r *ghttp.Request) {
		r.Response.RedirectBack()
	})
	s1.SetPort(8201)
	s1.Start()

	s2 := g.Server("custom")
	s2.BindHandler("/index", func(r *ghttp.Request) {
		r.Response.Writeln(`<a href="/home">home</a>`)
	})
	s2.BindHandler("/home", func(r *ghttp.Request) {
		r.Response.Writeln(`<h1>hello world</h1> <a href="/back">back</a>`)
	})
	s2.BindHandler("/back", func(r *ghttp.Request) {
		r.Response.RedirectBack()
	})
	s2.SetPort(8202)
	s2.Start()

	g.Wait()
}

/*
func main() {
	s := g.Server()
	s.BindHandler("/page", func(r *ghttp.Request) {
		r.Response.Writeln(`<a href="/back">back</a>`)
	})
	s.BindHandler("/back", func(r *ghttp.Request) {
		r.Response.RedirectBack()
	})
	s.SetPort(8199)
	s.Run()
*/
