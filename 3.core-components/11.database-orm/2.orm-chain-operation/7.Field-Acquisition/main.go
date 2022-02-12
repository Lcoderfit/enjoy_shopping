package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-字段获取
acquisition  获得 美[ˌækwɪˈzɪʃn]

一、FieldsStr/FieldsExStr字段获取
	NOTE：获取字段名，等效于sql中的: show columns from xxx
	1.1 FieldsStr(prefix) 获取所有字段,并添加在字段名前面添加prefix前缀
	1.2 FieldsExStr(fields, prefix) 获取除fields之外的字段,并添加prefix前缀
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/fields-str", func(r *ghttp.Request){
			// SHOW FULL COLUMNS FROM `user`
		    res := m.GetFieldsStr()
		    r.Response.Writeln("fields-str: ", res)

		    res = m.GetFieldsStr("prefix_")
		    r.Response.Writeln("fields-str-prefix: ", res)

			res = m.GetFieldsExStr("uid,num", "pre_")
			r.Response.Writeln("fields-ex-str:", res)
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}