package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
1.r.Session文件存储（使用StorageFile对象实现，gcache进程缓存模块控制数据过期）
	1.1 Session存储使用 内存+文件 的方式，数据操作完全基于内存，然后通过文件存储对数据进行持久化管理
	1.2 仅当更新Session的时候（此时会被标记为dirty），才会执行Session的序列化操作并存放到文件中进行持久化的存储
	1.3 当内存中的Session不存在时，才会对文件中数据进行反序列化恢复Session到内存中（这样可以减少IO）
		例如当你启动服务并设置session之后，突然服务异常中断了（假设session没有过期），此时内存中的Session数据肯定都失效了；
		但是文件持久化功能 会将数据恢复到内存中

	1.4 比较适用于读多写少的场景（因为每次写操作会更新session，也就会导致持久化操作，即IO操作会增加）

	1.5 降低IO开销的设计
		写操作，每次更新会立即更新Session的ttl时间
		读操作，每间隔一分钟才会更新 前一分钟内读取操作所对应的Session的ttl时间，（是否只针对存活时间大于1分钟的才是间隔1分钟更新一次???）
		注意：每次访问Session后，ttl会重新计数?????(续活)

	1.6 r.Session.Set(key, value)
		r.Session.Get(key)	返回interface{}类型
		r.Session.Map() 返回map[string]interface{}类型
		r.Session.Clear() 清空Session

	1.7 gtime.Timestamp(): 获取当前时间戳，例如：1631429927
*/

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge": time.Second * 10,
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/set", func(r *ghttp.Request) {
			r.Session.Set("time", gtime.Timestamp())
			r.Response.Write("ok")
		})
		group.ALL("/get", func(r *ghttp.Request) {
			r.Response.Writeln(r.Session.Map())
		})
		group.ALL("/del", func(r *ghttp.Request) {
			r.Session.Clear()
			r.Response.Write("ok")
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
