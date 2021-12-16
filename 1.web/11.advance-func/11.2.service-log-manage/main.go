package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

/*
1.日志配置属性
	1.1 首先日志分为三大类：access log/error log/WebServer log
			access log表示访问日志(默认是关闭的，通过AccessLogEnabled配置项打开)
			error log表示错误日志
			WebServer log 表示业务日志（即类似于nginx输出的server日志一样，通过glog或g.Log()来打印）
	1.2 默认情况下的配置
		日志直接输出到终端，不会输出到文件，access log默认关闭；仅有error log默认开启

	1.3 属性名
	Logger            *glog.Logger      // g.Log()和glog返回的对象
	LogPath           string            // 存放日志文件的根目录
	LogStdout         bool              // 是否打印Server日志(即g.Log()或glog打印的部分)到终端
	ErrorStack        bool              // Logging stack information when error.
	ErrorLogEnabled   bool              // Enable error logging files.
	ErrorLogPattern   string            // Error log file pattern like: error-{Ymd}.log
	AccessLogEnabled  bool              // 是否打印请求日志，默认是false，默认只打印到终端，如果设置了AccessLogPattern则会打印到文件
	AccessLogPattern  string            // 打印请求日志到文件: access-{Ymd}.log，{Ymd}用于设置时间格式（需要加大括号，否则就是普通字符）

2.s.Setxxxx
	2.1 s.SetAccessLogEnabled(true)可以设置开启请求日志的打印
	2.2 需要注意的是，s实例没有SetAccessLogPattern和SetErrorLogPattern方法，所以需要根据SetConfigWithMap或config.toml配置

3.请求日志
	3.1 日志格式: 访问时间(精确到毫秒) HTTP状态码 "请求方式 请求前缀 请求地址 请求协议” 执行时间(毫秒) 客户端IP 来源URL UserAgent。
		其中，请求前缀为http或者https，请求协议往往为HTTP/1.0或者HTTP/1.1。

4.Server日志
	4.1 Server日志需要通过g.Log()或glog进行配置，如果需要设置Server日志输出到文件，首先需要关闭输出到终端的开关
		glog.SetStdoutPrint(false)
		glog.SetPath("./log")
		glog.SetFile("glog.{Y-m-d}.log")
		glog.Print("hello") // 有了上面三行配置，就会将hello输出到/log/glog.20210916-09-16.log中

5.自定义log
	5.1 通过添加需要的日志配置项创建新的*glog.Logger对象
	5.2 通过添加中间件，例如设置一个前置中间件来打印请求日志，可以自定义请求日志的格式（例如值打印状态码，url等信息）
*/

func main() {
	s := g.Server()
	// 请求日志默认是关闭的，通过SetAccessLogEnabled开启请求日志，开启后默认是打印到终端（而不是文件）
	s.SetAccessLogEnabled(true)
	s.SetConfigWithMap(g.Map{
		"LogPath":          "./log",
		"AccessLogPattern": "access.{Ymd}.log",
	})
	s.BindHandler("/", func(r *ghttp.Request) {
		glog.SetStdoutPrint(false)
		glog.SetPath("./log")
		glog.SetFile("glog.{Y-m-d}.log")
		glog.Println("hello")
		r.Response.Writeln("ok")
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
