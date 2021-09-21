package main

import (
	"github.com/gogf/csrf"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
	"time"
)

/*
1.CSRF（跨站脚本攻击，cross site request forgery）原理
	1.1 首先用户A登录网站B，然后网站B会返回cookie信息到浏览器；
		当用户A访问恶意网站C时，恶意网站C可以自动发送一个请求到网站B，并且由于此时浏览器保留了之前用户A访问网站B时返回的cookie信息
		所以此时恶意网站自动发送的请求也可以携带该cookie信息（服务器会认为是A发送的请求），从而导致并非用户A本人意愿的操作
	1.2 防御手段：https://segmentfault.com/q/1010000004957432
	对提交的表单携带上token信息(添加一个隐藏的输入框,例如<input type="hidden" name="X-Token" value="{{.token}}">),
	而恶意网站他无法获取这个token,所以无法发送带有token的表单数据(所以网站B会拒绝请求)

	1.3 疑问：如果恶意网站先发送get请求到网站B，获取到响应的html页面中的token，然后再发送POST请求，不久仍然可以csrf成功吗？

2.默认的CSRF设置
	DefaultCSRFConfig = Config{
		// Name设置存放到Cookie中的token字段名
		Cookie: &http.Cookie{
			Name: "_csrf",
		},
		// Cookie有效期
		ExpireTime:      time.Hour * 24,
		// Token长度(32位数字和大小写字母组合)
		TokenLength:     32,
		// 该设置用于从请求中通过该字段获取token值(r.Get("X-CSRF-Token"))
		// 例如： <input type="hidden" name="X-Token" value="{{.token}}">，后端可以通过r.Get("X-Token")获取value（即token值）
		// 然后再将该token与服务器端的token进行比对，如果一致则通过校验
		TokenRequestKey: "X-CSRF-Token",
	}
	首先，第一次请求时，响应的Set-Cookie头会携带token信息，例如：
	Set-Cookie:_csrf=aycUQu77MsZVi7yWSrbjJWDs48ZaWCCW; Expires=Mon, 20 Sep 2021 08:52:11 GMT; Secure; SameSite=None
	然后浏览器会保存该cookie信息，之后的请求中的Cookie头部将包含此cookie信息，例如：
	Cookie:_csrf=AVUxizHfx0IYrFHrOKjQWDB16MXtQoOe

2.r.Cookie.Set()和csrf.New()或csrf.NewWithCfg()
	2.1 r.Cookie.Set是在*ghttp.Request实例的data属性（即r.data，一个字典）中添加key value键值对
		r.Cookie.Get也是从r.data中获取值

	2.2 csrf.New()或csrf.NewWithCfg()
		在响应中(r.Response)添加Set-Cookie头部并将token键值对赋值给该头部
		group.Middleware(csrf.NewWithCfg(csrf.Config{
			TokenLength:     32,
			// 可携带在请求中的带有token信息的字段名
			TokenRequestKey: "X-Token",
			ExpireTime:      time.Hour * 24,
			Cookie: &http.Cookie{
				Name:     "_csrf",
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			},
		}))



客户端：get请求不允许修改数据
服务端：csrf_token

*/

func GetResponseCookie(r *ghttp.Request, key string) string {
	rsp := http.Response{Header: r.Response.Header()}
	for _, v := range rsp.Cookies() {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(csrf.New())
		group.ALL("/csrf", func(r *ghttp.Request) {
			r.Response.Writeln(r.Method + ":" + r.RequestURI)
		})
	})
	s.SetPort(8300)
	s.Start()

	s1 := g.Server("s1")
	s1.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(csrf.NewWithCfg(csrf.Config{
			TokenLength:     32,
			TokenRequestKey: "X-Token",
			ExpireTime:      time.Hour * 24,
			Cookie: &http.Cookie{
				Name:     "_csrf",
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			},
		}))
		group.ALL("/csrf", func(r *ghttp.Request) {
			r.Response.Writeln(r.Method + ":" + r.RequestURI)
		})
		group.ALL("/index", func(r *ghttp.Request) {
			r.Response.WriteTpl("index.html", g.Map{
				"token": GetResponseCookie(r, "_csrf"),
			})
		})
	})
	s1.SetPort(8301)
	s1.Start()

	s2 := g.Server("s2")
	s2.BindHandler("/", func(r *ghttp.Request) {
		r.Cookie.Set("a", "b")
		g.Log().Line(true).Println(r.Cookie.Get("a"))
		r.Response.Writeln(r.Cookie.Get("a"))
	})
	s2.SetPort(8302)
	s2.Start()

	g.Wait()
}
