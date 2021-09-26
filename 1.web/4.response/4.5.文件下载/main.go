package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.r.Response.ServeFile()
	如果第一个参数是文件名，则访问url会显示出文件内容；
	如果是目录，需要传入第二个参数，即allowIndex，allowIndex为true则表示允许展示目录下的文件，
	默认为false（访问url时候显示Forbidden）
	注意：即使展示了目录下的文件，但是如果文件没有对应的访问url，则你在页面点击时会显示not found
	例如a目录下有b.txt文件，你设置为：
		s.BindHandler("/", func(r *ghttp.Request){
			r.Response.ServeFile("a/", true)
		})
	则你访问http://localhost:8200/时，会列出a目录下的b.txt文件（每个文件都对应一个链接），当你点击
	这b.txt这个链接时，url跳转到:http://localhost:8200/b.txt，而由于你没有注册这个url，所以前端显示not found

2.r.Response.ServeFileDownload("test.txt", "a.md")，第一个参数为文件的相对或绝对路径，第二个参数为下载的文件的名称，
	如果不设置第二个参数，则默认是文件名，该方法用于下载文件，例如这里是下载test.txt文件，但是由于这里设置了下载时的名称，
	所以点击链接时会下载一个名叫a.md的文件（如果设置为a.zip,也可以下载，但是貌似打开后是个空的文件....）

	注意：如果第一个参数是一个目录，则通过url访问时会显示"无法访问此网站";
*/

func main() {
	s := g.Server()
	// 显示文件内容，或展示目录下的文件
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.ServeFile("test.txt")
	})
	s.BindHandler("/show-doc", func(r *ghttp.Request) {
		r.Response.ServeFile("../", true)
	})

	// 下载文件
	s.BindHandler("/download", func(r *ghttp.Request) {
		r.Response.ServeFileDownload("test.txt", "a.md")
	})
	// 第一个参数不能设置为一个目录，不然会显示"无法访问此网站"
	s.BindHandler("/download-doc", func(r *ghttp.Request) {
		r.Response.ServeFileDownload("../", "d.zip")
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
