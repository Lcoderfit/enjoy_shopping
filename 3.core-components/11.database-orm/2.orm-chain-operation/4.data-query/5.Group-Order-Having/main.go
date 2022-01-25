package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-Group/Order/Having
https://goframe.org/pages/viewpage.action?pageId=17204068


一、分组排序(Group/Order)
	// SELECT COUNT(*) total,age FROM `user` GROUP BY age
	m.Fields("count(*) total, age").Group("age").All()

	// select uid, name, age from user order by name (默认升序排序)
    m.Fields("uid, name, age").Order("name").All()

二、条件过滤(having)
	// select count(*) total, age from user group by age having total > 0;
	m.Fields("count(*) total, age").Group("age").Having("total>", 0).All()

	注意: Fields()中可以使用别名和函数 => Fields("count(*) total, age")
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.Model("user").Safe()
		group.GET("/group", func(r *ghttp.Request) {
			// all-group: [{"age":2,"total":2}]
			all, err := m.Fields("count(*) total, age").Group("age").All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-group: ", err.Error())
			} else {
				r.Response.Writeln("all-group: ", all)
			}

			// all-order: [{"age":2,"name":"jsy","uid":2},{"age":2,"name":"lh","uid":1}]
			all, err = m.Fields("uid, name, age").Order("name").All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-order: ", err.Error())
			} else {
				r.Response.Writeln("all-order: ", all)
			}

			// having 分组条件过滤
			all, err = m.Fields("count(*) total, age").Group("age").
				Having("total>", 1).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("all-having: ", err.Error())
			} else {
				r.Response.Writeln("all-having: ", all)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
