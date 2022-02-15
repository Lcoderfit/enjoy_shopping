package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"time"
)

/*ORM链式操作-查询缓存
https://goframe.org/pages/viewpage.action?pageId=1114346

一、gdb.CacheOption
	// Duration < 0 则清空缓存
	// Duration == 0 永久保留
	// Duration > 0 缓存时间
	// Name 设置缓存的键名
	// Force 为true表示无论查询结果是否为nil，都会缓存
	type CacheOption struct {
		// Duration is the TTL for the cache.
		// If the parameter `Duration` < 0, which means it clear the cache with given `Name`.
		// If the parameter `Duration` = 0, which means it never expires.
		// If the parameter `Duration` > 0, which means it expires after `Duration`.
		Duration time.Duration

		// Name is an optional unique name for the cache.
		// The Name is used to bind a name to the cache, which means you can later control the cache
		// like changing the `duration` or clearing the cache with specified Name.
		Name string

		// Force caches the query result whatever the result is nil or not.
		// It is used to avoid Cache Penetration.
		Force bool
	}

二、缓存对象
	g.DB().GetCache() *gcache.Cache 返回缓存对象，可以通过该缓存对象缓存各种数据

三、缓存适配
	1.默认使用*gache.Cache缓存对象提供单进程内存缓存(如果服务采用多节点部署，会导致数据不一致的情况，所以只适用于单点)
	2.分布式redis缓存，适用于服务多节点部署的情况

四、实例
	// NOTE: 如果是链式安全操作,则每次使用缓存时也要加上Cache()方法,否则是不会使用缓存的
	m.Cache(time.Hour, "k1").Where("uid=1").One()

	// 删除缓存
	m.Cache(-1, "k1").Update(g.Map{"name": "xxm"}, "uid", 3)
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/cache", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/", func(r *ghttp.Request) {
			// 这里如果不需要对k1进行主动删除操作,也可以不设置k1键名
			res, err := m.Cache(time.Hour, "k1").Where("uid=1").One()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("set-cache-one: ", res)
			}

			// 由于上面缓存了该查询语句,所以这里不会再执行sql查询,直接使用缓存,可以看到DEBUG信息中只有一条sql打印出来(说明这一次没有查询sql)
			// NOTE: 如果是链式安全操作,则每次使用缓存时也要加上Cache()方法,否则是不会使用缓存的
			res, err = m.Cache(time.Hour, "k1").Where("uid=1").One()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("use cache: ", res)
			}

			// 更新数据,删除k1缓存
			_, err = m.Cache(-1, "k1").Update(g.Map{"name": "xxm"}, "uid", 3)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("cache: ", )
			}

			res, err = m.Cache(time.Hour, "k1").Where("uid", 1).One()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("get cache: ", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
