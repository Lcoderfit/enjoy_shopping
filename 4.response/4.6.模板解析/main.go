package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

/*
1.ParseTplxxx 和 WriteTplxxx
	func (r *Response) WriteTpl(tpl string, params ...gview.Params) error
	func (r *Response) WriteTplContent(content string, params ...gview.Params) error
	func (r *Response) WriteTplDefault(params ...gview.Params) error

	func (r *Response) ParseTpl(tpl string, params ...gview.Params) (string, error)
	func (r *Response) ParseTplContent(content string, params ...gview.Params) (string, error)
	func (r *Response) ParseTplDefault(params ...gview.Params) (string, error)

	1.1 Tpl和TplContent和TplDefault的区别
		Tpl方法接收两个参数，第一个参数是模板文件的路径，第二个参数是需要解析到模板的数据（是一个map[string]interface{}类型）
		TplContent 第一个参数是模板字符串，例如`{.Name}`，第二个参数是map[string]interface{}类型，作用是将数据根据模板解析成对应的值
		TplDefault 接收一个参数，这个参数是用来解析到默认模板文件(即index.html)的变量，也是个map[string]interface{}
	注意：params参数可以传入多个数据，均会被解析成对应的值


	1.2 ParseTplxxx 和 WriteTplxxx 的区别
		ParseTplxxx根据data的值将模板语法解析成对应的string（返回一个string和一个error类型）
		WriteTplxxx根据data的值将模板语法解析成对应的数据然后写入到响应体中(返回一个error类型)

2.解析规则
	2.1 data必须为一个map[string]interface{}类型，可以根据gconv.Map()将结构体类型转换为map[string]interface类型
		params参数（不定参）可以传递多个data，均会被解析成对应的值

	2.2 如果是结构体根据gconv.Map转换为map[string]interface{}类型的，则结构体字段中只有可导出的字段能够被转换成功，
		因为不可导出的字段在gconv包中是无法访问的，所以转换成的map[string]interface{}中是不包含不可导出字段对应的键值对的

	2.3 如果data是一个map，且data的key是大写的，则模板中名称也需要大写: {{.Name}}; .表示的是当前的变量，.Name表示的是当前变量中的Name字段

	2.4 嵌套结构体的内层结构体中的字段通过gconv.Map的转换后字段都会变为最外层的key,
			注意：如果外层结构体与内层结构体的字段名冲突，则转换成map后只取最外层的结构体字段值作为map的value
		例如下面这个结构体，将转换为 map[string]interface{"Link": "link", conf: "conf"};
			Target{
				Link: "link",
				Config{
					conf: "conf",
					Server{
						// 内层的Link将被最外层的Link覆盖
						Link: "link",
					 },
				},
			}
		但是自定义变量无法覆盖内置变量：（例如Session，Request，Config等）

	2.5 {{dump .}} 表示将当前所有的变量（包括传入的参数和config.toml文件中的内置变量）, 序列化(格式化之后)之后的样子

	2.6 内置变量
	{{.Config.配置项}} // config.toml配置项
	{{.Cookie.键名}}	// 当前请求的Cookie对象参数值
	{{.Session.键名}} //当前请求的Session对象参数值
	{{.Query.键名}}	//url中的请求参数
	{{.Form.键名}}	//表单请求参数
	{{.Request.键名}} //当前请求参数中，不区分提交方式（包括了Query和Form）

*/

type Server struct {
	Link    string
	Fuck    string
	Session string
}

type Config struct {
	Server
	Fuck string
}

type Target struct {
	name string
	Age  int
	Config
}

func main() {
	t := Target{
		name: "coder",
		Age:  10,
		Config: Config{
			Fuck: "lcoderfit",
			Server: Server{
				Link:    "link",
				Fuck:    "link-fuck",
				Session: "session",
			},
		},
	}
	target := gconv.Map(t)

	f := func(t Target) {
		g.Log().Println(t.name)
	}
	f(t)
	g.Log().Line(true).Info(t)

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			st := `{{.Session}}`
			res, _ := r.Response.ParseTplContent(st, target)
			g.Log().Println(res)
			_ = r.Response.WriteTplContent(`{{dump .}}, {{.Config.server.Address}}, {{.Config.Server.Link}}, {{.Config.Fuck}}`, target)
			_ = r.Response.WriteTplDefault(target)
		})
	})
	s.SetPort(8200)
	s.Start()

	s1 := g.Server("s1")
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Cookie.Set("theme", "default")
		r.Session.Set("name", "john")
		content := `Cookie:{{.Cookie.theme}}, Session:{{.Session.name}}`
		// 仅使用内置变量时，第二个params参数就传入nil即可
		_ = r.Response.WriteTplContent(content, nil)
	})
	s1.SetPort(8201)
	s1.Start()

	g.Wait()
}
