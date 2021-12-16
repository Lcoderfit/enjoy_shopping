package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gview"
)

/*
1.动态分页
	1.1 r.GetPage(totalSize, pageSize): 第一个为总数量，第二个为每页的数量
	1.2 page.GetContent() 有四种模式
		模式1：只显示当前页码
		模式2：显示当前页码，并且有下拉框选择页码
		模式3：显示所有页码，并显示当前页码（5/10）和数据总条数(98条)
		模式4：显示所有页码

	1.3 前端通过Get传递QueryString参数，默认分页参数名为page,前端只需要传入第几页，后端控制数据总量和总共分多少页
*/

func main() {
	s := g.Server()
	s.BindHandler("/page/demo", func(r *ghttp.Request) {
		page := r.GetPage(98, 10)
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
			"page4": page.GetContent(4),
		})
		r.Response.Write(buffer)
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}
