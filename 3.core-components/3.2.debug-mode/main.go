package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*调试模式
1.简介
	在一些关键节点打印出带有[INTE]前缀的消息，只能输出到标准输出，不能重定向到文件
2.打开调试模式的3中办法
	2.1 命令行方式打开:
		2.1.1 首先编译main.go文件: gf build -o bin/main.exe main.go；会在bin/路径下生成main.exe文件
		2.1.2 命令行通过: ./bin/main.exe --gf.debug=true 或者 ./bin/main.exe --gf.debug=1开启
		2.1.3 通过g.SetDebug(true)打开调试模式
		2.1.4 通过设置环境变量: GF_DEBUG=true 或者 GF_DEBUG=1设置(edit configurations)
*/

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s.SetPort(8199)
	s.Start()
	g.Wait()
}
