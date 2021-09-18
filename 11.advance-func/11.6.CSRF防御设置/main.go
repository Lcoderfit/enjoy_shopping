package main

import (
	"github.com/gogf/csrf"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
	"time"
)

func main() {
	s := g.Server()
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(csrf.NewWithCfg(csrf.Config{
			TokenLength:     32,
			TokenRequestKey: "X-Token",
			ExpireTime:      time.Hour * 24,
			Cookie: &http.Cookie{
				Name: "_csrf", // token name in cookie
			},
		}))
		group.ALL("/csrf", func(r *ghttp.Request) {
			r.Response.Writeln(r.Method + ":" + r.RequestURI)
		})
	})
}

/*
s := g.Server()
s.Group("/api.v2", func(group *ghttp.RouterGroup) {
	group.Middleware(csrf.NewWithCfg(csrf.Config{
		Cookie: &http.Cookie{
			Name: "_csrf",// token name in cookie
		},
		ExpireTime:      time.Hour * 24,
		TokenLength:     32,
		TokenRequestKey: "X-Token",// use this key to read token in request param
	}))
	group.ALL("/csrf", func(r *ghttp.Request) {
		r.Response.Writeln(r.Method + ": " + r.RequestURI)
	})
})
*/
