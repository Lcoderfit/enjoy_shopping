package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
1.字段与属性转换时的默认映射规则
	1.1 映射的目标结构体的属性必须是可导出的，否则无法映射
	1.2 请求参数与目标结构体属性进行匹配时，不区分大小写,且会忽略 连接符,下划线，空格  （参数和属性双方均忽略）
	1.3 如果匹配成功，则将参数值赋给属性，否则忽略该键值，属性值默认为对应类型的零值

2.自定义映射关系
	通过p标签(或者param/params标签)设置特定的参数名映射到属性，例如： Pass1 string `p:"password1"` 即是设置参数password1映射到 Pass1属性
	注意：即使设置了p标签，默认的参数与属性映射关系仍然成立，即如果参数名为pass1, 仍可映射到Pass1属性

3. r.Parse()
	3.1 r.Parse接收一个结构体指针或结构体双重指针作为参数（该结构体为转换后的目标结构体），如果传入一个空结构体指针（即未初始化），如果符合转换规则
		则在Parse内部会为该空指针初始化并进行参数映射,
	3.2 如果目标结构体中可导出字段包含v标签，则调用Parse时会对字段进行校验

{"code":0,"error":"","data":{"Name":"john","Pass1":"123","Pass2":"456"}}

4.Only one usage of each socket address (protocol/network address/port) is normally permitted.
	每个套接字地址（协议/网络地址/端口）通常只允许使用一次
	一般是两个server没有传递参数进行区分，由于单例模式导致对同一个server调用了两次Start方法，所以就会产生两个两个进程对同一端口进行监听
	对同一个Server设置多端口是可以的，s.SetPort(8200, 8201)
	注意：同一个Server，两次调用SetPort设置不同端口也可以，但是后面的设置会覆盖前面的设置，例如：

	s.SetPort(8200)
	s.SetPort(8201)

	则最终只有最后一个端口生效(8021)

5. r.Response.WriteJsonExit()
	会以json格式返回数据，并且相当于return，即跳出当前服务方法（不再执行当前服务方法内的后续逻辑），但是当前服务方法后如果还存在其他请求处理流程，
	则仍会执行，例如调用WriteJsonExit()后，仍会返回MiddlewareTest中间件执行	g.Log().Line(true).Println("this is a test")

6.json tags
	curl "http://localhost:8201/register?name=john&password1=123&password2=456"
	返回：{"code":0,"error":"","data":{"Name":"john","Pass1":"123","Pass2":"456"}}
	RegisterRes中设置了json标签，所以返回的数据会转换为JSON格式（返回数据的字段名为json标签设置的字段名）
	但是由于RegisterReq没有设置json标签，所以返回的数据中字段名为属性名
	可以通过添加json标签设置：
	// 参数映射时，参数可以传入 password1=1或pass1=1均可映射成功，然后通过WriteJsonExit()会将属性名转换为json标签设置的名字
	type RegisterReq struct {
		Name  string
		Pass1 string `p:"password1" json:"p1"`
		Pass2 string `p:"password2" json:"p2"`
	}

7. curl post

curl -d "name=john&password1=123&password2=456" -X POST "http://127.0.0.1:8199/register"
{"code":0,"error":"","data":{"Name":"john","Pass1":"123","Pass2":"456"}}
*/

type User struct {
	Id    int
	Name  string
	Pass1 string `p:"password1"`
	Pass2 string `p:"password2"`
}

type RegisterReq struct {
	Name  string
	Pass1 string `p:"password1"`
	Pass2 string `p:"password2"`
}

type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func MiddlewareTest(r *ghttp.Request) {
	r.Middleware.Next()
	g.Log().Line(true).Println("this is a test")
}

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		user := new(User)
		if err := r.Parse(&user); err != nil {
			r.Response.Writeln(err)
		} else {
			r.Response.Writeln(user)
		}
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("parse")
	s1.Use(MiddlewareTest)
	s1.BindHandler("/register", func(r *ghttp.Request) {
		var req *RegisterReq
		if err := r.Parse(&req); err != nil {
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
