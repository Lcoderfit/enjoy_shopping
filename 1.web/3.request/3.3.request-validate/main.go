package main

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

/*
1.r.Response.WriteJsonExit()
	直接传入一个结构体，不用传入指针

2.请求结构体不用设置json标签（仅用于数据校验和接受请求参数）;响应结构体需设置json标签

3.校验错误处理
	3.1 如果不设置 #请输入账户|账户长度为:min到:max位.... 则默认返回的错误信息是英文的，例如：
		curl "http://127.0.0.1:8199/register?password1=123456&password2=123456"
		{"code":1,"error":"The username field is required","data":null}

	3.2 v标签用法：`v:"required|length:4,30#请输入账户|账户长度为:min到:max位"`, 注意，如果要使用length的最大最小值，
		则需要通过:min和:max获取，如果仅适用min的话只是"min"字段串，而":min"会被解析为4

	3.2 一般情况下，字段校验会返回所有校验失败的错误信息：
		curl "http://127.0.0.1:8199/register"
		{"code":1,"error":"请输入账号; 账号长度为4到30位; 请输入密码; 密码长度不够; 请确认密码; 密码长度不够; 两次密码不一致","data":null}

		一般只需要返回第一条错误信息即可，通过类型断言或者gerror.Current
		3.2.1 类型断言
			if v, ok := err.(gvalid.Error); ok {
				r.Response.WriteJsonExit(RegisterRes{
					Code:  1,
					Error: v.FirstString(),
				})
			}
		curl "http://127.0.0.1:8199/register?password1=123456&password2=123456"
		{"code":1,"error":"请输入账户","data":null}

		3.2.2 gerror
		gerror.Current(err).Error() 获取第一条错误信息
		curl "http://127.0.0.1:8199/register?name=lco"
		{"code":1,"error":"账户长度为4到30位","data":null}

	3.4 gf从1.6版本开始，HTTP请求数据的校验不再受结构体默认值影响，也就是说，如果结构体属性没有匹配到参数值，则结构体属性会为其类型对应的零值，
		但是该零值不会对校验产生影响，校验还是针对参数？？？？？？？？？
*/

// 注册请求数据结构
type RegisterReq struct {
	// p标签都是可选的，默认使用不区分大小写和忽略-_空格的方式匹配
	Name  string `p:"username" v:"required|length:4,30#请输入账户|账户长度为:min到:max位"`
	Pass  string `p:"password1" v:"required|length:6,30#请输入密码|密码长度不够"`
	Pass2 string `p:"password2" v:"required|length:6,30|same:password1#请输入密码|密码长度不够|两次密码不一致"`
}

// 注册返回数据结构
type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/register", func(r *ghttp.Request) {
			var req *RegisterReq
			if err := r.Parse(&req); err != nil {
				//r.Response.WriteJsonExit(RegisterRes{
				//	Code:  1,
				//	Error: err.Error(),
				//})

				// 使用类型断言，如果判读是校验错误，则只返回第一条校验错误
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
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("gerror")
	s1.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/register", func(r *ghttp.Request) {
			var req RegisterReq
			if err := r.Parse(&req); err != nil {
				r.Response.WriteJsonExit(RegisterRes{
					Code: 1,
					// gerror.Current获取第一条报错信息
					Error: gerror.Current(err).Error(),
				})
			}
			r.Response.WriteJsonExit(RegisterRes{
				Data: req,
			})
		})
	})
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
