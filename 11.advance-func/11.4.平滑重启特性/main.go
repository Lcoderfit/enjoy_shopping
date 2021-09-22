package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gproc"
	"time"
)

/*
1.gproc.Pid()
	1.1 返回当前进程的PID

2.平滑重启（热重启）和平滑更新（热更新）
	2.1 平滑重启：用程序A覆盖程序B，平滑过渡，不会中断当前正在处理的请求
	2.2 平滑更新：修改程序后自动重新编译运行
	2.3 要开启平滑重启功能，需要添加两个配置项
	s.SetConfigWithMap(g.Map{
		"graceful": true,
		"gracefulTimeout": 2,
	})

3.s.Restart() 和 s.EnableAdmin s.ShutDown
	3.1 函数签名
		func (s *Server) Restart(newExeFilePath...string) error 参数可以指定程序重启时所运行的可执行文件路径
		func (s *Server) Shutdown() error	关闭服务
		func (s *Server) EnableAdmin(pattern ...string) 为用户提供了一个管理WebServer的界面
				（界面有restart和shutdown两个功能链接，点击可启动相关服务）;参数也可以指定程序重启时所运行的可执行文件路径
			注意：EnableAdmin提供的restart接口可以通过url传参的方式控制重启时运行的可执行文件路径，
			localhost:8300/debug/admin?newExeFilePath=xxxxxxx

4.支持HTTPS服务的平滑重启
5.支持多服务和多端口重启
	s2 := g.Server("s2")
	s2.EnableAdmin()
	s2.SetPort(8302, 8303)
	s2.Start()

6.在windows和Linux系统上的区别
	在windows系统上只能进行完整重启(先杀掉后重启)
	Linux上是平滑重启

7.通过命令进行平滑重启 和 平滑关闭
	7.1 平滑重启: kill -SIGUSR1 pid
	// SIGINT/SIGQUIT/SIGKILL/SIGTERM/SIGHUB
	7.2 平滑关闭： kill -SIGTERM pid
*/

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"graceful":        true,
		"gracefulTimeout": 2,
	})
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello")
	})
	s.BindHandler("/pid", func(r *ghttp.Request) {
		r.Response.Writeln(gproc.Pid())
	})
	s.BindHandler("/sleep", func(r *ghttp.Request) {
		r.Response.Writeln(gproc.Pid())
		time.Sleep(10 * time.Second)
		r.Response.Writeln(gproc.Pid())
	})
	s.EnableAdmin()
	s.SetPort(8300)
	s.Start()

	// HTTPS平滑重启
	s1 := g.Server("s1")
	s1.SetConfigWithMap(g.Map{
		"graceful":        true,
		"gracefulTimeout": 2,
	})
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello")
	})
	s1.EnableHTTPS(
		"/etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem",
		"/etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem",
	)
	s1.EnableAdmin()
	s1.Start()

	// 多服务和多端口平滑重启
	s2 := g.Server("s2")
	s2.EnableAdmin()
	s2.SetPort(8302, 8303)
	s2.Start()

	g.Wait()
}
