package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
orm使用配置:
https://goframe.org/pages/viewpage.action?pageId=1114245

1.配置项
	1.1 配置项含义
	[database]
		# 每一个分组中可以配置多个节点，一个master(主节点)，多个slave(从节点);
		# 参考D:\PrivateProject\Gf-Tags\learn-gf\3.core-components\11.database-orm\config\config.toml中的配置
		[[database.分组名称]]
			host                 = "地址"
			port                 = "端口"
			user                 = "账号"
			pass                 = "密码"
			name                 = "数据库名称"
			type                 = "数据库类型(mysql/pgsql/mssql/sqlite/oracle)"
			link                 = "(可选)自定义数据库链接信息，当该字段被设置值时，以上链接字段(Host,Port,User,Pass,Name)将失效，但是type必须有值"
			role                 = "(可选)数据库主从角色(master/slave)，不使用应用层的主从机制请均设置为master"
			debug                = "(可选)开启调试模式"
			prefix               = "(可选)表名前缀"
			dryRun               = "(可选)ORM空跑(只读不写)"
			charset              = "(可选)数据库编码(如: utf8/gbk/gb2312)，一般设置为utf8"
			weight               = "(可选)负载均衡权重，用于负载均衡控制，不使用应用层的负载均衡机制请置空"
			timezone             = "(可选)时区配置，例如:local"
			maxIdle              = "(可选)连接池最大闲置的连接数"
			maxOpen              = "(可选)连接池最大打开的连接数"
			maxLifetime          = "(可选)连接对象可重复使用的时间长度"
			createdAt            = "(可选)自动创建时间字段名称"
			updatedAt            = "(可选)自动更新时间字段名称"
			deletedAt            = "(可选)软删除时间字段名称"
			timeMaintainDisabled = "(可选)是否完全关闭时间更新特性，true时CreatedAt/UpdatedAt/DeletedAt都将失效"
			ctxStrict            = "(可选)是否严格限制在ORM操作时必须调用Ctx方法传递上下文变量，否则执行任何的SQL报错，默认关闭"
	1.2 link的作用
		link可用于配置的简化，MySQL的link格式（其他数据库参考最上面的goframe.org文档链接）:
		mysql:user:pass@tcp(host:port)/name (link中已配置的字段，就不需要再单独配置了)


2.单分区多节点配置实例
[database]
    # 每一个分组中可以配置多个节点，一个master(主节点)，多个slave(从节点)
	# 两个database.group1块，一个是master节点，一个是slave节点
    [[database.group1]]
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "lcoder124541"
        name = "lc-sql"
        type = "mysql"
        # link
        role = "master"
        debug = true
        prefix = "gf_"
        dryRun = 0
        charset = "utf8"
        weight = "50"
        timezone = "local"
        maxIdle = "10"
        maxOpen = "100"
        maxLifetime = "10s"
        createAt = "create_at"
        updateAt = "update_at"
        deleteAt = "delete_at"
        timeMaintainDisabled = false
        ctxStrict = false
    [[database.group1]]
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "lcoder124541"
        name = "lc-sql"
        type = "mysql"
        # link
        role = "slave"
        debug = true
        prefix = "gf_"
        dryRun = 0
        charset = "utf8"
        weight = "50"
        timezone = "local"
        maxIdle = "10"
        maxOpen = "100"
        maxLifetime = "10s"
        createAt = "create_at"
        updateAt = "update_at"
        deleteAt = "delete_at"
        timeMaintainDisabled = false
        ctxStrict = false

3.多分组实例
	3.1 支持多节点负载均衡的配置(分组两边使用[[]])
	[database]
		[[database.group1]]
			link = "mysql:root:123456@tcp(127.0.0.1:3306)/test1"
		[[database.group2]]
			link = "mysql:root:123456@tcp(127.0.0.1:3306)/test2"

	3.2 不使用多节点负载均衡的配置(分组两边使用[])
	[database]
		[database.default]
			link = "mysql:root:123456@tcp(127.0.0.1:3306)/test1"
		[database.user]
			link = "mysql:root:123456@tcp(127.0.0.1:3306)/user"

4.单数据库节点(不使用分组)
	[database]
		link = "mysql:root:123456@tcp(127.0.0.1:3306)/user"

5.日志输出配置
	# gdb的日志输出配置,当没有配置时，会使用日志组件的默认配置
	[database]
		[database.logger]
			path   = "/var/log/gf-app/sql"
			level  = "all"
			stdout = true
		# 数据库分组配置
		[database.primary]
			link = "mysql:root:123456@tcp(127.0.0.1:3306)/user"
			debug = true

	日志组件默认配置如下:
		path   = "/var/log" # 日志文件保存路径
		level  = "all" # 日志输出级别
		stdout = false # 日志是否同时输出到终端

6.gdb原生配置 TODO

7.sql.Result
	7.1 result.LastInsertId() 返回(int64, error)类型, 返回的是最后插入的数据的id值
	7.2 result.RowsAffected() 返回(int64, error)类型, 返回的是影响的数据条数,例如插入N条数据,就返回N

8.db.Model
	// Model可以传入数据库表名或结构体名称
	8.1 db.Model("t").Insert(
		// 插入数据时需要设置的字段值
		g.Map{
			"a": 1,
			"b": 2,
		}
	)
	8.2 由于config.toml配置了createAt = "create_at" ... 等,插入数据时,一开始表t的create_at和update_at字段会自动赋值为当前时间
		delete_at是软删除时间, 当数据被删除时会自动赋值; update_at是更新时间
*/

func main() {
	// 对createAt/updateAt/deleteAt的数据库配置进行测试
	// 用gf启动，则项目路径为main.go所在的路径
	g.Cfg().SetPath("../config")
	db := g.DB("g1")
	s := g.Server()
	// 操作数据库的语句最好放在接口中调用,如果放到接口外,gf run main.go启动时,每次更新文件,gf会自动热更新,每次热更新都会重新运行一次main.go
	s.BindHandler("/", func(r *ghttp.Request) {
		// 可以是表名,也可以是结构体名称
		result, err := db.Model("t").Insert(g.Map{
			"a": "first",
		})
		if err != nil {
			g.Log().Line(true).Error(err)
			return
		}
		lastId, err := result.LastInsertId()
		if err != nil {
			g.Log().Line(true).Error(err)
			return
		}
		affectedCount, err := result.RowsAffected()
		if err != nil {
			g.Log().Line(true).Error(err)
			return
		}
		g.Log().Line(true).Info("affectedCount: ", affectedCount)
		g.Log().Line(true).Info("lastInsertId: ", lastId)
	})
	s.SetPort(8303)
	s.Start()

	g.Wait()
}
