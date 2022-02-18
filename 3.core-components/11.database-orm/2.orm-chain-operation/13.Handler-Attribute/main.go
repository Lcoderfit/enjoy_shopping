package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-Handler特性

一、Handler复用常用逻辑
	1.1 如果有重复的查询逻辑，可以封装成一个handler方法
		func(m *gdb.Model) *gdb.Model {
			// 添加对m的链式操作然后返回
			return m.Where(xxx, xxxx)
		}

	1.2 分页操作
		// SELECT * FROM `user` LIMIT 0,2
		res, err := m.Handler(Paginate(r)).All()
		Paginate是一个闭包函数,返回func(m *gdb.Model) *gdb.Model类型
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/handler", func(r *ghttp.Request) {
			res, err := m.Handler(NameLengthCondition, GenderCondition).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("res: ", res)
			}
		})

		group.GET("/pagination", func(r *ghttp.Request) {
			// SELECT * FROM `user` LIMIT 0,2
			res, err := m.Handler(Paginate(r)).All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("res: ", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}

func Paginate(r *ghttp.Request) func(m *gdb.Model) *gdb.Model {
	// 这是一个闭包
	return func(m *gdb.Model) *gdb.Model {
		type Pagination struct {
			PageNum  int
			PageSize int
		}

		var pagination Pagination
		// 从r中获取前端传入的page_num和page_size参数，然后解析到结构体字段中
		_ = r.Parse(&pagination)
		if pagination.PageNum > 100 {
			pagination.PageNum = 100
		}
		if pagination.PageSize <= 0 {
			pagination.PageSize = 2
		}
		return m.Page(pagination.PageNum, pagination.PageSize)
	}
}

func NameLengthCondition(m *gdb.Model) *gdb.Model {
	return m.Where("length(name) >= ?", 2)
}

func GenderCondition(m *gdb.Model) *gdb.Model {
	return m.Where("gender=", 1)
}
