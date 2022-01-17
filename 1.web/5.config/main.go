package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
配置注意事项
[test]
	# 这里必须要有上一级的配置块名称test,否则无法获取
	[test.t1]
		m = 1

# 程序获取时只能需要通过最内层的配置块名称获取, test.test.t1是获取不到的
m := g.Cfg().GetInt("test.t1.m")
 */

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			m := g.Cfg().Get("test.t1.m")
			g.Log().Line(true).Info(m)

			r.Response.WriteJson("hello world")
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}