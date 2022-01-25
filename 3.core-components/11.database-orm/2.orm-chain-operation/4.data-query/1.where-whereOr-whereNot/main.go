package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-Where/WhereOr/WhereNot
https://goframe.org/pages/viewpage.action?pageId=17204072

func (m *Model) Where(where interface{}, args...interface{}) *Model
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model
func (m *Model) WhereBetween(column string, min, max interface{}) *Model
func (m *Model) WhereLike(column string, like interface{}) *Model
func (m *Model) WhereIn(column string, in interface{}) *Model
func (m *Model) WhereNull(columns ...string) *Model
func (m *Model) WhereLT(column string, value interface{}) *Model
func (m *Model) WhereLTE(column string, value interface{}) *Model
func (m *Model) WhereGT(column string, value interface{}) *Model
func (m *Model) WhereGTE(column string, value interface{}) *Model

func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model
func (m *Model) WhereNotLike(column string, like interface{}) *Model
func (m *Model) WhereNotIn(column string, in interface{}) *Model
func (m *Model) WhereNotNull(columns ...string) *Model

func (m *Model) WhereOr(where interface{}, args ...interface{}) *Model
func (m *Model) WhereOrBetween(column string, min, max interface{}) *Model
func (m *Model) WhereOrLike(column string, like interface{}) *Model
func (m *Model) WhereOrIn(column string, in interface{}) *Model
func (m *Model) WhereOrNull(columns ...string) *Model
func (m *Model) WhereOrLT(column string, value interface{}) *Model
func (m *Model) WhereOrLTE(column string, value interface{}) *Model
func (m *Model) WhereOrGT(column string, value interface{}) *Model
func (m *Model) WhereOrGTE(column string, value interface{}) *Model

func (m *Model) WhereOrNotBetween(column string, min, max interface{}) *Model
func (m *Model) WhereOrNotLike(column string, like interface{}) *Model
func (m *Model) WhereOrNotIn(column string, in interface{}) *Model
func (m *Model) WhereOrNotNull(columns ...string) *Model

一、Where/WhereOr
	// 注意： MySQL中And的优先级高于Or
	// c1 and c2 or c3 and c4  <=>  (c1 and c2) or (c3 and c4)
	1.1 Where的五种查询方式
		1.1.1 字符串			 Where("uid=1")/Where("name='lh'")
		1.1.2 字段和值分开传递: Where("name", "lh")/Where("uid>=", 0)
		1.1.3 Map			Where(g.Map{"name", "lh"})/Where(g.Map{"uid>=", 0})
		1.1.4 占位符:		Where("uid>=?", 0)/Where("uid>=? or name=?", 2, "lh")
			Note: 如果是子查询,则要带上占位符?
		1.1.5 Struct参数
			type Condition struct {
				Sex int `orm:sex`
				Age int `orm:age`
			}
			cd := Condition{1, 18}
			Where(cd)

		1.1.6 Where().Or()  Where().And()

		case:
			condition := g.Map{
				"title like ?"         : "%九寨%",
				"online"               : 1,
				"hits between ? and ?" : g.Slice{1, 10},
				"exp > 0"              : nil,
				"category"             : g.Slice{100, 200},
			}
			// SELECT * FROM article WHERE title like '%九寨%' AND online=1 AND hits between 1 and 10
			// AND exp > 0 AND category IN(100,200)
			db.Model("article").Where(condition).All()

二、WherePri主键查询
	WherePri(1) <=> Where("uid", 1)
	WherePri(g.Slice{1,2,3}) <=> Where("uid in(?)", g.Slice{1,2,3})
*/

type User struct {
	Uid      uint64
	Name     string
	Nickname string
	Age      int
	Num      int
}

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.GET("/where-whereOr", func(r *ghttp.Request) {
			var user *User
			// SELECT `uid`,`name`,`nickname`,`age`,`num` FROM `user` WHERE (uid <= 0) OR (name='lh') AND (name='jsy') LIMIT 1
			// And优先级>Or, 相当于: ((uid <= 0) OR (name='lh')) AND (name='jsy')
			err := m.Where("uid <= ?", 0).WhereOr("name=?", "lh").
				Where("name=?", "jsy").Scan(&user)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				return
			}
			r.Response.Writeln("where-where-or: ok")
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
