package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	// 获取一个默认的Server实例，采用单例模式（多次调用该方法返回同一个Server对象）
	s := g.Server()
	// 默认情况下不支持静态文件处理
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("哈喽世界!")
	})
	// 执行Server的监听运行，默认监听80端口
	s.Run()
}
