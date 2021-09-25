package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.开启PProf服务性能分析(注意：pprof会对性能产生影响，所以除了分析性能以外，其他情况不开启)
	1.1 开启
	s.EnablePProf(pattern) pattern可以自定义分析结果页面的url，默认为/debug/pprof
	访问：http://localhost:8300/debug/pprof (前面有数字的那一行是可以点击的，
			例如allocs对应链接http://localhost:8300/debug/pprof/allocs)
profiles:
6	allocs
0	block
10	goroutine
6	heap
3	mutex
11	threadcreate

full goroutine stack dump

	1.2 开启对锁事件和阻塞操作的跟踪
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)

2.如果有多个不同的Server，则每个Server都需要开启EnablePProf，否则未开启的Server无法访问/debug/pprof

3.PProf性能指标
heap: 报告内存分配样本；用于监视当前和历史内存使用情况，并检查内存泄漏。
threadcreate: 报告了导致创建新OS线程的程序部分。
goroutine: 报告所有当前goroutine的堆栈跟踪。
block: 显示goroutine在哪里阻塞同步原语（包括计时器通道）的等待。默认情况下未启用，需要手动调用runtime.SetBlockProfileRate启用。
mutex: 报告锁竞争。默认情况下未启用，需要手动调用runtime.SetMutexProfileFraction启用。

4.详细的性能分析
$go tool pprof "http://localhost:8301/debug/pprof/profile"
	4.1 pprof会进行30s左右的接口信息采集（此时webserver应该不断有流量进入，可以不停访问hello world页面做测试）
	4.2 生成性能分析报告后终端会进入$(pprof), 此时可以输入top3 查看占用cup前三的进程
	4.3 web 命令会根据pprof采集的信息生成一个网页，将分析信息转变成可视化的流程图
	4.4 CPU性能分析: go tool pprof "http://localhost:8300/debug/pprof/profile"
		heap性能分析：go tool pprof "http://localhost:8300/debug/pprof/head"
		go tool pprof url 本质上就是对该url接口进行性能分析????

5.快速开启一个独立的PProf Server
	5.1 go ghttp.StartPProfServer(8300, pattern) // pattern可以自定义访问分析页面时的url, 8300是自定义端口
	5.2 一般用于一些没有HTTP Server的常驻进程中（定时任务，GRPC）
*/

func main() {
	// 开启对锁事件的跟踪，rate参数表示平均有1/rate的事件被报告
	//runtime.SetMutexProfileFraction(1)
	// 开启对阻塞操作的跟踪
	//runtime.SetBlockProfileRate(1)

	s := g.Server()
	s.EnablePProf()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s.SetPort(8300)
	s.Start()

	// 开启一个独立的PProf Server，常用于一些没有HTTP Server的常驻进程中（定时任务，GRPC）
	go ghttp.StartPProfServer(8200)
	s1 := g.Server("s1")
	s1.EnablePProf()
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("hello world")
	})
	s1.SetPort(8301)
	s1.Start()

	g.Wait()
}
