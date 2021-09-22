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


*/

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

	g.Wait()
}
