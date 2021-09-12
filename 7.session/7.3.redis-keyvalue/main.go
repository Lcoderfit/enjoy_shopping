package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.redis key-value 存储模式
	1.1 SessionMaxAge: session最大存活时间
	1.2 SessionStorage: gsession.NewStorageRedis(g.Redis()) 设置使用 内存+redis的存储模式
		SessionStorage: gsession.NewStorageMemory() 设置使用 内存 存储模式
	注意：默认使用的 内存+文件 的存储模式

	s.SetConfigWithMap(g.Map{
		"SessionMaxAge": time.Minute,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})

2.与文件存储的区别
	2.1 每一次请求时如果需要对Session进行操作，都会从redis拉取一份最新的数据；而文件存储只会在内存中不存在Session时才会从文件读取
	2.2 每一次请求结束后，都会所有的Session数据进行JSON序列化，然后以key-value形式存储到redis

3.应用场景
	单个用户下Session数据量不大的场景，用 redis key-value存储模式，
	单个用户下数据量较大(>10MB)的场景，用 redis-HashTable模式


4.redis配置：https://goframe.org/pages/viewpage.action?pageId=17203977

	// 默认使用default实例的配置，default和cache表示不同实例分组，可以通过该名称来获取不同的redis实例
	[redis]
		default = "127.0.0.1:6379,db,pass"
		cache = "127.0.0.1:6379,db,pass"

*/

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Minute,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Writeln("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Writeln(r.Session.Map())
		})
		group.ALL("/del", func(r *ghttp.Request) {
			r.Session.Clear()
			r.Response.Writeln("delete ok")
		})
	})
	s.Start()

	g.Wait()
}
