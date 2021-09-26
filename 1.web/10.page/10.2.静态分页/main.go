package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gview"
)

/*
1.静态路由
	路由定义中需要给定一个page参数（字段匹配，命名匹配，模糊匹配均可）

*/

func main() {
	s := g.Server()
	s.BindHandler("/page/static/*page", func(r *ghttp.Request) {
		page := r.GetPage(100, 10)
		buffer, _ := gview.ParseContent(r.Context(), `
        <html>
            <head>
                <style>
                    a,span {padding:8px; font-size:16px;}
                    div{margin:5px 5px 20px 5px}
                </style>
            </head>
            <body>
                <div>{{.page1}}</div>
                <div>{{.page2}}</div>
                <div>{{.page3}}</div>
                <div>{{.page4}}</div>
            </body>
        </html>
		`, g.Map{
			"page1": page.GetContent(1),
			"page2": page.GetContent(2),
			"page3": page.GetContent(3),
			"page":  page.GetContent(4),
		})
		r.Response.Writeln(buffer)
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
