package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.同名参数: k=1&k=2
	r.Get("name")方法，即可以获取自定义的动态路由参数，也可以获取url参数，但是如果出现了同名参数，后面的参数值会覆盖前面的
	例如：curl -XGET "http://localhost:8200?k=lu&k=lcoder" r.Get("name")会返回lcoder（后面的参数值会覆盖前面的参数值）
2.数组参数: k[]=v1&k[]=v2
	浏览器访问：http://localhost:8200/?k[]=v1&k[]=v2	r.Get("k")返回["v1","v2"]
3.Map参数：k[a]=m&k[b]=n，并且支持多级map：k[a][a]=m&k[a][b]=n
	浏览器访问：http://localhost:8200/?k[a]=m&k[b]=n  r.Get("k")返回 {"a":"m","b":"n"}
	浏览器访问：http://localhost:8200/?k[a][a]=m&k[a][b]=n  r.Get("k")返回{"a":{"a":"m","b":"n"}}
*/

func main() {
	s := g.Server()
	//curl -XGET "http://localhost:8200?k=lu&k=lcoder" 会返回lcoder（后面的参数值会覆盖前面的参数值）
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(r.Get("k"))
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
