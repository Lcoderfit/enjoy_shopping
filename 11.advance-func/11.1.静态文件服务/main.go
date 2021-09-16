package main

import "github.com/gogf/gf/frame/g"

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

2.
*/

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.SetServerRoot("./")
	s.AddSearchPath("../../")
	// 必须设置在SetServerRoot AddSearchPath AddStaticPath之后才有效
	// s.SetFileServerEnabled(false)
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("s1")
	s1.SetIndexFolder(true)
	s1.SetServerRoot(".")
	s1.AddSearchPath("../../")
	s1.AddStaticPath("/index", "../")
	s1.AddStaticPath("/index/bin", "")
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
