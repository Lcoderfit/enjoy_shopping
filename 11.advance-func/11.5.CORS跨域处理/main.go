package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

/*
1.CORS跨域请求处理(注意：如果不设置跨域，则默认是不允许跨域请求的，如果设置为默认的跨域，即CORSDefault()，情况如下)
	1.1 func (r *Response) DefaultCORSOptions() CORSOptions
		返回默认的CORS跨域设置(5个, AllowOrigin, AllowMethods, AllowCredentials,AllowHeaders,MaxAge)
		AllowOrigin:      "*",	//允许进行跨域请求的域名
		AllowMethods:     supportedHttpMethods,	// 允许的HTTP请求方法
		AllowCredentials: "true", // 是否允许后续请求携带认证信息（cookies），只能是true，否则响应头不会返回该字段
		AllowHeaders:     defaultAllowHeaders, // 跨域请求时允许的请求头字段
		MaxAge:           3628800,	// 预检结果缓存时间（在跨域请求时会先发送Options预检请求，得到指定uri所支持的HTTP请求方法，然后将结果缓存）

	1.2 func (r *Response) CORSDefault()
		使用默认的CORS跨域设置处理跨域请求
	1.3 func (r *Response) CORS(options CORSOptions)
		接收自定义的跨域设置作为参数，根据自定义跨域设置处理跨域请求
	1.4 func (r *Response) CORSAllowedOrigin(options CORSOptions) bool
		传入自定义跨域设置作为参数，判断该当前请求是否被允许跨域

2.预检请求
	2.1 简单请求：https://developer.mozilla.org/zh-CN/docs/Glossary/Simple_header
		简单头部为：Accept, Accept-Language, Content-Language,
		Content-Type(且值为text/plain, application/x-www-form-urlencoded, multipart/form-data三者之一)
		当只包含简单头部时，该请求则为简单请求且不会发送预检请求(preflight request)

	2.2 预检请求流程: https://developer.mozilla.org/zh-CN/docs/Glossary/Preflight_request
		CORS预检请求用于检查服务器是否支持CORS（跨域资源共享）
		首先如果该请求为非简单请求，则浏览器会先发送一个预检请求（一般会自动发）,预检请求是一个Options请求，一般包含
			Access-Control-Request-Method // 通知服务器在正式请求中将会使用的HTTP method
			Access-Control-Request-Headers // 通知服务器在正式请求中将会使用的请求头
			Origin	// 指示当前请求来源于哪一个站点
		三个首部，然后服务端会在响应中使用:
			Access-Control-Allow-Method  // 指明了实际请求中允许使用的HTTP方法
			Access-Control-Allow-Headers // 指明了实际请求中允许携带的首部字段

3.跨域请求配置项
	3.1 AllowDomain      []string // Used for allowing requests from custom domains
		允许进行跨域请求的域名
	3.2 ExposeHeaders
		该选项设置Access-Control-Expose-Headers头部
		XMLHttpRequest对象的getResponseHeader()方法只能获取一些基本的响应头
		（Cache-Control, Content-Type, Content-Language, Expires, Last-Modified, Praga）
		而Access-Control-Expose-Headers头部用于设置可被客户端访问的头部字段

		AllowDomain      []string // Used for allowing requests from custom domains
		AllowOrigin      string   // Access-Control-Allow-Origin
		AllowCredentials string   // Access-Control-Allow-Credentials
		ExposeHeaders    string   // Access-Control-Expose-Headers
		MaxAge           int      // Access-Control-Max-Age
		AllowMethods     string   // Access-Control-Allow-Methods
		AllowHeaders     string   // Access-Control-Allow-Headers

4.Chrome调试工具发送ajax请求
	4.1 访问baidu.com, 然后按F12-> console->输入
	$.post("http://localhost:8300/api.v1/order", function(result){
		console.log(result)
	});

	如果无法进行跨域则会报错

5.AllowOrigin 和AllowDomain
	5.1 AllowOrigin配置项可以设置Access-Control-Allow-Origin响应首部字段，表示允许跨域请求的域名；
		如果不设置AllowDomain，并且使用默认的CORSOptions，则会先对AllowOrigin初始化为"*"；
		如果请求中不包含Origin字段，则根据Referer字段获取域名，如果Referer也为空，则返回的响应中Access-Control-Allow-Origin字段为"*"
		如果请求中包含Origin字段，则Access-Control-Allow-Origin即为Origin的值
	5.2 AllowDomain
		该选项主要用于对跨域设置允许进行跨域请求的域名；
		CORSAllowedOrigin会根据CORS配置项对跨域请求进行检查，先判断AllowDomain是否为空，为空则返回true；然后判断Origin是否为空，为空则返回true
		之后会遍历AllowDomain切片,如果域名解析失败,则返回false;最后判断请求中的Origin字段值是否为AllowDomain中设置的域名的子域名,只要有
		一个符合条件,即返回true,否则最终返回false

		注意：CORSAllowedOrigin会根据请求中Origin是否是AllowDomain中设置的域名的子域名来判断跨域请求是否被允许

6.自定义检测授权（如果检测不通过直接返回，不再走后续流程）
	6.1 在MiddlewareCORS1中会出现这样一个问题，即使请求无法跨域，然后后续的请求流程仍会继续（最后仍会调用服务函数打印出"order"字符串）;只不过
		r.Response.Write()的内容不会被返回给前端而已
	6.2 自定义CORS授权
		判断跨域请求是否被允许，如果不允许则直接返回，不再进行后续请求：
		if !r.Response.CORSAllowedOrigin(corsOptions) {
			r.Response.WriteStatus(http.StatusForbidden)
		}

7.也可以使用Hook函数来实现跨域，不过一般用中间件居多
*/

func Order(r *ghttp.Request) {
	g.Log().Line().Println("order")
	r.Response.Write("GET")
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareCORS1(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"goframe.org"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

func MiddlewareCORS2(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"goframe.org"}
	if !r.Response.CORSAllowedOrigin(corsOptions) {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

func MiddlewareTest(r *ghttp.Request) {
	g.Log().Line(true).Println("test middleware")
	r.Middleware.Next()
}

func main() {
	s := g.Server()
	s.Group("/api.v1", func(group *ghttp.RouterGroup) {
		// 只有在中间件注册之后注册的路由该中间件才会执行，在中间件注册之前注册的服务函数执行时，是不会执行该中间件的
		group.ALL("/test", func(r *ghttp.Request) {
			g.Log().Println("test url")
		})
		//group.Middleware(MiddlewareTest, MiddlewareCORS)
		//group.Middleware(MiddlewareCORS)
		//group.Middleware(MiddlewareCORS1)
		group.Middleware(MiddlewareCORS2)
		group.POST("/order", Order)
	})
	s.SetPort(8300)
	s.Start()

	g.Wait()
}
