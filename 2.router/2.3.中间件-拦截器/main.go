package main

import "github.com/gogf/gf/net/ghttp"

/*
1.通过r.Middleware属性控制请求流程
	r.Middleware.Next()跳出当前流程，执行下一流程
2.中间件类型
	前置中间件：先执行中间件逻辑，然后r.Middleware.Next()跳到下一个流程
	后置中间件：先执行r.Middleware.Next()跳到下一流程，然后再执行中间件逻辑

3.全局中间件和分组路由中间件
	全局中间件：通过s.Use()注册，遵循路由匹配规则？？？
	分组路由中间件：绑定到分组路由，按照注册的顺序先后执行；执行分组路由下的所有服务接口前都会先执行该中间件
*/

func MiddlewareCORS(r *ghttp.Request) {
	// CORSDefault设置允许任意的跨域请求
	r.Response.CORSDefault()
	r.Middleware.Next()
	r.Exit()
}

func main() {

}
