package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 静态文件服务设置: https://goframe.org/pages/viewpage.action?pageId=17204011
/*
1.静态资源设置
2.多实例创建
3.多端口监听
4.多域名绑定
*/

func main() {
	// 创建一个Server实例，是一个单例模式（即多次调用该方法返回的是同一个Server对象）
	// g.Server()与g.Server("s")创建的是两个不同的Server实例
	s := g.Server("s")

	// 设置是否允许列出静态文件(通过SetServerRoot AddSearch AddStaticPath设置的)列表
	// 如果设置为false，则文件不会列出，如果访问设置的静态文件目录及其子目录或文件，则会显示Forbidden
	s.SetIndexFolder(true)

	// 静态文件服务总开关，默认是关闭的（即false），当调用SetServerRoot AddSearchPath AddStaticPath时会自动设置为true(即打开静态文件服务)
	// 如果设置为false，则所有设置的静态文件目录将均无法访问
	// s.SetFileServerEnabled(true)

	// 设置默认的静态文件查找路径（会被添加到SearchPath的第一个搜索路径）
	// SetServerRoot设置的路径下的文件，如果需要通过浏览器访问或下载，其URI与该文件的相对路径一直
	// 例如设置为"."表示当前路径（即1.hello-world目录下），则浏览器通过 http://localhost/bin/main.exe即可下载main.exe文件
	// 访问 http://localhost/bin 则会列出bin目录下的所有文件
	s.SetServerRoot(".")

	// 添加静态文件检索路径，可以调用多次，每次添加一个
	// ../（即learn-gf目录下），则通过 http://localhost/README.md 则可访问到README.md的内容
	s.AddSearchPath("../")

	// 添加URI与目录路径的映射关系，可以自定义静态文件的的URI访问规则
	// 则访问 http://localhost/home/README.md 即可访问到文件内容
	s.AddStaticPath("/home", "../")

	// 默认情况下不支持静态文件处理
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})

	// 设置多端口监听，访问以下端口将得到相同的结果
	// 注意，设置了其他端口之后，默认端口80就失效了，需要要使用80需要显示设置
	s.SetPort(8100, 8101, 8102)

	// 监听Server运行，默认监听80端口
	// 如果使用s.Run() 则不会再执行Run()之后的语句，Start则可以
	s.Start()

	// 同时也支持创建多个实例，通过传入一个name参数，如果参数名相同则为同一个实例，否则为不同的实例
	s1 := g.Server("s1")
	s1.SetPort(80)
	s1.SetIndexFolder(true)
	s1.SetServerRoot(".")
	// 这个/bin是在域名基础上加上：localhost/bin, 而与/s1这个路由无关
	// 通过 localhost:80/bin 即可列出bin目录下的文件
	s1.AddStaticPath("/bin", "/bin")

	// 如果这里不设置domain，且路由与多域名绑定的路由冲突，则该路由会覆盖多域名绑定的路由
	s1.BindHandler("/s2", func(r *ghttp.Request) {
		r.Response.Write("this is s1 server")
	})
	// 多域名绑定, 访问127.0.0.1时候调用Hello1，访问localhost时调用Hello2
	s1.Domain("127.0.0.1").BindHandler("/s1", Hello1)
	s1.Domain("localhost").BindHandler("/s1", Hello2)
	// 还可以多个域名写在一个Domain函数中，用英文逗号分开
	// 注意：Domain只支持绑定精确的域名，不支持通配符
	s1.Domain("127.0.0.2,127.0.0.3").BindHandler("/s1", Hello3)

	s1.Start()

	// s2
	s2 := g.Server("s2")
	s2.SetIndexFolder(true)
	s2.SetServerRoot(".")
	s2.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request) {
		r.Response.Writef(
			"%v %v %v %v",
			r.Get("class"),
			r.Get("course"),
			r.Get("name"),
			r.Get("act"),
		)
	})
	s2.SetPort(8199)
	s2.Start()

	// 会一直阻塞，知道所有启动的Server实例都关闭
	g.Wait()
}

func Hello1(r *ghttp.Request) {
	r.Response.Write("this domain: 127.0.0.1")
}

func Hello2(r *ghttp.Request) {
	r.Response.Write("this domain: localhost")
}

func Hello3(r *ghttp.Request) {
	r.Response.Write("this domain: 127.0.0.2, 127.0.0.3")
}
