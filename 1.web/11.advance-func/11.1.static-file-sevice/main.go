package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.静态文件管理
	s.SetFileServerEnabled(false) 静态文件服务总开关，默认是false，如果设置了SetServerRoot AddSearchPath AddStaticPath
							则会自动设置为true，如果设置为false，则所有静态文件均无法访问
			注意：SetFileServerEnabled必须要在所有的(SetServerRoot AddSearchPath AddStaticPath之后设置才有效)
			如果最终静态服务总开关设置为false，则访问静态资源时将显示网页"无法访问此网站"
		例如：
			s.SetServerRoot(".")
			s.SetFileServerEnabled(false)
			// 由于这里设置了静态文件路径，所以静态服务总开关又会重新自动设置为true，
			// 只要为true，则上面那个SetServerRoot(".")的设置也将生效
			s.AddSearchPath("../../")

	s.SetIndexFolder(true) 表示是否允许展示静态目录的文件列表（通过SetServerRoot AddSearchPath AddStaticPath），
							默认是false，如果设置为false，则访问设置的静态文件及其子目录，会显示Forbidden
	s.SetServerRoot("./") 设置默认的静态文件查找路径，访问localhost:8200/时会默认优先使用SetServerRoot设置的路径
	s.AddSearchPath("../../") 添加静态文件搜索路径，可以添加多个，但参数只能传一个路径（不要误认为是不定参）
	s.AddStaticPath(uri, directory) 添加uri与静态目录的映射，
					注意：是与静态目录的映射，如果设置为文件的路径，则浏览器会显示Not Found
	s.

2.路由注册优先级
	2.1 服务函数注册的路由优先级要高于静态服务设置的路由优先级（与注册顺序无关），所以请求/index会调用服务函数返回"BindHandler index"
	2.2 通过SetRewrite/SetRewriteMap重写的uri会覆盖服务函数注册的uri，具有最高优先级，如果是前后注册了好几个相同的uri，则后面的会覆盖前面的

	注意:这个重写太恶心了,
	// 必须要设置SetServerRoot才能使SetRewrite有效，否则设置将无效（而且SetServerRoot必须要是当前目录才行，否则也无效！！！）
	// 这个特性太恶心了
	s4.SetServerRoot(".")
	s4.AddSearchPath("../")
	// 下面这一行是无效的,因为映射的路径在当前目录之外
	s4.SetRewrite("/t", "../")
	s4.SetRewriteMap(g.MapStrStr{
		"/t1": "t1.html",
		"/t2": "t2.html",
	})

		// SetRewrite/SetRewriteMap具有最高优先级，会覆盖服务函数注册的路由
		s2.SetRewrite("/index", "./main.go")
		// 如果重写的uri相同，则后面的uri会覆盖前面的uri映射关系
		s2.SetRewrite("/index", "./bin/main.exe")

		// uri可以映射到目录，如果SetIndexFolder为true则会展示静态目录下的文件列表
		// 但有一个奇怪的现象，如果你在浏览器上点击文件链接，即使点击main.exe,也会继续访问 xxxx/main.exe的url(浏览器显示Not Found)，
		// 并不会下载main.exe文件，而且不像SetSerRoot那三个函数一样（这三个函数设置的路径，可以一直访问目录下的子目录），
		// 但是SetRewrite如果设置的是uri映射到目录，则只能显示当前目录下的文件列表，再访问子目录的话，由于添加了子目录名的url没有注册，所以会显示Not Found
		// SetServerRoot AddSearchPath AddStaticPath是可以下载的
		s2.SetRewrite("/index/bin", ".")

		// 如果uri映射到main.go,则会显示，如果是main.exe，则会下载，不过是以uri最后一段的名字作为下载的文件名
		s2.SetRewrite("/index/a/b", "./bin/main.exe")
		// 无法将uri映射到当前目录之外
		s2.SetRewrite("/index1", "../")

	2.3 静态目录映射的优先级（AddStaticPath）
		s3.AddStaticPath("/index", ".")
		s3.AddStaticPath("/index/bin", "../")

		// 如果./目录下有bin目录，则/index/bin会访问到../，永远无法访问到./bin
		// 但是请求/index时还是可以看到文件列表，如果点击/bin对应的链接，则会跳转到/index/bin对应的目录

	2.4 g.MapStrStr g.MapStrInt .... （内置map，value值为各种类型的都有）

