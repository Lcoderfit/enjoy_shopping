package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.redis HashTable存储模式
	对Session的增删改查都是通过直接访问redis实现
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge": time.Minute,
		"SessionStorage": gsession.NewStorageRedisHashTable(g.Redis()),
	})

2.与redis key-value模式的区别
	redis key-value模式每次请求时会拉取最新的数据到内存中，然后请求完后再将内存中全量数据更新到redis
	redis HashTable模式，对Session的增删该查是直接操作redis实现的，如果中断服务后重启，数据会从redis中恢复;
						如果手动修改redis的键值对数据，则页面访问时也会刷新

3.redis查看key的类型： type key_name
	redis查询hash类型不能使用get方法：hget key field 或 hgetall key
	如果使用get方法会报错：WRONGTYPE Operation against a key holding the wrong kind of value
	修改redis hash的字段值：hset "ppu7ev015ccdk2ce752lkvr2cc100279" name john

4. 注意：只有file存储模式下，你读取Session变量时才会自动重置ttl时间，而其他模式不会重置，
		也就是说一旦设置了session的存活时间（除非重新设置session），无论读取多少次都不会重置ttl
*/

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  70 * time.Second,
		"SessionStorage": gsession.NewStorageRedisHashTable(g.Redis()),
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("name", gtime.Timestamp())
			r.Response.Writeln("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Writeln(r.Session.Map())
		})
		group.ALL("/del", func(r *ghttp.Request) {
			r.Session.Clear()
			r.Response.Writeln("ok")
		})
	})
	s.Start()

	g.Wait()
}
