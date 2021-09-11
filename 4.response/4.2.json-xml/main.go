package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.WriteJson开头的方法接收任意类型的参数(字节数组，整型。。。等任意类型均可)，返回json格式数据
	Content-Type为application/json

HTTP/1.1 200 OK
Content-Type: application/json
Server: GF HTTP Server
Date: Sat, 11 Sep 2021 02:50:10 GMT
Content-Length: 22


2.WriteXml开头的方法，接收任意类型的参数，返回xml格式数据
	Content-Type为application/xml

curl -i "http://localhost:8200/xml"
HTTP/1.1 200 OK
Content-Type: application/xml
Server: GF HTTP Server
Date: Sat, 11 Sep 2021 02:59:10 GMT
Content-Length: 38

<doc><id>1</id><name>john</name></doc>

3.curl -i http://127.0.0.1:8199/json
-i参数表示返回的响应中携带头部信息



4.json与jsonp的区别
json是一种数据格式，而jsonp是传递这种数据格式的方式

jsonp请求可以携带一个callback参数，例如http://localhost:8200/jsonp?callback=MyCallback，
这个例子中callback参数值为MyCallback, 该请求会返回json数据{"id": 1, "name": "john"}，但是由于
加了一个callback参数，返回结果为： MyCallback({"id": 1, "name": "john"})，

即服务端会将这个callback参数作为函数名来包裹JSON数据，这样客户端就可以随意定制自己的函数来自动处理返回数据

5.带有Exit的响应处理方法，表示退出当前请求流程(相当于return)，执行下一流程
r.Response.WriteJson
r.Response.WriteJsonExit

r.Response.WriteXml
r.Response.WriteXmlExit

r.Response.WriteJsonP
r.Response.WriteJsonPExit
*/

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/json", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{
				"id":   1,
				"name": "john",
			})
		})
		group.ALL("/jsonp", func(r *ghttp.Request) {
			r.Response.WriteJsonP(g.Map{
				"id":   1,
				"name": "john",
			})
		})
		group.ALL("/xml", func(r *ghttp.Request) {
			r.Response.WriteXml(g.Map{
				"id":   1,
				"name": "john",
			})
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
