package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-Struct/Structs/Scan

1.Struct(pointer)
	将查询结果转换为一个struct对象(字段匹配则赋值,否则为类型的零值)
	pointer应该为一个\*struct/\*\*struct类型

	推荐方式(这样只在查询的时候才会执行初始化分配内存)
	var user *User
	Struct(&user)

	SELECT `uid`,`name`,`age` FROM `user` WHERE uid>0 LIMIT 1

2.Structs(pointer)
	将查询结果转换为一个struct对象切片
	pointer应该为一个\*[]\*struct 或 \*[]struct 类型

	// 推荐方式(这样只会在查询的时候执行初始化分配内存)
	var users []*User
	Structs(&users)

3.Scan(pointer)--struct和structs方法将被弃用,推荐使用Scan
	如果传入\*struct/\*\*struct, 则会调用struct方法
	如果传入*[]struct *[]*struct, 则会调用structs方法

SELECT `uid`,`name`,`age` FROM `user` WHERE uid>0
*/

type User struct {
	Uid  uint64
	Name string
	Age  int
}

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/", func(r *ghttp.Request) {
			var user *User
			// SELECT `uid`,`name`,`age` FROM `user` WHERE uid>0 LIMIT 1
			err := m.Where("uid>0").Struct(&user)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("Struct: ", err.Error())
			} else {
				r.Response.Writeln("Struct: ", user)
			}

			var users []*User
			// SELECT `uid`,`name`,`age` FROM `user` WHERE uid>0
			err = m.Where("uid>0").Structs(&users)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("structs: ", err.Error())
			} else {
				r.Response.Writeln("structs: ", users)
			}

			// SELECT `uid`,`name`,`age` FROM `user` WHERE uid>0
			err = m.Where("uid>0").Scan(&users)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("Scan: ", err.Error())
			} else {
				r.Response.Writeln("Scan: ", users)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
