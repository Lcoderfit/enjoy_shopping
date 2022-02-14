package main

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-事务处理
https://goframe.org/pages/viewpage.action?pageId=1114300

一、什么时候开启事务？？？
		当涉及多次数据修改或插入操作时，且数据必须保持一致性，要么全成功，要么全失败，这就需要使用事务
        例如：添加用户接口，首先添加用户需要一次写入操作，然后权限表设置用户的权限需要一次，岗位表又需要写入一次,这三张表的数据需要保持一致性
        就要用到事务(因为一旦中途操作失败了,会进行数据回滚,如果不开启事务则会出现修改了一个表,但是其他几个表的数据没有修改的情况,数据的一致性就被破坏了)

		在比例本例中，在user表中添加一行数据，需要在user_detail中添加一条相应的数据，且uid要与user表中插入数据的ID保持一致,所以需要开始事务；
		确保两张表数据的一致性

TODO:什么是TX事务???

二、g.DB().Transaction(ctx context.Context, f func(ctx context.Context, tx *TX) error) error
	2.1 当Transaction闭包中返回err为nil时，则闭包执行结束后自动提交commit事务，err!=nil或产生panic中断，则会导致事务的回滚(rollback)
	2.2 在goframe项目中，一般的dao操作可以通过不断传递ctx绑定到同一个事务
		dao.User.Transaction(ctx, func(ctx context.Context), tx *gdb.Tx) error {
			// 通过绑定到同一个ctx,绑定到同一个事务
			dao.User.Ctx(txt).....
			dao.UserDetail.Ctx(txt)
		}

	2.3 也可以通过同一个tx对象进行操作,使所有表操作在同一个事务内进行
		g.DB().Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			tx.Model("user").xx
			tx.Model("user_detail")
		}

	2.4 通过g.DB().Table()或g.DB().Model()方法的TX(tx)方法绑定到同一个tx事务对象
			// 1.通过g.DB().Begin()获取一个事务操作对象tx
			// 2.通过g.DB().Model().TX(tx) 或 g.DB().Table().TX(tx) 绑定到同一个事务操作对象进行操作
			tx, err := g.DB().Begin()
			if err != nil {
				return err
			}


		func Register() error {
			var (
				uid int64
				err error
			)
			tx, err := g.DB().Begin()
			if err != nil {
				return err
			}
			// 方法退出时检验返回值，
			// 如果结果成功则执行tx.Commit()提交,
			// 否则执行tx.Rollback()回滚操作。
			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()
			// 写入用户基础数据
			uid, err = AddUserInfo(tx, g.Map{
				"name":  "john",
				"score": 100,
				//...
			})
			if err != nil {
				return err
			}
			// 写入用户详情数据，需要用到上一次写入得到的用户uid
			err = AddUserDetail(tx, g.Map{
				"uid":   uid,
				"phone": "18010576259",
				//...
			})
			return err
		}

		func AddUserInfo(tx *gdb.TX, data g.Map) (int64, error) {
			result, err := g.Table("user").TX(tx).Data(data).Insert()
			if err != nil {
				return 0, err
			}
			uid, err := result.LastInsertId()
			if err != nil {
				return 0, err
			}
			return uid, nil
		}

		func AddUserDetail(tx *gdb.TX, data g.Map) error {
			_, err := g.Table("user_detail").TX(tx).Data(data).Insert()
			return err
		}

*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB()
		group.GET("/tx", func(r *ghttp.Request) {
			err := m.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
				d := g.Map{
					"name":   "x2",
					"gender": 1,
					"num":    10,
				}
				res, err := tx.Model("user").Insert(d)
				if err != nil {
					return err
				}

				// 在执行事务时如果返回err!= nil, 则会rollback
				// return gerror.New("test-error")

				id, err := res.LastInsertId()
				res, err = tx.Model("user_detail").Insert(g.Map{
					"uid":   id,
					"level": "p5",
				})
				return err
			})
			if err != nil {
				g.Log().Line(true).Error(err.Error())
			} else {
				r.Response.Writeln("tx: ok")
			}
		})

		group.GET("/raw-tx", func(r *ghttp.Request) {

		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
