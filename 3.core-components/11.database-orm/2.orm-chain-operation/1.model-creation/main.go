package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
orm链式操作-模型创建

1.创建基于数据表的Model对象
	1.1 g.Model("user") 创建一个基于user表的Model对象(*Model类型)
	1.2 g.DB().Model("user") 读取默认的数据库配置并创建Model对象
	1.3 切换当前的数据库对象
		m := g.Model("user")
		m = m.DB(g.DB("other-user")) 从default配置项切换到other-user配置项

		等效于: m := g.DB("other-user").Model("user")

	1.4 创建一个基于原始sql的Model对象
		1.4.1 g.DB().Raw()
			sql := "select * from t where id in(?)"
			// LT => less than: a < b
			result, err := g.DB().Raw(sql, g.Slice{2, 3}).Where("b", "second").All()
		1.4.2 g.Model("user").Raw()
			sql := "select a, b from t where id>?"
			// args传入字符串和整型都可以,内部应该会根据数据表字段类型进行转换
			result, err := g.Model("user").Raw(sql, "1").All()

2.链式安全
	2.1 链式不安全: 会修改当前Model对象
		user := g.DB().Model("user")
		// user对象会被修改, user对象相当于 select * from user where id=1;
		user.Where("id", 1)

	2.2 链式安全: 每一个链式操作都会返回一个新的Model对象,不会对当前Model对象进行修改
		user := g.DB().Model("user").Safe()
		// Where条件不会修改user这个Model对象,user相当于 select * from user; (无where条件)
		// 同时返回一个新的Model对象m,m对应 select * from user where id=1;
		m := user.Where("id", 1)

3.链式操作(default/Clone/Safe)
	3.1 默认情况下,gdb链式操作是非链式不安全的(链式操作的每一个方法都会对当前Model对象进行修改)
			user := g.Model("t")
			m := user.Where("id in (?)", g.Slice{1, 2, 3})

			// m和user的地址是一样的,链式操作会修改当前的Model
			// m和user均相当于: select * from t where id in (1,2,3);
			g.Log().Line(true).Infof("m:%p, t:%p", m, user)

	3.2 Clone()方法 对当前模型进行克隆,创建一个新的Model对象,然后对新的Model对象进行链式操作,这样就不会修改原有的Model
			t := g.Model("t")
			// t相当于 select * from t where id > 0;
			// m相当于 select * from t where id > 0 and b='second';
			// t是非链式安全的,t.Where()会对t进行修改,叠加查询条件
			t.Where("id>?", 0)
			m := t.Clone()
			m.Where("b", "second")

	3.3 Safe()方法 设置当前模型为链式安全的对象,后续的每一个链式操作对象都将返回一个新的Model对象,
		如果需要对模型属性进行修改或叠加查询, 使用m = m.xxx 对原有对象进行覆盖
			// t 相当于 select * from t;
			// m相当于 select * from t where id > 0 and b='second';
			t := g.Model("t").Safe()
			m := t.Where("id>?", 0)
			// 查询条件的叠加,通过m=m.xxx对原有对象进行覆盖
			m = m.Where("b", "second")

		注意:使用Safe()方法, m.Where()会返回一个新的Model对象,但是不会修改m本身
			m = m.Where() 由于返回的新Model对象又赋值给了m,所以m包含了新的Where条件

*/

func main() {
	s := g.Server()
	// gf run main.go 启动时当前项目路径为main.go所在的路径
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		// 创建基于原始sql的model对象
		group.GET("/row", func(r *ghttp.Request) {
			sql := "select * from t where id in(?)"
			// LT => less than: a < b
			result, err := g.DB().Raw(sql, g.Slice{2, 3}).Where("b", "second").All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info(result.List())
			r.Response.Writeln("hello world")
		})

		group.GET("/model-raw", func(r *ghttp.Request) {
			sql := "select a, b from t where id>?"
			// args传入字符串和整型都可以,内部应该会根据数据表字段类型进行转换
			result, err := g.Model("user").Raw(sql, "1").All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info(result.List())
			r.Response.WriteJson("hello world")
		})

		// 非链式安全操作
		group.GET("/chain-unsafe", func(r *ghttp.Request) {
			user := g.Model("t")
			m := user.Where("id in (?)", g.Slice{1, 2, 3})

			// m和user的地址是一样的,链式操作会修改当前的Model
			// m和user均相当于: select * from t where id in (1,2,3);
			g.Log().Line(true).Infof("m:%p, t:%p", m, user)

			result, err := m.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("m: ", result.List())

			result, err = user.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("t: ", result.List())
			r.Response.WriteJson("hello world")
		})

		// Clone方法
		group.GET("/clone", func(r *ghttp.Request) {
			t := g.Model("t")
			// t相当于 select * from t where id > 0;
			// m相当于 select * from t where id > 0 and b='second';
			t.Where("id>?", 0)
			m := t.Clone()
			m.Where("b", "second")

			result, err := t.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("t: ", result.List())

			result, err = m.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("m: ", result.List())

			r.Response.WriteJson("hello world")
		})

		// Safe()方法
		group.GET("/safe", func(r *ghttp.Request) {
			// t 相当于 select * from t;
			// m相当于 select * from t where id > 0 and b='second';
			t := g.Model("t").Safe()
			m := t.Where("id>?", 0)
			m = m.Where("b", "second")

			result, err := t.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("t: ", result.List())

			result, err = m.All()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Exit()
			}
			g.Log().Line(true).Info("m: ", result.List())

			r.Response.WriteJson("hello world")
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
