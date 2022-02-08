package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-常用操作示例

g.Model("user") 读取默认的数据库配置(default), 基于user表创建一个*Model对象
g.DB().Model("user")与上面等效
g.DB("other").Model("user")读取配置中的other数据库配置，然后基于user表生成一个*Model对象

Where
Wherexxx/WhereNotxxx/WhereOrxxx/WhereOrNotxxx

一、in查询
	Where五种方式:
	1.1 直接传入slice: m.Where("uid", g.Slice{1})
	1.2 使用一个占位符,传递slice: Where("uid in (?)", g.Slice{1, 2})
	1.3 使用多个占位符,每一个占位符对应一个参数: Where("uid in (?,?)", 1, 2)

	1.4 map: m.Where(g.Map{"uid", g.Slice{1,2}})
	1.5 struct: m.Where(User{"uid", g.Slice{1,2}})

	NOTE: 当传递的Slice参数为空或nil时,查询不会报错,而是转换为一个false语句
	// SELECT * FROM `user` WHERE 0=1
	db.Model("user").Where("uid", g.Slice{}).All()
	db.Model("user").Where("uid", nil).All()

	1.6 WhereIn/WhereNotIn/WhereOrIn/WhereOrNotIn
		// SELECT * FROM `user` WHERE `gender`=1 AND `type` IN(1,2,3)
		db.Model("user").Where("gender", 1).WhereIn("type", g.Slice{1,2,3}).All()

		// SELECT * FROM `user` WHERE `gender`=1 AND `type` NOT IN(1,2,3)
		db.Model("user").Where("gender", 1).WhereNotIn("type", g.Slice{1,2,3}).All()

		// SELECT * FROM `user` WHERE `gender`=1 OR `type` IN(1,2,3)
		db.Model("user").Where("gender", 1).WhereOrIn("type", g.Slice{1,2,3}).All()

		// SELECT * FROM `user` WHERE `gender`=1 OR `type` NOT IN(1,2,3)
		db.Model("user").Where("gender", 1).WhereOrNotIn("type", g.Slice{1,2,3}).All()

二、like查询
	2.1 Where的一种方法:
		m.Where("name like ?", "l%")
	2.2 Wherexxx系列方法(都是与Where链式操作一起使用)
		func (m *Model) WhereLike(column string, like interface{}) *Model
		func (m *Model) WhereNotLike(column string, like interface{}) *Model
		func (m *Model) WhereOrLike(column string, like interface{}) *Model
		func (m *Model) WhereOrNotLike(column string, like interface{}) *Model

		// SELECT * FROM `user` WHERE `gender`=1 AND `name` LIKE 'john%'
		db.Model("user").Where("gender", 1).WhereLike("name", "john%").All()

		// SELECT * FROM `user` WHERE `gender`=1 AND `name` NOT LIKE 'john%'
		db.Model("user").Where("gender", 1).WhereNotLike("name", "john%").All()

		// SELECT * FROM `user` WHERE `gender`=1 OR `name` LIKE 'john%'
		db.Model("user").Where("gender", 1).WhereOrLike("name", "john%").All()

		// SELECT * FROM `user` WHERE `gender`=1 OR `name` NOT LIKE 'john%'
		db.Model("user").Where("gender", 1).WhereOrNotLike("name", "john%").All()

三、min/max/avg/sum
	3.1 直接将聚合函数放入Fieldsxxx()函数中
		// SELECT MIN(score) FROM `user` WHERE `uid`=1
		db.Model("user").Fields("MIN(score)").Where("uid", 1).Value()

		// SELECT MAX(score) FROM `user` WHERE `uid`=1
		db.Model("user").Fields("MAX(score)").Where("uid", 1).Value()

		// SELECT AVG(score) FROM `user` WHERE `uid`=1
		db.Model("user").Fields("AVG(score)").Where("uid", 1).Value()

		// SELECT SUM(score) FROM `user` WHERE `uid`=1
		db.Model("user").Fields("SUM(score)").Where("uid", 1).Value()

	3.2 直接使用GF内置的聚合函数, 返回(float64, error)类型
		func (m *Model) Min(column string) (float64, error)
		func (m *Model) Max(column string) (float64, error)
		func (m *Model) Avg(column string) (float64, error)
		func (m *Model) Sum(column string) (float64, error)

		// SELECT MIN(`score`) FROM `user` WHERE `uid`=1
		db.Model("user").Where("uid", 1).Min("score")
		
		// SELECT MAX(`score`) FROM `user` WHERE `uid`=1
		db.Model("user").Where("uid", 1).Max("score")
		
		// SELECT AVG(`score`) FROM `user` WHERE `uid`=1
		db.Model("user").Where("uid", 1).Avg("score")
		
		// SELECT SUM(`score`) FROM `user` WHERE `uid`=1
		db.Model("user").Where("uid", 1).Sum("score")

