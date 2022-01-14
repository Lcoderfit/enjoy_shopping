package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/genv"
)

/*
环境变量-genv
https://goframe.org/pages/viewpage.action?pageId=1114354

1.获取所有的环境变量
	1.1 genv.All() 返回一个[]string, 每个环境变量为key=value格式
		["A=A", "B=b"]
	1.2 genv.Map() 返回一个map[string]string类型
		{"A":"a", "B":"b"}
	1.3 genv.Get(key, def) 获取指定key的环境变量（忽略大小写）,如果没有获取到，则返回def设置的默认值
		注意：上面这三个获取环境变量的方法，设置的key是啥就是啥，除了忽略大小写外没有区别，例如:
		genv.SetMap(g.MapStrStr{"ADDR.PORT": "xx"}), 则只能通过addr.port这个key获取（小写字母可以改成大小字母）
		但是.是不能写成_的
	1.4 genv.GetWithCmd(key, def)
		1.4.1 优先从环境变量获取，且会将key中字母全部转换为大写，.转换为_
			genv.Set("gf.debug", "1") 通过 genv.GetWithCmd("gf.debug")是获取不到的,因为从环境变量获取时会将key转换为"GF_DEBUG"
			genv.Set("GF_DEBUG", "1") 通过 genv.GetWithCmd("gf.debug")是可以获取的

		1.4.2 如果没有，则从命令行中获取，会将key中字母全部转换为小写，_转换为.
			./bin/main --gf.debug 1 通过 genv.GetWithCmd("gf.debug")是可以获取到的,因为从环境变量获取时会将key转换为"gf.debug"

		1.4.3 如果仍没有，则返回def参数设置的默认值

2.设置环境变量(只能设置字符串类型的key value)
	注意：设置的key是啥就是啥，除了忽略大小写外没有区别，例如:
	genv.Set("gf.debug", "1"), 则只能通过"gf.debug"这个key获取（小写字母可以改成大小字母）,但是.是不能写成_的

	2.1 genv.Set(key, value)
	2.2 genv.SetMap(map[string]string{"A":"a", "B":"b"})

3.删除环境变量
	genv.Remove(key ...string) 可以删除一个或多个

4.判断环境变量是否存在
	genv.Contains(key)
*/

func main() {
	genv.SetMap(g.MapStrStr{
		"gf.debug":  "this is debug mode",
		"ADDR.PORT": "127.0.0.1:6379",
	})

	s := g.Server()
	s.BindHandler("/all", func(r *ghttp.Request) {
		// 返回一个[]string, 每个元素是key=value的格式
		sArr := genv.All()
		// 把数据进行json序列化以后返回
		r.Response.WriteJson(sArr)
	})

	s.BindHandler("/map", func(r *ghttp.Request) {
		m := genv.Map()
		r.Response.WriteJson(m)
	})

	s.BindHandler("/get", func(r *ghttp.Request) {
		s := genv.Get("addr.port", "127.0.0.1:8200")
		r.Response.WriteJson(s)
	})

	s.BindHandler("/set", func(r *ghttp.Request) {
		genv.Set("fuck", "shit")
		r.Response.WriteJson(genv.Get("fuck"))
	})

	s.BindHandler("/setmap", func(r *ghttp.Request) {
		genv.SetMap(g.MapStrStr{
			"A": "a",
			"B": "b",
		})
		r.Response.WriteJson(genv.Contains("A"))
	})

	s.BindHandler("/remove", func(r *ghttp.Request) {
		genv.Remove("A")
		r.Response.Writeln(genv.Contains("A"))
		r.Response.Writeln(genv.Build(g.MapStrStr{
			"C": "c",
			"D": "d",
		}))
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
