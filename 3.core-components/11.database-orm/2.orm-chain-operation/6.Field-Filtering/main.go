package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-字段过滤
TODO:查看评论下的问题
https://goframe.org/pages/viewpage.action?pageId=1114229

一、Fields/FieldsEx字段过滤
	1.1 Fields("xx", "xx") 指定需要操作的字段，例如查询/写入/更新字段， Fields("name").Where(xx) => select name from xxx where xx
	1.2 FieldsEx 指定需要排除的字段, FieldsEx("name").Where(xx) 表示查询时不选择name字段
		m.FieldsEx(xx).Data(xx).Update 更新时过滤一些不必要的字段
二、OmitEmpty 空值过滤 (对Data/Where/Update/Insert/Delete/Save/Replace 等会传入数值参数的函数均有效)
	当map/struct中包含0, nil, "" 等空值字段时，gdb会当作正常的参数更新到数据表，而通过OmitEmpty可以过滤为空值的字段
	2.1 OmitEmpty会过滤Data()和Where()及其他可传入数据参数的函数中为空值的字段

	2.2 OmitEmptyWhere 过滤传入Where()的数据参数中为空的字段(只对Where条件有效，Update/Insert/xx等无效)
		// select * from user where num=10, 过滤了name和gender的where条件
		OmitEmptyWhere().Where(g.Map{"name": "", "gender": 0, "num": 10})

	2.3 OmitEmptyData 过滤传入Data中为空的字段(Update/Insert...均有效,但是Where()中的不受影响)

	NOTE: 批量写入/更新操作中OmitEmpty方法将会失效，因为在批量操作中，必须保证每个写入记录的字段是统一的

	关于omitempty标签与OmitEmpty方法：
		针对于struct的空值过滤大家会想到omitempty的标签。该标签常用于json转换的空值过滤，也在某一些第三方的ORM库中用作struct到数据表字段的空值过滤，
		即当属性为空值时不做转换。
		omitempty标签与OmitEmpty方法所达到的效果是一样的。在ORM操作中，我们不建议对struct使用omitempty的标签来控制字段的空值过滤，
		而建议使用OmitEmpty方法来做控制。因为该标签一旦加上之后便绑定到了struct上，没有办法做灵活控制；而通过OmitEmpty方法使得开发者可以选择性地、
		根据业务场景对struct做空值过滤，操作更加灵活。

三、OmitNil 空值过滤
	与OmitEmpty相似，但是只会对map/struct中为nil的字段进行过滤
	3.1 OmitNil
	3.2 OmitWhereNil 删除where条件中为nil的条件
	3.3 OmitDataNil 对Data中为nil的字段进行过滤
*/

type User struct {
	Uid    uint64
	Name   string
	Gender int8
	Num    int
}

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/fields", func(r *ghttp.Request) {
			// select name, gender from user where 1;
			res, err := m.Fields("name", "gender").Where(1).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("fields: ", res)
			}
		})

		group.GET("/fields-ex", func(r *ghttp.Request) {
			res, err := m.FieldsEx("uid, num").Where(1).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("fields-ex: ", res)
			}
		})

		group.GET("/omit-empty", func(r *ghttp.Request) {
			res, err := m.OmitEmpty().Update(g.Map{
				"gender": nil,
				"num":    12,
			}, "uid", 2)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("omit-empty: ", res)
			}

			// SELECT * FROM `user` WHERE `name`='lh'
			user := &User{
				Uid:    0,
				Name:   "lh",
				Gender: 0,
				Num:    0,
			}
			_, err = m.OmitEmptyWhere().Where(user).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("omit-empty-where: ", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
