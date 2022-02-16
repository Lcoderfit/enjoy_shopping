package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-时间维护
https://goframe.org/pages/viewpage.action?pageId=1114139

一、对数据创建，更新，删除时间的自动填充(该特性仅对链式操作有效)
	1.1 字段名规定为create_at, update_at, delete_at,也可以通过配置文件进行自定义配置(不过一般使用规定的名称)
	1.2 一旦表中包含create_at, update_at, delete_at 则该特性自动开启
	1.3 Replace方法也会更新create_at, 因为这个本质上是插入数据,如果存在主键或唯一索引冲突,则删除所有冲突数据然后
		再插入
	1.4 当插入一条数据时,create_at和update_at均会被填充(所以Replace也会填充update_at)
	1.5 如果启用了软删除特性,则查询语句都会自动添加一个条件: delete_at is null(即选取未被软删除的数据)
	1.6 忽略时间特性: Unscoped()
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user")
		group.GET("/time", func(r *ghttp.Request) {
			_, err := m.Update(g.Map{
				"name": "xh",
			}, "uid", 6)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("update_at: ok")
			}

			_, err = m.Replace(g.Map{
				"uid":  4,
				"name": "x4",
			})
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("replace: ok")
			}

			_, err = m.Replace(g.Map{
				"uid": 4,
				"name": "x4Unscoped",
			})
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("unscoped: ok")
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
