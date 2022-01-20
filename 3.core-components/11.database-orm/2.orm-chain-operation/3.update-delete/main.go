package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-更新删除
https://goframe.org/pages/viewpage.action?pageId=17203961

注意:Update()和Delete()参数都必须带有where条件,可以通过Where()函数添加条件或通过Update()/Delete()函数本身添加Where条件

一、Update
	1.1 Data()函数的三种传值方式
		Map:    g.DB().Model("user").Data(g.Map{"nickname":"ph"}).Where("name": "lh").Update()
		string: g.DB().Model("user").Data("nickname='ph'").Where("name": "lh").Update()
		字段跟值作为两个参数传入: g.DB().Model("user").Data("nickname", "ph").Where("name": "lh").Update()

	1.2 如果遇到不需要where条件的情况,直接将Where设置为1即可
		// update user set nickname='ph' where 1
		string: g.DB().Model("user").Data("nickname='ph'").Where(1).Update()

	1.3 Update()函数传参 (注意:只有两种传递update的字段跟值的方式)
		// Map可以是[string]interface [string]string .... 等等
		1.3.1 Map:    g.DB().Model("user").Update(g.Map{"name": "lh"}, "uid", "1")
		1.3.2 string: g.DB().Model("user").Update("name='lh'", "uid", "1")
 					  g.DB().Model("user").Update("name='lh'", 1) 第二个参数为1相当于where(1)

		Update只会将第一个参数识别为update的字段跟值, 从第二个参数开始会识别为Where条件(注意,这里指的是g.DB().Model("user").Update())
		如果是g.DB().Update(), 则第一个参数是表名,从第二个参数开始才相当于是g.DB().Model("user").Update()的第一个参数

		注意:以下两种方式是错误的:
			g.DB().Model("user").Update("nickname", "t2", 1)
			g.DB().Model("user").Update("nickname", "t3", "age", nil)

	1.4 order+limit限制更新的条数，注意，必须要有where条件
		// update user set age=2 where 1 order by age desc limit 2; 限制更新前两条
		m.Data("age", 2).Where(1).Order("age desc").Limit(2).Update()

二、Counter/Increment/Decrement
	2.1 字段值的增减:
			data := g.Map{
				"num": gdb.Counter{
					Field: "num",
					Value: 1,
				},
			}
			// g.DB().Model("user").Update(), 只有第一个参数传递需要update的字段跟值,第二个参数开始是where查询条件
			// UPDATE `user` SET `num`=`num`+1 WHERE `uid`=1
			_, err := m.Update(data, "uid", 1)

	2.2 Increment/Decrement
			// UPDATE `user` SET `age`=`age`+20 WHERE `uid`=3
			_, err := m.Where("uid", 3).Increment("age", 20)

			// UPDATE `user` SET `num`=`num`-10 WHERE `uid`=4
			// 注意,如果num是null,则不会发生变化
			_, err = m.Where("uid", 4).Decrement("num", 10)

三、RawSQL
	3.1 gdb.Raw()
			// UPDATE `user` SET `age`=age+1 WHERE `uid`=1
			// 将sql片段嵌入sql语句中,不自动转换为字符串类型参数
			data := g.Map{
				"age": gdb.Raw("age+1"),
			}
			_, err := m.Update(data, "uid", 1)


