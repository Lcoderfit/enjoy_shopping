package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM方法操作

一、执行原生SQL
	1.Query(query string, args ...interface{}) (*sql.Rows, error)
		一般用于查询
	2.Exec(query string, args ...interface{}) (*sql.Result, error)
		用于更新和插入
	3.Prepare(query string) (*sql.Stmt, error)
		TODO:预处理

二、查询推荐使用Getxx
	// 数据表记录查询：
	// 查询单条记录、查询多条记录、获取记录对象、查询单个字段值(链式操作同理)
	GetAll(sql string, args ...interface{}) (Result, error)
	GetOne(sql string, args ...interface{}) (Record, error)
	GetValue(sql string, args ...interface{}) (Value, error)
	GetArray(sql string, args ...interface{}) ([]Value, error)
	GetCount(sql string, args ...interface{}) (int, error)
	GetScan(objPointer interface{}, sql string, args ...interface{}) error

三、单条数据插入
	// 数据单条操作
	Insert(table string, data interface{}, batch...int) (sql.Result, error)
	Replace(table string, data interface{}, batch...int) (sql.Result, error)
	Save(table string, data interface{}, batch...int) (sql.Result, error)

	// 数据修改/删除
	Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error)
	Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)

N:
	Insert/Replace/Save方法中的data参数支持的数据类型为：string/map/slice/struct/\*struct，当传递为slice类型时，自动识别为批量操作，此时batch参数有效。
	原生获取默认的数据库配置对象: gdb.Instance()
	N:需要设置原生配置:https://goframe.org/pages/viewpage.action?pageId=1114245

	// 注意Link中不能包含type(即mysql),需要单独定义type
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Link:  "root:Lcoder66242@tcp(47.101.48.37:3306)/lc-sql?charset=utf8",
				Role:  "master",
				Type:  "mysql",
				Debug: true,
			},
		},
	})

	gdb.SetConfig(gdb.Config {
		"default" : gdb.ConfigGroup {
			gdb.ConfigNode {
				Host     : "192.168.1.100",
				Port     : "3306",
				User     : "root",
				Pass     : "123456",
				Name     : "test",
				Type     : "mysql",
				Role     : "master",
				Weight   : 100,
			},
			gdb.ConfigNode {
				Host     : "192.168.1.101",
				Port     : "3306",
				User     : "root",
				Pass     : "123456",
				Name     : "test",
				Type     : "mysql",
				Role     : "slave",
				Weight   : 100,
			},
		},
		"user-center" : gdb.ConfigGroup {
			gdb.ConfigNode {
				Host     : "192.168.1.110",
				Port     : "3306",
				User     : "root",
				Pass     : "123456",
				Name     : "test",
				Type     : "mysql",
				Role     : "master",
				Weight   : 100,
			},
		},
	})
*/

func main() {
	s := g.Server()
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Link:  "root:Lcoder66242@tcp(47.101.48.37:3306)/lc-sql?charset=utf8",
				Role:  "master",
				Type:  "mysql",
				Debug: true,
			},
		},
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		// 原生获取默认的数据库配置对象
		db, err := gdb.Instance()
		if err != nil {
			g.Log().Line(true).Error(err.Error())
			return
		}
		group.GET("/raw", func(r *ghttp.Request) {
			res, err := db.Insert("user", g.Map{
				"name": "xx1",
			})
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("res-insert: ", res)
			}

			all, err := db.GetAll("select * from user limit 2")
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("res-get-all: ", all)
			}

			// 批量操作,没一批次插入2条数据
			_, err = db.Insert("user", g.List{
				g.Map{"name": "xx2"},
				g.Map{"name": "xx3"},
				g.Map{"name": "xx4"},
				g.Map{"name": "xx5"},
			}, 2)
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("res-batch: ok")
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
