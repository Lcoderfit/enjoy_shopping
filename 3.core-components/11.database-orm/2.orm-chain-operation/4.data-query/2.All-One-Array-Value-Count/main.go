package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-All/One/Array/Value/Count
https://goframe.org/pages/viewpage.action?pageId=17204071

一、数据查询常用方法
	1.1 用于查询多条记录，构成切片返回-可以查询所有数据,列任意选
		func (m *Model) All(where ...interface{} (Result, error)

	1.2 用于查询单条记录-查询所有数据,列任意选,但是会limit 1
		func (m *Model) One(where ...interface{}) (Record, error)

	1.3 用于查询单个字段列构成的切片-可以传入多个列字段名,但是最终只会返回单个列值构成的切片
		Array()的第一个参数一定是表示查询的列, 从第二个参数开始才是where条件
		Array("name", "uid>0")    select name from user where uid>0
		Array("name", "uid>", 0)  select name from user where uid>0
		Fields("name").Where("uid>", 0).Array() 同上
		func (m *Model) Array(fieldsAndWhere ...interface{}) ([]Value, error)

	1.4 查询并返回一个字段值-也是只能查询一个字段列的值,且返回数据limit 1,Value第一个参数一定是字段名
		注意: 返回的不是一个切片,是一个字段值
		Value("name", "uid>0") select name from user where uid>0 limit 1
		func (m *Model) Value(fieldsAndWhere ...interface{}) (Value, error)

	1.5 查询数据总条数
		Count("uid>", 0) => SELECT COUNT(1) FROM `user` WHERE uid>0
		func (m *Model) Count(where ...interface{}) (int, error)

	1.6 count(列名)
		count("name") => SELECT COUNT(name) FROM `user`
		func (m *Model) CountColumn(column string) (int, error)

	区别：
		1.Update/Delete/Insert 返回的是 (sql.Result, error)类型
		2.Wherexx 系列返回的是*Model
		3.All() (Result, error)
          One() (Record, error)
          Array() ([]Value, error)
		  Value() (Value, error)
  		  Count() (int, error)
		  CountColumn (int, error)

二、Find*支持主键的查询条件
	以下所有方法均支持主键查询,当直接传入条件参数不传入字段参数时,会智能识别为主键查询(相当于WherePri)
	注意:FindAll/FindOne第一个参数都是查询参数(主键查询)
		FindArray/FindValue 第一个参数还是列名参数

	func (m *Model) FindAll(where ...interface{}) (Result, error)
	func (m *Model) FindOne(where ...interface{}) (Record, error)
	func (m *Model) FindArray(fieldsAndWhere ...interface{}) (Value, error)
	func (m *Model) FindValue(fieldsAndWhere ...interface{}) (Value, error)
	func (m *Model) FindCount(where ...interface{}) (int, error)
	func (m *Model) FindScan(pointer interface{}, where ...interface{}) error

	// SELECT * FROM `scores` WHERE `id`=1
	Model("scores").FindAll(1)

	// SELECT * FROM `scores` WHERE `id`=1 LIMIT 1
	Model("scores").FindOne(1)

	// SELECT `name` FROM `scores` WHERE `id`=1
	// 返回 ["xxxx"]
	Model("scores").FindArray("name", 1)

	// SELECT `name` FROM `scores` WHERE `id`=1 LIMIT 1
	// 返回 "xxxx"
	Model("user").FindValue("name", 1)

	// SELECT COUNT(1) FROM `user`  WHERE `id`=1
	Model("user").FindCount(1)
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/", func(r *ghttp.Request) {
			// All: [{"age":2,"name":"lh","nickname":null,"num":13,"uid":1},{"age":2,"name":"jsy","nickname":null,"num":12,"uid":2}]
			res, err := m.Where("uid>", 0).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("All: ", err.Error())
			} else {
				r.Response.Writeln("All: ", res)
			}

			// one: {"age":2,"name":"lh","nickname":null,"num":13,"uid":1}
			one, err := m.Where("uid>", 0).One()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("one: ", one)
			}

			// array: ["jsy","lh"]
			// 注意: 这里无法查询多个字段,如果传入 m.Array("uid, name", "uid>", 0)则最后会查询到哪个字段不确定,
			// 但是一定只会返回一个字段列的值构成的切片
			array, err := m.Array("name", "uid>", 0)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("array: ", err.Error())
			} else {
				r.Response.Writeln("array: ", array)
			}

			// value: jsy
			// 只会返回一个字段的值
			value, err := m.Value("name", "uid>", 0)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("value: ", err.Error())
			} else {
				r.Response.Writeln("value: ", value)
			}

			// SELECT COUNT(1) FROM `user` WHERE uid>0
			count, err := m.Count("uid>", 0)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("count: ", err.Error())
			} else {
				r.Response.Writeln("count: ", count)
			}

			// SELECT COUNT(name) FROM `user`
			colCount, err := m.CountColumn("name")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("colCount: ", err.Error())
			} else {
				r.Response.Writeln("colCount: ", colCount)
			}
		})

		// Find*
		group.GET("/find", func(r *ghttp.Request) {
			res, err := m.FindAll(1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("FindAll: ", err.Error())
			} else {
				r.Response.Writeln("FindAll: ", res)
			}

			one, err := m.FindOne(1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("FindOne: ", err.Error())
			} else {
				r.Response.Writeln("FindOne: ", one)
			}

			// select name from user where id=1
			array, err := m.FindArray("name", 1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("FindArray: ", err.Error())
			} else {
				r.Response.Writeln("FindArray: ", array)
			}

			// select name from user where id=1
			value, err := m.FindValue("name", 1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("FindValue: ", err.Error())
			} else {
				r.Response.Writeln("FindValue: ", value)
			}

			// SELECT COUNT(1) FROM `user` WHERE id=1
			count, err := m.FindCount(1)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln("FindCount: ", err.Error())
			} else {
				r.Response.Writeln("FindCount: ", count)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
