package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM查询-常用操作示例

g.Model("user") 读取默认的数据库配置(default), 基于user表创建一个*Model对象
g.DB().Model("user")与上面等效
g.DB("other").Model("user")读取配置中的other数据库配置，然后基于user表生成一个*Model对象

一、in查询

二、like查询

三、min/max/avg/sum

四、count

五、distinct

六、between

七、null
*/

func main()  {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/in", func(r *ghttp.Request) {

		})

		group.GET("/like", func(r *ghttp.Request) {

		})

		group.GET("/jh", func(r *ghttp.Request) {

		})

		group.GET("/count", func(r *ghttp.Request) {

		})

		group.GET("/distinct", func(r *ghttp.Request) {

		})

		group.GET("/between", func(r *ghttp.Request) {

		})

		group.GET("/null", func(r *ghttp.Request) {

		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}