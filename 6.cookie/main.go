package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
type Cookie
    func GetCookie(r *Request) *Cookie
    func (c *Cookie) Contains(key string) bool
    func (c *Cookie) Flush()
    func (c *Cookie) Get(key string, def ...string) string
    func (c *Cookie) GetSessionId() string
    func (c *Cookie) Map() map[string]string
    func (c *Cookie) Remove(key string)
    func (c *Cookie) RemoveCookie(key, domain, path string)
    func (c *Cookie) Set(key, value string)
    func (c *Cookie) SetCookie(key, value, domain, path string, maxAge time.Duration, httpOnly ...bool)
    func (c *Cookie) SetHttpCookie(httpCookie *http.Cookie)
    func (c *Cookie) SetSessionId(id string)
*/

func main() {
	s := g.Server()
	s.BindHandler("/cookie", func(r *ghttp.Request) {
		datetime := r.Cookie.Get("datetime")
		r.Cookie.Set("datetime", gtime.Datetime())
		r.Response.Writeln(datetime)

		ts := r.Cookie.Get("ts")
		r.Cookie.SetCookie("ts", "10", "", "/", time.Second*60)
		r.Response.Writeln(ts)

	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