*/

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.SetServerRoot("./")
	s.AddSearchPath("../../")
	// 必须设置在SetServerRoot AddSearchPath AddStaticPath之后才有效
	// s.SetFileServerEnabled(false)
	s.SetPort(8207)
	s.Start()

	s1 := g.Server("s1")
	s1.SetIndexFolder(true)
	s1.SetServerRoot(".")
	s1.AddSearchPath("../../")
	s1.AddStaticPath("/index/bin", "../")
	// 服务函数注册的路由优先级要高于静态服务设置的路由优先级（与注册顺序无关），所以请求/index会调用服务函数返回"BindHandler index"
	s1.BindHandler("/index", func(r *ghttp.Request) {
		r.Response.Writeln("BindHandler index")
	})
	s1.AddStaticPath("/index", "./")
	s1.SetPort(8208)
	s1.Start()

	s2 := g.Server("s2")
	s2.SetIndexFolder(true)
	s2.SetServerRoot(".")
	s2.AddSearchPath("../../")
	//s2.AddStaticPath("/index", "..")

	// SetRewrite/SetRewriteMap具有最高优先级，会覆盖服务函数注册的路由
	s2.SetRewrite("/index", "./main.go")
	// 如果重写的uri相同，则后面的uri会覆盖前面的uri映射关系
	s2.SetRewrite("/index", "./bin/main.exe")
	// uri可以映射到目录，如果SetIndexFolder为true则会展示静态目录下的文件列表
	// 但有一个奇怪的现象，如果你在浏览器上点击文件链接，即使点击main.exe,也会继续访问 xxxx/main.exe的url(浏览器显示Not Found)，
	// 并不会下载main.exe文件，而且不像SetSerRoot那三个函数一样（这三个函数设置的路径，可以一直访问目录下的子目录），
	// 但是SetRewrite如果设置的是uri映射到目录，则只能显示当前目录下的文件列表，再访问子目录的话，由于添加了子目录名的url没有注册，所以会显示Not Found
	// SetServerRoot AddSearchPath AddStaticPath是可以下载的
	s2.SetRewrite("/index/bin", ".")
	// 如果uri映射到main.go,则会显示，如果是main.exe，则会下载，不过是以uri最后一段的名字作为下载的文件名
	s2.SetRewrite("/index/a/b", "./bin/main.exe")
	// 无法将uri映射到当前目录之外
	s2.SetRewrite("/index1", "../")

	s2.BindHandler("/index", func(r *ghttp.Request) {
		r.Response.Writeln("BindHandler index")
	})
	// s2.SetRewriteMap(map[string]string{
	// 	"/index": "../11.2.服务日志管理/",
	// })
	s2.SetPort(8202)
	s2.Start()

	s3 := g.Server("s3")
	s3.SetIndexFolder(true)
	s3.AddStaticPath("/index", ".")
	// 如果./目录下有bin目录，则/index/bin会访问到../，永远服务访问到./bin
	// 但是请求/index时还是可以看到文件列表，如果点击/bin对应的链接，则会跳转到/index/bin对应的目录
	s3.AddStaticPath("/index/bin", "../")
	s3.SetPort(8203)
	s3.Start()

	s4 := g.Server("s4")
	s4.SetIndexFolder(true)
	// 必须要设置SetServerRoot才能使SetRewrite有效，否则设置将无效（而且SetServerRoot必须要是当前目录才行，否则也无效！！！）
	// 这个特性太恶心了
	s4.SetServerRoot(".")
	s4.AddSearchPath("../")
	s4.SetRewrite("/t", "../")
	s4.SetRewriteMap(g.MapStrStr{
		"/t1": "t1.html",
		"/t2": "t2.html",
	})
	s4.SetPort(8204)
	s4.Start()

	g.Wait()
}
