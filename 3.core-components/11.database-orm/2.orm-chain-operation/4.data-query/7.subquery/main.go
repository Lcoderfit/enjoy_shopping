package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-子查询特性
https://goframe.org/pages/viewpage.action?pageId=17204086

注意：Where/Having/From子查询， 传入的参数都是*Model类型

一、Where子查询-Where中的第二个参数传入*Model类型
	// select * from user where name in (select name from other)
	g.Model("user").Where("name", g.Model("other").Fields("name"))

二、Having子查询-Having中第二个参数传入*Model类型
	// select avg(num) av, name from user group by name having av>= (select avg(num) from user);
	g.Model("user").Fields("avg(num) av, name").Group("name").Having("av>=?", g.Model("user").Fields("avg(num)"))

三、From
	// SELECT * FROM (SELECT `name` FROM `user` WHERE uid<=1) as u
	// 注意:派生表必须取别名,否则会报错: Error 1248: Every derived table must have its own alias,
	g.Model("? as u", g.Model("user").Fields("name").Where("uid<=?", 1)).All()

NOTE: 1.所有的子查询条件,都要写成占位符的形式: Where("age>", *Model)是错误的 => Where("age>?", *Model)是正确的
		Having("av>=?", *Model)
		g.Model("? as u", *Model) 1.取别名 2.占位符?

      2.如果是非子查询条件,则可以不带占位符
		Where("age>", 0)
		Having("av>", 0)
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/where", func(r *ghttp.Request) {
			// where 子查询
			//  SELECT * FROM `user` WHERE num>=(SELECT max(num) FROM `user`)
			all, err := g.Model("user").Where(
				"num>=?",
				g.Model("user").Fields("max(num)"),
			).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-where: ", err.Error())
			} else {
				r.Response.Writeln("all-where: ", all)
			}

			// having子查询
			// SELECT avg(num) av,name FROM `user` GROUP BY `name` HAVING av>=(SELECT avg(num) FROM `user`)
			subQuery := g.Model("user").Fields("avg(num)")
			all, err = g.Model("user").Fields("avg(num) av, name").
				Group("name").Having("av>=?", subQuery).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-group-subQuery: ", err.Error())
			} else {
				r.Response.Writeln("all-group-subQuery: ", all)
			}

			// from子查询
			// SELECT * FROM (SELECT `name` FROM `user` WHERE uid<=1) as u
			subQuery = g.Model("user").Fields("name").Where("uid<=?", 1)
			all, err = g.Model("? as u", subQuery).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-from: ", err.Error())
			} else {
				r.Response.Writeln("all-from: ", all)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
