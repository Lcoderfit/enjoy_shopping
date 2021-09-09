package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

/*
1.接口类型与其他内置类型比较
	r.Get("name") == "123" 是可以的，r.Get()返回interface{}类型，interface{}类型具有动态类型和动态值，当其动态类型和动态值与"123"相等
	时，则语句为true

2.s.SetPort()可以写到路由注册的前面，但是s.Start()不行,

3.缓冲控制
	3.1. func (r *Response) Buffer() []byte  将r.buffer（缓冲区）中的数据以[]byte的形式返回
	func (r *Response) BufferLength() int	返回r.buffer中内容的大小, 单位是字节
	func (r *Response) BufferString() string 以字符串形式返回r.buffer中的内容
	func (r *Response) Flush()				先将r.buffer中的内容写入r.writer，然后情况缓冲区(r.buffer) (即将缓冲区数据输出到客户端然后情况缓冲区)
	func (r *Response) SetBuffer(data []byte)  先清空缓冲区，然后将data写入缓冲区
	func (r *Response) ClearBuffer() 直接调用r.buffer.Reset()清空缓冲区
 */

// 后置中间件，都异常进行统一处理
func MiddlewareErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.Status >= http.StatusInternalServerError {
		r.Response.ClearBuffer()
		r.Response.Writeln("服务器开了小差，请稍候重试")
	}
}

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareErrorHandler)
		group.ALL("/", func(r *ghttp.Request) {
			panic("db error: sql is xxxxxxxx")
		})
	})
	s.Start()

	g.Wait()
}
