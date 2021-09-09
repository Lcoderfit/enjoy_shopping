package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.设置属性默认值，通过d标签或default标签设置
	如果v标签设置了required, 且输入参数为空是时，则Name属性会取d/default标签设置的默认值，即required校验规则不会报错
*/

type Test struct {
	Name string `v:"required#请输入姓名" d:"john"`
}

type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		var t Test
		if err := r.Parse(&t); err != nil {
			r.Response.WriteJsonExit(Response{
				Code:  1,
				Error: err.Error(),
			})
		}
		r.Response.WriteJsonExit(Response{
			Data: t,
		})
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
