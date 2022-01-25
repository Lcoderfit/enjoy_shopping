package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-LeftJoin/RightJoin/InnerJoin
https://goframe.org/pages/viewpage.action?pageId=17204066

1.并不推荐使用Join进行联表查询，特别是在数据量比较大、并发请求量比较高的场景中，容易产生性能问题，也容易提高维护的复杂度。
2.建议: 数据库只负责存储数据和简单的单表操作，通过ORM提供的功能实现数据聚合。

一、LeftJoin/RightJoin/InnerJoin
	1.g.DB().Model("user u").LeftJoin("ads a", "u.uid=a.ad_id").Fields("u.*, a.*").One()
	2.g.DB().Model("user", "u").RightJoin("ads", "a", "u.uid=a.ad_id").Fields("u.name").Array()
	3.g.DB().Model("user u, ads a").Fields("a.ad_id").Group("ad_id").Order("timestamp").All()

设置别名的方式:
	g.DB().Model("user u")
	g.DB().Model("user", "u")
	g.DB().Model("user u, ads a")

Join查询参数
	LeftJoin("ads a", "u.uid=a.ad_id")  第二个参数是on条件
	LeftJoin("ads", "a", "u.uid=a.ad_id") 如果只有两个参数,则第二个参数是on查询条件
										如果有三个参数,则第二个是别名,第三个是on查询条件


二、自定义数据表别名--(通过 Model("user").As("u") 方法设置 user表的别名)
	// SELECT * FROM `user` AS u LEFT JOIN `user_detail` as ud ON(ud.id=u.id) WHERE u.id=1 LIMIT 1
	db.Model("user", "u").LeftJoin("user_detail", "ud", "ud.id=u.id").Where("u.id", 1).One()
	db.Model("user").As("u").LeftJoin("user_detail", "ud", "ud.id=u.id").Where("u.id", 1).One()

*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/join", func(r *ghttp.Request) {
			// left join
			// select * from user u left join ads a on u.uid=a.ad_id limit 1
			one, err := g.DB().Model("user u").LeftJoin("ads a", "u.uid=a.ad_id").
				Fields("u.*, a.*").One()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("One: ", err.Error())
			} else {
				r.Response.Writeln("one: ", one)
			}

			// right join
			// select a.ad_id from user u right join ads on u.uid=a.ad_id
			value, err := g.DB().Model("user", "u").RightJoin("ads", "a", "u.uid=a.ad_id").
				Fields("u.name").Array()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("value: ", err.Error())
			} else {
				r.Response.Writeln("value: ", value)
			}

			// inner join
			// SELECT a.ad_id FROM `user` u,`ads` a GROUP BY `ad_id` ORDER BY `timestamp`
			all, err := g.DB().Model("user u, ads a").Fields("a.ad_id").
				Group("ad_id").Order("timestamp").All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all: ", err.Error())
			} else {
				r.Response.Writeln("all: ", all)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
