package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)


/*orm链式操作-写入保存
https://goframe.org/pages/viewpage.action?pageId=17203727

一、数据插入
	Data()和Insert()/InsertIgnore()/Save()/Replace()均可以接受参数，用来设置插入的数据的字段值
	参数可以是:
		map[string]string/map[string]interface{}/map[interface{}]interface{}
		struct/\*struct/[]struct/[]\*struct (这个'\'可以忽略，是为了转译\*使用的)

	这些函数均返回(sql.Result, error)类型

	1.1 g.DB.Model("user").Data(data).Insert() / g.DB.Model("user").Insert(data)
		插入数据，如果存在主键或唯一索引重复，则返回错误

	1.2 g.DB.Model("user").Data(data).InsertIgnore() / g.DB.Model("user").InsertIgnore(data)
		插入数据，如果存在主键或唯一索引重复，则忽略错误(不再插入)

	1.3 g.DB.Model("user").Data(data).Save() / g.DB.Model("user").Save(data)
		插入数据，如果存在主键或唯一索引重复，则更新数据，否则插入;
		需要注意：如果数据更新后仍存在主键或唯一索引冲突，则仍会报错，参考三、duplicate key update

	1.4 g.DB.Model("user").Data(data).Replace() / g.DB.Model("user").Replace(data)
		插入数据，如果存在主键或唯一索引重复，则将所有重复的数据删除，然后再执行插入

	返回(int64, error)类型
	id, err := g.DB.Model("user").Data(data).InsertAndGetId(g.Map{"name": "t9"})

二、batch 数据批量写入
	2.1 g.DB.Model("user").Data(g.List{
		g.Map{"name":"1"},
		g.Map{"name": "2'},
		g.Map{"name": "3},
	}).Insert() 或者直接将参数传入Insert()

	2.2 分批写入，没批写入2条数据，默认是每一批写入100条
		g.DB.Model("user").Batch(2).Data(g.List{
		g.Map{"name":"1"},
		g.Map{"name": "2'},
		g.Map{"name": "3},
	}).Insert() 或者直接将参数传入Insert()

三、RawSQL语句嵌入
	// uid默认为0
	// gdb.Raw("now()") 解析到sql: insert into user(`create_time`) values(now());
	// insert into user(`uid`, `name`) values(uid+200, 'raw')
	_, err = m.Insert(g.Map{
		"uid": gdb.Raw("uid+200"),
		"name": "raw",
	})

四、insert ignore into
	插入时如果遇到主键或唯一主键重复，则忽略错误
	insert ignore into `user`(`uid`, `name`) values('2', 'tianyi')

五、on duplicate key update
当插入的数据存在主键或唯一索引冲突时,会报错:
    -- 这个表示插入的数据总的name字段存在重复(name字段是一个唯一索引)
    ERROR 1062 (23000): Duplicate entry 'tianyi' for key 'name'
    -- 插入数据时存在主键重复,重复值为2
    ERROR 1062 (23000): Duplicate entry '2' for key 'PRIMARY'

    当存在主键或唯一索引重复时则更新:
        -- `name`=values(`name`)可以动态设置字段name的值,即如果存在主键或唯一索引冲突,更新name的值为values('2', 'tianyi')中name字段的值
        insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`=values(`uid`),`name`=values(`name`);
        -- 插入('2', 'tianyi')时,如果存在重复数据,则更新uid为2, name为'lh'
        insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`='2',`name`='lh';
        -- 注意,如果表中已经同时有uid为2和3的数据,则下面的语句仍会报错,因为插入uid为2的数据会导致主键冲突,此时会更新这条uid为2的数据为('3', 'lh');
        -- 但是uid=3的数据也是存在的,所以仍会: Duplicate entry '3' for key 'PRIMARY'
        insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`='3',`name`='lh';

六、replace into `user`(`uid`, `name`) values('3', 'tianyi');
    在插入数据时,如果遇到主键或唯一索引重复的数据,会先把所有重复数据删除,然后再插入;
    例如表中有('2', 'tianyi') ('3', 'xm'); 则会先删除这两条数据,然后再插入('3', 'tianyi')
*/
func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user").Safe()
		group.Group("/basic-use", func(group *ghttp.RouterGroup) {
			// 1.如果写入的数据中存在主键或者唯一索引冲突时，返回失败，否则写入一条新数据
			// 通过insert into写入数据
			// Error 1062: Duplicate entry 'jsy' for key 'name', INSERT INTO `user`(`name`) VALUES('jsy')
			group.GET("/insert", func(r *ghttp.Request) {
				_, err := m.Data(g.Map{"name": "jsy"}).Insert()
				if err != nil {
					g.Log().Line(true).Error(err)
					r.Response.WriteJsonExit(err.Error())
				}
				r.Response.WriteJsonExit("yes")
			})
			// 2.InsertIgnore插入时,如果写入的数据存在主键或唯一索引冲突(重复),
			group.GET("/insertignore", func(r *ghttp.Request) {
				_, err := m.Data(g.Map{"uid": 2, "name": "lhx"}).InsertIgnore()
				if err != nil {
					g.Log().Line(true).Error(err)
					return
				}
				r.Response.Writeln("insertignore success")
			})
			// 3.Save() 如果写入的数据存在主键或唯一索引冲突,则更新数据,否则新增数据
			// 有则更新,无则新增
			// 如果换成m.Data(g.Map{"name": "tianyi"}),则数据id不会变换(因为是更新数据,而Replace会删除后再新增,所以id会变化)
			group.GET("/save", func(r *ghttp.Request) {
				_, err := m.Data(g.Map{"uid": 2, "name": "tianyi"}).Save()
				if err != nil {
					g.Log().Line(true).Error(err)
					return
				}
				return
			})
			// 4.Replace() 如果遇到主键或唯一索引导致的数据冲突(数据重复),会删除原有数据,必定会新增一条数据
			// 有则删除后新增,无则更新
			// user表中name字段是唯一索引,多次调用该接口,每次会删除表中的数据,然后新增一条,所以每次调用该接口,表中的name为"tianyi"的数据id增加1
			// 如果换成m.Data(g.Map{"uid":2, "name": "tianyi"}) 表中数据id一直为2
			group.GET("/replace", func(r *ghttp.Request) {
				_, err := m.Data(g.Map{"name": "tianyi"}).Replace()
				if err != nil {
					g.Log().Line(true).Error(err)
					return
				}
				r.Response.Writeln("replace success")
			})
		})

		// 使用写入/保存方法直接传递参数
		group.GET("/param", func(r *ghttp.Request) {
			// insert: Error 1062: Duplicate entry 'jsy' for key 'name', INSERT INTO `user`(`name`) VALUES('jsy')
			_, err := m.Insert(g.Map{"name": "jsy"})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("insert: ", err.Error())
			} else {
				r.Response.Writeln("insert: ", "ok")
			}

			_, err = m.InsertIgnore(g.Map{"name": "jsy"})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("insert-ignore: ", err.Error())
			} else {
				r.Response.Writeln("insert-ignore: ", "ok")
			}

			_, err = m.Save(g.Map{"uid": 2, "name": "tianyi-s"})
			if err != nil {
				g.Log().Line(true).Error("save: ", err)
				r.Response.Writeln("save: ", err.Error())
			} else {
				r.Response.Writeln("save: ", "ok")
			}

			_, err = m.Replace(g.Map{"name": "lh"})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("replace: ", err.Error())
			} else {
				r.Response.Writeln("replace: ", "ok")
			}
		})

		// 使用Struct作为参数插入数据
		group.GET("/struct", func(r *ghttp.Request) {
			user := &User{
				Uid:  3,
				Name: "xm",
			}
			_, err := m.Data(user).Save()
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("save1: ", err.Error())
			} else {
				r.Response.Writeln("save1: ", "ok")
			}

			_, err = m.Save(user)
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("save2: ", err.Error())
			} else {
				r.Response.Writeln("save2: ", "ok")
			}
		})

		// 数据批量保存
		group.GET("/batch", func(r *ghttp.Request) {
			_, err := m.Save(g.List{
				{"uid": "10001", "name": "first"},
				{"uid": "10002", "name": "second"},
				{"uid": "10003", "name": "third"},
			})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("batch-save1:", err.Error())
			} else {
				r.Response.Writeln("batch-save1: ", "ok")
			}

			// 分批保存,每批保存两条(会发送两个sql请求,第一个写入前两条数据,第二个请求写入第三条数据)
			_, err = m.Batch(2).Save(g.List{
				{"uid": "10001", "name": "first"},
				{"uid": "10002", "name": "second"},
				{"uid": "10003", "name": "third"},
			})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln("batch-save2:", err.Error())
			} else {
				r.Response.Writeln("batch-save2: ", "ok")
			}
		})

		// 使用原生sql语句嵌入,不会转换为字符串类型
		group.GET("/raw", func(r *ghttp.Request) {
			// insert into user(`uid`, `name`) values('uid+200', 'raw')
			// Error 1366: Incorrect integer value: 'uid+2' for column 'uid' at row 1, INSERT INTO `user`(`uid`,`name`) VALUES('uid+2','not-raw')
			_, err := m.Insert(g.Map{
				"uid":  "uid+2",
				"name": "not-raw",
			})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("not-raw: ", "ok")
			}

			// uid默认为0
			// gdb.Raw("now()") 解析到sql: insert into user(`create_time`) values(now());
			// insert into user(`uid`, `name`) values(uid+200, 'raw')
			_, err = m.Insert(g.Map{
				"uid": gdb.Raw("uid+200"),
				"name": "raw",
			})
			if err != nil {
				g.Log().Line(true).Error(err)
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("raw-sql: ", "ok")
			}
		})

		// 插入数据并获取最后插入的一条数据的id
		// 返回(int64, error)类型
		group.GET("/insert-and-get-id", func(r *ghttp.Request) {
			id, err := m.InsertAndGetId(g.Map{"name": "t9"})
			if err != nil {
				g.Log().Line(true).Error(err)
				return
			}
			g.Log().Line(true).Info("insert-and-get-id success: ", id)
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}

type User struct {
	Uid  uint64
	Name string
}
