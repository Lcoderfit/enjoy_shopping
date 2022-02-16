package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-数据库切换
https://goframe.org/pages/viewpage.action?pageId=1114237

g.DB().Schema("godmin").Model("sys_user') 切换单例数据库对象(从lc-sql切换到godmin)
g.DB().Model("sys_user").Schema("godmin") 从default中的数据库lc-sql切换到godmin中数据库
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/database", func(group *ghttp.RouterGroup) {
		m := g.DB()
		group.GET("/", func(r *ghttp.Request) {
			// 通过g.DB().Model().Schema()切换数据库
			res, err := m.Schema("godmin").Model("sys_user").All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("user_detail: ", res)
			}

		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}