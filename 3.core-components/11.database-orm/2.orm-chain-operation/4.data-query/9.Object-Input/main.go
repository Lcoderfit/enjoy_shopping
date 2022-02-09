package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-对象输入

Data 支持 map[string]string/map[string]interface{}/map[interface{}]interface{}
		 struct/\*struct/[]struct/[]\*struct (这个'\'可以忽略，是为了转译\*使用的)
Where/Wherexxx/Or/And 支持任意的 string/map/slice/struct/\*struct

1.当参数为struct/\*struct类型时，会自动解析为map类型，只有可导出字段能被解析
2.支持orm/gconv/json标签， 设置的是 结构体字段与表字段对应关系，如果同时有orm和json，则orm优先级更高
	NOTE：一般用orm设置结构体字段到表字段的映射，json用来实现结构体字段到json数据字段的映射
	NOTE: 无论什么标签，统一使用下划线格式

			user = &User{
				Uid:  1,
				Name: "lh",
			}
			//  SELECT * FROM `user` WHERE `uid`=1 AND `name`='lh'
			res, err = m.Where(user).All()

			user := &User{}
			// 这样使用Data是没有效果的,相当于 select * from user (Data没起到任何作用)
			res, err := m.Data(user).All()
*/

type User struct {
	Uid  int64  `orm:"uid"`
	Name string `orm:"name"`
}

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.Model("user").Safe()
		group.GET("/data", func(r *ghttp.Request) {
			user := &User{}
			// 这样使用Data是没有效果的,相当于 select * from user (Data没起到任何作用)
			res, err := m.Data(user).All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("data: ", res)
			}

			user = &User{
				Uid:  1,
				Name: "lh",
			}
			//  SELECT * FROM `user` WHERE `uid`=1 AND `name`='lh'
			res, err = m.Where(user).All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("data2:", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
