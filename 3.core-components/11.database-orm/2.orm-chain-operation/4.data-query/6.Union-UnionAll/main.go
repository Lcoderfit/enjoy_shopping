package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-Union/UnionAll
https://goframe.org/pages/viewpage.action?pageId=17204084

一、Union/UnionAll方法操作(union内的参数都是*Model类型的参数)
	1.Union会去除重复数据，UnionAll不会去除重复数据
	2.直接通过g.DB().Union()进行数据合并(可以对多个表的数据进行合并,但是合并的列类型应该保持一致)
		all, err := g.DB().Union(
			g.Model("user").Where("uid", 1),
			g.Model("user").Where("uid", 2),
		).All()

	3.通过g.DB().Model("user").Union()进行合并
		// 这样做也是可以的
		all, err = g.Model("user").UnionAll(
			g.Model("user").Where("uid<=2").Fields("uid"),
			g.Model("ads").Where("ad_id <= 5").Fields("ad_id"),
		).All()

二、Dao链式操作
	dao.User.Union(
		dao.User.Where("id", g.Slice{1,2,3}),
		dao.User.Where("id", 4),
	)
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			// (select * from user where id=1) union (select * from user where id=2)
			// union的每一行数据类型应该保持一致
			// all: [{"age":2,"name":"lh","nickname":null,"num":13,"uid":1},{"age":2,"name":"jsy","nickname":null,"num":12,"uid":2}]
			// (SELECT * FROM `user` WHERE `uid`=1) UNION (SELECT * FROM `user` WHERE `uid`=2)
			all, err := g.DB().Union(
				g.Model("user").Where("uid", 1),
				g.Model("user").Where("uid", 2),
			).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all: ", err.Error())
			} else {
				r.Response.Writeln("all: ", all)
			}

			// all: [{"uid":2},{"uid":1},{"uid":1},{"uid":2},{"uid":3},{"uid":5}]
			// UnionAll 的行与行之间的相同列的数据类型应该保持一致
			//  (SELECT `uid` FROM `user` WHERE uid<=2) UNION ALL (SELECT `ad_id` FROM `ads` WHERE ad_id <= 5)
			all, err = g.Model("user").UnionAll(
				g.Model("user").Where("uid<=2").Fields("uid"),
				g.Model("ads").Where("ad_id <= 5").Fields("ad_id"),
			).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-union-fail: ", err.Error())
			} else {
				r.Response.Writeln("all-union: ", all)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
