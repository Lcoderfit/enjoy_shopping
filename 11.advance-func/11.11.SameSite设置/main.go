package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

/*
1.SameSite等基础知识
https://web.dev/samesite-cookies-explained/

https://web.dev/samesite-cookie-recipes/

https://web.dev/schemeful-samesite/

Chrome89开始协议不同也属于跨站请求（跨站跟跨域是不同的，参见11.6CSRF防御设置中的注释，同站和同域的区别）
https://www.chromestatus.com/feature/5096179480133632

2.r.Cookie.SetCookie() 传入key value domain path等Cookie选项设置cookie
r.Cookie.SetHttpCookie() 传入一个*http.Cookie对象作为参数设置cookie

*/

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Cookie.SetHttpCookie(&http.Cookie{
			Name:   "test",
			Value:  "1234",
			Secure: true,
			// 如果SameSite=None，则Secure必须为true，否则请求中无法携带第三方cookie;
			// 这样设置则第三方cookie必须通过HTTPS请求发送
			SameSite: http.SameSiteNoneMode,
		})
	})
}
