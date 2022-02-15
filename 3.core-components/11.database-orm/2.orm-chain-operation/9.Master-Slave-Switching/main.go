package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-主从切换
https://goframe.org/pages/viewpage.action?pageId=1114171

TODO:需要部署MySQL主从集群 进行验证


我们来一个简单的示例。我们有一个订单系统，每天的流量比较大，因此数据库在主从同步时往往会存在1-500ms时间的延迟。在业务需求中，
创建订单后需要立即展示订单列表页面。可以预料到如果该订单列表页面默认往从节点读取数据的话，很有可能用户在创建订单后在订单列表页面看不到最新创建的订单（因为数据库主从同步延迟）。
这个问题，我们可以在订单列表页面设置为往主节点读取数据即可解决。

一、主从切换
	1.1 配置主从节点(role="xxx" weight=50 指定权重)
	[database]
		[[database.default]]
			link = "mysql:root:Lcoder66242@tcp(47.101.48.37:3306)/lc-sql"
			# 与数据库进行数据交互时会打印sql语句到终端
			debug = true
			role = "master"
		[[database.default]]
			link = "mysql:root:Lcoder66242@tcp(47.101.48.37:3306)/lc-sql"
			# 与数据库进行数据交互时会打印sql语句到终端
			debug = true
			role = "slave"

	1.2 如果不配置主从，则默认写请求在master节点，读请求在slave节点（分摊读写，减小数据库压力）
	1.3 问题：对于有些业务，例如创建了订单之后需要立即显示订单，如果直接写master，然后在从slave节点读，因为从节点从主节点同步数据存在时间延迟
			所以可能会出现往主节点写入之后，数据还没有同步到从节点，此时再从从节点读取数据，就读取不到刚刚创建的订单了

		解决：对于特定的业务，从master节点进行读和写两个操作
		m.Master().Where(xx)  指定从master节点读取数据

		m.Slave().Where(xx) 指定对slave节点执行查询操作
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user")
		group.GET("/master-slave", func(r *ghttp.Request) {
			// 从节点读取
			res, err := m.Where("uid>0").All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("slave-read:", res)
			}

			// 主节点写入
			result, err := m.Insert(g.Map{
				"name": "xg",
			})
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("master-write: ", result)
			}

			// 主节点读取
			res, err = m.Master().Where("name", "lh").All()
			if err != nil {
			    g.Log().Line(true).Error(err.Error())
			    r.Response.Writeln(err.Error())
			} else {
			    r.Response.Writeln("master-read:", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