四、Delete
	4.1 m.Where("uid", 5).Delete()
	4.2 order+limit限制删除的条数，注意，必须要有where条件
		// 这个必须要有where条件,否则会报错:
		// there should be WHERE condition statement for DELETE operation
		// DELETE FROM `user` WHERE 1 ORDER BY `age` asc LIMIT 1
		_, err = m.Order("age asc").Limit(1).Delete(1)

	4.3 _, err = m.Delete("uid>", 3)
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.Model("user").Safe()
		// 通过Data()传递更新字段值
		group.GET("/data-update", func(r *ghttp.Request) {
			// 1.Data()接收update字段值,Where()接收where条件
			_, err := m.Data(g.Map{"name": "lh-map"}).Where("name", "lh").Update()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.WriteJson(err.Error())
			} else {
				r.Response.Writeln("map: ok")
			}

			_, err = m.Data("name='lh-string'").Where("name", "lh1").Update()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.WriteJson(err.Error())
			} else {
				r.Response.Writeln("string: ok")
			}

			// 如果Update()函数中不传递参数,则必须加上Where条件,否则会报错:
			// there should be WHERE condition statement for UPDATE operation
			_, err = m.Data("age", 2).Where(1).Order("age desc").Limit(2).Update()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.WriteJson(err.Error())
			} else {
				r.Response.Writeln("fields: ok")
			}

			_, err = m.Data("nickname", "test").Where(1).Update()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.WriteJson(err.Error())
			} else {
				r.Response.Writeln("fields-all-update: ok")
			}

		})

		// 通过Update()传递更新字段值和where条件
		group.GET("/update-where", func(r *ghttp.Request) {
			_, err := m.Update(g.Map{"name": "lh"}, "name", "lh-map")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("update-map: ok")
			}

			_, err = m.Update("name='lh1'", "name", "lh-string")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("string-update: ok")
			}

			// 不能使用m.Update("nickname", "t2", 1) 会把t2=1识别成where条件
			// Update只会将第一个参数识别为update的字段跟值, 从第二个参数开始会识别为Where条件(注意,这里指的是g.DB().Model("user").Update())
			// 如果是g.DB().Update(), 则第一个参数是表名,从第二个参数开始才相当于是g.DB().Model("user").Update()的第一个参数
			// Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server
			// version for the right syntax to use near 'WHERE `t2`=?' at line 1, UPDATE `user` SET nickname WHERE `t2`=1
			_, err = m.Update("nickname", "t2", 1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("string-update-1: ok")
			}

			// 下面语句会报错:
			// Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server
			// version for the right syntax to use near 'WHERE `t3`=?' at line 1, UPDATE `user` SET nickname WHERE `t3`='age'
			// 使用Update时不能将要更新的字段名跟字段值分开写
			_, err = m.Update("nickname", "t3", "age", nil)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("fields-update: ok")
			}
		})

		// Counter更新特性
		// 1.实现字段的自增
		// 2.实现非自身字段的新增
		group.GET("/counter", func(r *ghttp.Request) {
			data := g.Map{
				"num": gdb.Counter{
					Field: "num",
					Value: 1,
				},
			}
			// g.DB().Model("user").Update(), 只有第一个参数传递需要update的字段跟值,第二个参数开始是where查询条件
			// UPDATE `user` SET `num`=`num`+1 WHERE `uid`=1
			_, err := m.Update(data, "uid", 1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("counter-update: ok")
			}

			data = g.Map{
				"num": gdb.Counter{
					Field: "age",
					Value: 10,
				},
			}
			// g.DB().Update(), 第一个参数是表名,第二个参数传递需要update的字段跟值,从第三个参数开始是where查询条件
			// UPDATE `user` SET `num`=`age`+10 WHERE `uid`=2
			_, err = g.DB().Update("user", data, "uid", 2)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("counter-update: ok")
			}
		})

		// Increment()/Decrement()
		group.GET("/inc-dec", func(r *ghttp.Request) {
			// UPDATE `user` SET `age`=`age`+20 WHERE `uid`=3
			_, err := m.Where("uid", 3).Increment("age", 20)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("increment: ok")
			}

			// UPDATE `user` SET `num`=`num`-10 WHERE `uid`=4
			// 注意,如果num是null,则不会发生变化
			_, err = m.Where("uid", 4).Decrement("num", 10)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("dec: ok")
			}
		})

		// RawSQL
		group.GET("/raw", func(r *ghttp.Request) {
			// UPDATE `user` SET `age`=age+1 WHERE `uid`=1
			// 将sql片段嵌入sql语句中,不自动转换为字符串类型参数
			data := g.Map{
				"age": gdb.Raw("age+1"),
			}
			_, err := m.Update(data, "uid", 1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("raw update: ok")
			}
		})

		// Delete()
		group.GET("/delete", func(r *ghttp.Request) {
			_, err := m.Where("uid", 5).Delete()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("where-delete: ok")
			}

			// 这个必须要有where条件,否则会报错:
			// there should be WHERE condition statement for DELETE operation
			// DELETE FROM `user` WHERE 1 ORDER BY `age` asc LIMIT 1
			_, err = m.Order("age asc").Limit(1).Delete(1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("limit-delete: ok")
			}

			_, err = m.Delete("uid>", 3)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("delete: ok")
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}