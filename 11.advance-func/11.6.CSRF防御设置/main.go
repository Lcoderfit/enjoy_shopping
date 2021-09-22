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

	2.3 Cookie设置（尤其是SameSite属性：https://zhuanlan.zhihu.com/p/121048298）
		2.3.1 http.Cookie配置
		http.Cookie{
			// 在Set-Cookie或Cookie字段中携带的token信息的键
			Name: "_csrf",
			Secure: true,
			SameSite: http.SameSiteNoneMode,
		}

		2.3.2 SameSite配置
			2.3.2.1 SameSite三种配置
				SameSite用于限制第三方Cookie，有三种取值：Strict，Lax，None(对应http.SameSiteStrictMode http.SameSiteLaxMode, http.SameSiteNoneMode)
				http.SameSiteStrictMode就相当于设置Cookie的SameSite属性值为Strict,例如: Set-Cookie:xxx;SameSite=Strict;xxx

			2.3.2.2 SameSiteDefaultMode
				http.SameSiteDefaultMode会设置Cookie的SameSite属性，但是不会赋值(例如：Set-Cookie:	xxx;SameSite;xxx)，即表示取Chrome默认的SameSite设置
					注意：同站和同域的区别：具有相同的二级域名则为同站（eTLD+1），例如a.lcoderfit.com和b.lcoderfit.com的二级域名均为lcoderfit.com，所以属于同站
						但是a.lcoderfit.com和b.lcoderfit.com属于不同域

			2.3.2.3 三种设置的含义(注意：无论是何种设置，同站的情况下均会发送cookie)：
			https://zhuanlan.zhihu.com/p/121048298
			http://www.ruanyifeng.com/blog/2019/09/cookie-samesite.html
				Strict：表示禁止所有第三方Cookie（只有同站的Cookie才被允许）
				Lax：允许发送部分第三方Cookie（当不用站时,post请求、iframe、ajax、image均不会携带cookie）
				None:无论是否跨站均会发送Cookie,由于Chrome浏览器从80版本后默认设置由None变为Lax，所以会导致一些网站的ajax，post，iframe，
					image跨站请求无法携带第三方cookie；
					此时可以采用一种临时解决方案：强制将SameSite设置为None，然后设置Cookie的Secure属性为true（第三方cookie只能通过https发送，否则无效）
					例如：跨站脚本攻击时，当用户访问恶意网站B时，网站B自动发送一个请求到A企图修改用户的数据，此时若该请求不是https请求，则不会携带cookie
					Set-Cookie:SameSite=None; Secure
					注意：要是SameSitem=None设置有效，则必须设置Secure，否则该设置无效



3.域名，同站和同源
	3.1 域名
		3.1.1 根域名 .
		3.1.2 顶级域名（一级域名）: .com .cn .org
		3.1.3 二级域名: a.com b.com c.com

	3.2 同源（SameOrigin）
		3.2.1 协议+域名+端口均相同即为同源，否则不同源

	3.3 同站(SameSite)
		3.3.1 site(站)的定义：eTLD+1， 例如https://my-project.github.io 中,github.io即为eTLD，eTLD+1为 my-project.github.io（即为站）
							其实site可以理解为TLD/eTLD+1 (即二级域名)
		3.3.2 TLD和eTLD
			TLD(top-level domain), 顶级域名，如.com .cn
			eTLD（有效顶级域名，effective top-level domain），简单来说就一个等效于一个顶级域名的二级域名
				例如有一个.io的顶级域名，github公司注册了github.io，然后该公司想将该
				域名开放，可以根据github.io注册a.github.io或b.github.io，
				此时如果仅通过TLD+1来作为站，则a.github.io和b.github.io属于同一个站，
				但是github.io和a.github.io与b.github.io属于不同的网站，
				所以需要将github.io看作一个整体，那么a.(github.io)和b.(github.io)就是两个不同的网站了，看作整体的这一部分就是eTLD

		3.3.3 如何确定是否同站？
			eTLD+1相同即为同站(也可以理解为具有相同的二级域名即为同站)

	3.4 schemeful same-site
		当协议不同时也属于不同站，例如 http://a.lcoderfit.com和https://a.lcoderfit.com是不同站

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
