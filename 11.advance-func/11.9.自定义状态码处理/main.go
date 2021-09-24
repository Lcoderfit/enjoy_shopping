package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
func (s *Server) BindStatusHandler(status int, handler HandlerFunc)
func (s *Server) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc)
func (s *Domain) BindStatusHandler(status int, handler HandlerFunc)
func (s *Domain) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc)

1.针对指定状态码进行自定义处理
	s.BindStatusHandler(404, func(r *ghttp.Request){})
	例如访问一个不存在的url: /index, 本来应该显示Not Found, 可以通过该函数进行自定义处理，
	当返回404时，即调用BindStatusHandler的第二个参数注册的处理器

	s.BindStatusHandlerByMap(handlerMap map[int]HandlerFunc)
		可一次性设置多个对不同状态码的自定义处理操作

2.引导跳转到指定的错误页面
	s1.BindHandler("/status/:status", func(r *ghttp.Request) {
		r.Response.Write("status: ", r.Get("status"), " found")
	})
	s1.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.RedirectTo("/status/404")
	})
	当访问一个不存在的url时，例如/index（404 Not Found）, 会先调用BindStatusHandler注册的状态码处理函数，
	然后跳转到/status/404

3.注意事项
	如果在调用状态码处理函数之前r.Buffer中已存在其他内容，需要先通过r.Response.ClearBuffer()清空缓冲区内容

*/

func MiddleWareClearBuffer(r *ghttp.Request) {
	r.Response.ClearBuffer()
	g.Log().Println("h2")
	r.Middleware.Next()
}

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Writeln("This is customized 404 page")
	})
	s.SetPort(8300)
	s.Start()

	// 通过r.Response.RedirectTo引导跳转到错误页面
	s1 := g.Server("s1")
	s1.BindHandler("/status/:status", func(r *ghttp.Request) {
		r.Response.Write("status: ", r.Get("status"), " found")
	})
	s1.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.RedirectTo("/status/404")
	})
	s1.SetPort(8301)
	s1.Start()

	// 批量注册状态码处理函数
	s2 := g.Server("s2")

	s2.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteStatus(403, "Forbidden")
	})
	s2.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
		403: func(r *ghttp.Request) {
			ClearByStatus(r, 403)
		},
		404: func(r *ghttp.Request) {
			ClearByStatus(r, 404)
		},
		500: func(r *ghttp.Request) {
			ClearByStatus(r, 500)
		},
	})
	s2.SetPort(8302)
	s2.Start()

	g.Wait()
}

func ClearByStatus(r *ghttp.Request, status int) {
	r.Response.ClearBuffer()
	r.Response.Writeln(status)
}
