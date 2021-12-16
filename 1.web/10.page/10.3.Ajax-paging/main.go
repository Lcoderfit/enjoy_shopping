package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gview"
)

/*
1.Ajax分页
	page.AjaxActionName = "DoAjax"

	<script src="https://cdn.bootcss.com/jquery/2.0.3/jquery.min.js"></script>
	<script>
	function DoAjax(url) {
		 $.get(url, function(data,status) {
			 $("body").html(data);
		 });
	}
	</script>
*/

func main() {
	s := g.Server()
	s.SetPort(8200)
	s.BindHandler("/page/ajax", func(r *ghttp.Request) {
		page := r.GetPage(100, 10)
		page.AjaxActionName = "DoAjax"
		buffer, _ := gview.ParseContent(r.Context(), `
        <html>
            <head>
                <style>
                    a,span {padding:8px; font-size:16px;}
                    div{margin:5px 5px 20px 5px}
                </style>
                <script src="https://cdn.bootcss.com/jquery/2.0.3/jquery.min.js"></script>
                <script>
                function DoAjax(url) {
                     $.get(url, function(data,status) {
                         $("body").html(data);
                     });
                }
                </script>
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
	s.Start()

	g.Wait()
}