四、count
	4.1 select count(*) from user where uid>0, NOTE: 返回(int, error)类型
		m.Where("uid>0").Count()
	4.2 对指定字段进行count统计
		func (m *Model) CountColumn(column string) (int, error)
		// select count(num) from user, NOTE: num值为null的行不会统计
		m.CountColumn("num")
	4.3 在Fieldsxx中使用count???

五、distinct
	5.1 对名字去重: select distinct name from user
		m.Fields("distinct name").All()
	5.2 通过Distinct()方法去重
		func (m *Model) Distinct() *Model

		// SELECT COUNT(DISTINCT `name`) FROM `user`
		db.Model("user").Distinct().CountColumn("name")

		// SELECT COUNT(DISTINCT uid,name) FROM `user`
		db.Model("user").Distinct().CountColumn("uid,name")

六、between
	6.1 直接使用Where+占位符
		Where("num between ? and ?", 0, 30)
	6.2 Wherexx系列方法
		func (m *Model) WhereBetween(column string, min, max interface{}) *Model
		func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model
		func (m *Model) WhereOrBetween(column string, min, max interface{}) *Model
		func (m *Model) WhereOrNotBetween(column string, min, max interface{}) *Model

		// SELECT * FROM `user` WHERE `gender`=0 AND `age` BETWEEN 16 AND 20
		db.Model("user").Where("gender", 0).WhereBetween("age", 16, 20).All()

		// SELECT * FROM `user` WHERE `gender`=0 AND `age` NOT BETWEEN 16 AND 20
		db.Model("user").Where("gender", 0).WhereNotBetween("age", 16, 20).All()

		// SELECT * FROM `user` WHERE `gender`=0 OR `age` BETWEEN 16 AND 20
		db.Model("user").Where("gender", 0).WhereOrBetween("age", 16, 20).All()

		// SELECT * FROM `user` WHERE `gender`=0 OR `age` NOT BETWEEN 16 AND 20
		db.Model("user").Where("gender", 0).WhereOrNotBetween("age", 16, 20).All()

七、null
	7.1 Where直接使用sql语句
		m.Where("name is not null").All()
	7.2 Where系列方法:
		// where name is null and gender is null
		m.WhereNull("name, gender").All()

		func (m *Model) WhereNull(columns ...string) *Model
		func (m *Model) WhereNotNull(columns ...string) *Model
		func (m *Model) WhereOrNull(columns ...string) *Model
		func (m *Model) WhereOrNotNull(columns ...string) *Model
*/

func main()  {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/in", func(r *ghttp.Request) {
			// 1.不使用占位符,直接传递slice
			res, err := m.Where("uid", g.Slice{1}).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("in1: ", res)
			}

			// 2.使用一个占位符,传递slice
			res, err = m.Where("uid in (?)", g.Slice{1, 2}).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("in2: ", res)
			}

			// 3.使用多个占位符,每一个占位符对应一个参数
			res, err = m.Where("uid in (?,?)", 1, 2).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("in3: ", res)
			}
		})

		group.GET("/like", func(r *ghttp.Request) {
			// 直接使用占位符
			// where name like 'l%'
			res, err := m.Where("name like ?", "l%").All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("like: ", res)
			}
		})

		group.GET("/jh", func(r *ghttp.Request) {
			res, err := m.Fields("min(num)").Value()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("jh: ", res)
			}
			
			result, err := m.Sum("num")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("jh2: ", result)
			}
		})

		group.GET("/count", func(r *ghttp.Request) {
			res, err := m.Where("uid>0").Count()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("count: ", res)
			}

			// select count(num) from user, NOTE: num值为null的行不会统计
			res, err = m.CountColumn("num")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("count2: ", res)
			}
		})

		group.GET("/distinct", func(r *ghttp.Request) {
			res, err := m.Fields("distinct name").All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("distinct: ", res)
			}
		})

		group.GET("/between", func(r *ghttp.Request) {
			res, err := m.Where("num between ? and ?", 0, 30).All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("between: ", res)
			}
		})

		group.GET("/null", func(r *ghttp.Request) {
			res, err := m.Where("name is not null").All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("null: ", res)
			}

			res, err = m.WhereNull("name").All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("null2: ", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}