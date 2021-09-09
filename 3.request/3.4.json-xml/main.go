package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

/*
1.r.Response.Writef("%s", "coder")

2.r.Get接收各种格式的参数
	2.1 url参数
	curl "http://localhost:8200?name=john&pass=123"
	name:john, pass:123

	2.2 form表单提交
	curl -d "name=john&pass=123" "http://localhost:8200"
	name:john, pass:123

	2.3 JSON数据格式(linux系统跟windows系统使用curl发送json数据格式不太一样)
	linux:
		curl -d '{"name":"john","pass":"123"}' "http://127.0.0.1:8200/"
		name:john, pass:123
	windows:
		curl -d "{\"name\":\"john\",\"pass\":\"123\"}" "http://127.0.0.1:8200/"
		name:john, pass:123

	2.4 xml格式数据
	linux:
		curl -d '<?xml version="1.0" encoding="UTF-8"?><doc><name>john</name><pass>123</pass></doc>' "http://127.0.0.1:8200/"
		name:john, pass:123

		curl -d '<doc><name>john</name><pass>123</pass></doc>' "http://127.0.0.1:8199/"
		name:john, pass:123


3.gf工具链使用：周六/周日
	https://hub.fastgit.org/gogf/gf-cli

4.对象转换及校验
	4.1 same校验规则： same:password1 表示Pass2属性值必须与password1参数相同
		4.1.1 情况1，如果same指定参数与Pass1属性的p标签所指定参数同名，则可以认为Pass2属性的值必须与Pass1的属性值相等
		参数无论传入password1或pass1均可，只要转换后Pass1属性值与Pass2属性值相等即校验通过
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"same:password1"`

		4.1.2 情况2，same标签与Pass1属性同名，但是与Pass1属性设置的p标签不同名，则如果参数传入password1是没有效果的，
		始终会报两字段不相等的校验错误(可以认为此时判读的是Pass2属性值是否与pass1参数(注意是pass1参数，不是Pass1属性)相等)
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"same:pass1"`

	4.2 linux(r.Parse方法自带了对客户端json/xml数据格式的解析功能)
		curl -d '{"username":"johngcn","password1":"123456","password2":"123456"}' "http://127.0.0.1:8201/register"
		{"code":0,"error":"","data":{"Name":"johngcn","Pass1":"123456","Pass2":"123456"}}

		curl -d '<?xml version="1.0" encoding="UTF-8"?><doc><username>johngcn</username><password1>123456</password1><password2>123456</password2></doc>' "http://127.0.0.1:8201/register"
		{"code":0,"error":"","data":{"Name":"johngcn","Pass1":"123456","Pass2":"123456"}}

*/

type RegisterReq struct {
	Name  string `p:"username" v:"required|length:6,30#请输入账号|账号长度:min到:max位"`
	Pass1 string `p:"password1" v:"required|length:6,30#请输入密码|密码长度不够"`
	Pass2 string `p:"password2" v:"required|length:6,30|same:password1#请确认密码|密码长度不够|两次密码不一致"`
}

type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writef("name:%v, pass:%v", r.Get("name"), r.Get("pass"))
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("gvalid")
	s1.BindHandler("/register", func(r *ghttp.Request) {
		var req *RegisterReq
		if err := r.Parse(&req); err != nil {
			if v, ok := err.(gvalid.Error); ok {
				r.Response.WriteJsonExit(RegisterRes{
					Code:  1,
					Error: v.FirstString(),
				})
			}
			r.Response.WriteJsonExit(RegisterRes{
				Code:  1,
				Error: err.Error(),
			})
		}
		r.Response.WriteJsonExit(RegisterRes{
			Data: req,
		})
	})
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
