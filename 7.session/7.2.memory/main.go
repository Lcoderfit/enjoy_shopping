package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.内存储存(基于StorageMemeory对象实现，性能高效（基于内存实现），但是没有持久化存储功能，应用程序重启之后便会丢失Session数据)
*/

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Minute,
		"SessionStorage": gsession.NewStorageMemory(),
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Writeln(r.Session.Map())
		})
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Writeln("ok")
		})
		group.ALL("/del", func(r *ghttp.Request) {
			r.Session.Clear()
			r.Response.Writeln("ok")
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}

/*
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Minute,
		"SessionStorage": gsession.NewStorageMemory(),
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Write("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Write(r.Session.Map())
		})
		group.ALL("/del", func(r *ghttp.Request) {
			r.Session.Clear()
			r.Response.Write("ok")
		})
	})
	s.SetPort(8199)
	s.Run()
*/
