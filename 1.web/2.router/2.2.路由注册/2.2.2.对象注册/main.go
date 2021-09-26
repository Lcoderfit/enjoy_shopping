package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
对象注册(带有BindObject支持:对象注册，对象方法注册;  BindHandler支持： 函数注册、包方法注册、对象方法注册三种方式，注意：不支持对象注册)
1.BindObject（面试的时候以举例子的方式说明使用的方式）
	1.1 对象注册
	被注册的对象方法必须为 func(*ghttp.Request)的形式
	Index是一个特殊的方法，当注册路由规则为/object时候，HTTP请求到/object/index /object都将映射到Index方法

	1.2 路由内置变量:
	{.struct} 表示路由注册时当前注册的对象名
	{.method} 表示路由注册时当前注册的方法名

	1.3 命名风格规则设置
	UriTypeDefault  /user/show-list
	UriTypeFullName /User/ShowList
	UriTypeAllLower /user/showlist
	UriTypeCamel	/user/showList
	通过SetNameToUriType设置Uri命名风格
	s.SetNameToUriType(ghttp.UriTypeDefault)

	1.4 对象方法注册(将路由绑定到指定的对象方法执行)
	BindObject("/object", b, "Show")

	1.5 指定注册的对象方法
	例如：s5.BindObject("/object", c, "Show,Order")， 指定注册c的Show和Order方法
	注意：BindObject方法虽然第三个参数是不定参，但是源码中会判断输入的不定参数是否超过1个，如果超过1个则只会取第一个，
	另外如果传入空字符会报错：invalid method name

2.BindObjectMethod
	2.1 对象方法注册
	BindObjectMethod和BindObject的区别

	BindObjectMethod 将路由绑定到指定的对象方法，当访问注册的路由时就会调用绑定到的对象方法（且BindObjectMethod的第三个参数只能指定一个方法名）

	BindObject的第三个参数可以指定多个方法名（用英文逗号隔开），且是根据指定的方法名来生成路由的， 例如（注册/object路由，然后指定"Show"方法）
	则HTTP请求到/object/show 时才会调用Show方法



3.BindObjectRest
	3.1 RESTful 对象注册 （只有与HTTP Method同名且可导出的对象方法才会被注册，HTTP Method请求会被映射到同名的对象方法）
		如果对象方法与HTTP method不同名，则不会被注册（对外不可见）

      SERVER     | DOMAIN  | ADDRESS | METHOD |  ROUTE  |          HANDLER           | MIDDLEWARE
-----------------|---------|---------|--------|---------|----------------------------|-------------
  BindObjectRest | default | :8207   | DELETE | /object | main.(*Controller1).Delete |
-----------------|---------|---------|--------|---------|----------------------------|-------------
  BindObjectRest | default | :8207   | GET    | /object | main.(*Controller1).Get    |
-----------------|---------|---------|--------|---------|----------------------------|-------------
  BindObjectRest | default | :8207   | POST   | /object | main.(*Controller1).Post   |
-----------------|---------|---------|--------|---------|----------------------------|-------------

4.Init和Shut
Init和Shut是HTTP请求流程中被Server自动调用的回调方法
Init在Server接收到请求时会自动调用，在服务接口被调用之前执行
Shut在请求结束时被Server自动调用
*/

// 对象注册
type Controller struct{}

func (c *Controller) Index(r *ghttp.Request) {
	r.Response.Writeln("index")
}

func (c *Controller) Show(r *ghttp.Request) {
	r.Response.Writeln("show")
}

// 路由内置变量
type Order struct{}

func (o *Order) List(r *ghttp.Request) {
	r.Response.Writeln("list")
}

// 命名风格规则设置
type User struct{}

func (u *User) ShowList(r *ghttp.Request) {
	r.Response.Writeln("show list")
}

// RESTful方法注册
type Controller1 struct{}

func (c *Controller1) Get(r *ghttp.Request) {
	r.Response.Writeln("Get")
}

func (c *Controller1) Post(r *ghttp.Request) {
	r.Response.Writeln("Post")
}

func (c *Controller1) Delete(r *ghttp.Request) {
	r.Response.Writeln("Delete")
}

// 该方法将无法注册(只有跟HTTP请求方法同名的方法才会被注册，否则不会被注册)
func (c *Controller1) Hello(r *ghttp.Request) {
	r.Response.Writeln("Hello")
}

// 构造方法和析构方法
type Controller2 struct{}

func (c *Controller2) Init(r *ghttp.Request) {
	r.Response.Writeln("Init")
}

func (c *Controller2) Hello(r *ghttp.Request) {
	r.Response.Writeln("Hello")
}

func (c *Controller2) Shut(r *ghttp.Request) {
	r.Response.Writeln("Shut")
}

func main() {
	s := g.Server()

	// 1.1.对象注册
	c := new(Controller)
	s.BindObject("/object", c)
	// 1.2.路由内置变量, 当HTTP请求/order-list时，则会调用Order结构体的List方法
	o := new(Order)
	s.BindObject("/{.struct}-{.method}", o)

	s.SetPort(8199)
	s.Start()

	// 1.3.命名风格规则配置
	u := new(User)
	s1 := g.Server("UriTypeDefault")
	s2 := g.Server("UriTypeFullName")
	s3 := g.Server("UriTypeAllLower")
	s4 := g.Server("UriTypeCamel")

	// 设置Uri命名风格
	// 默认通过英文连接符连接： show-list
	s1.SetNameToUriType(ghttp.UriTypeDefault)
	// 首字母大写（大驼峰）: ShowList
	s2.SetNameToUriType(ghttp.UriTypeFullName)
	// 全小写: showlist
	s3.SetNameToUriType(ghttp.UriTypeAllLower)
	// 小驼峰： showList
	s4.SetNameToUriType(ghttp.UriTypeCamel)

	// /user/show-list
	s1.BindObject("/{.struct}/{.method}", u)
	// /User/ShowList
	s2.BindObject("/{.struct}/{.method}", u)
	// /user/showlist
	s3.BindObject("/{.struct}/{.method}", u)
	// /user/showList
	s4.BindObject("/{.struct}/{.method}", u)

	s1.SetPort(8201)
	s2.SetPort(8202)
	s3.SetPort(8203)
	s4.SetPort(8204)

	s1.Start()
	s2.Start()
	s3.Start()
	s4.Start()

	// 1.4. 对象方法注册(如果一个对象有多个方法，但是我只想注册其中的一部分，则可以使用BindObject的第三个参数)
	s5 := g.Server("method register")
	// 只注册 Controller的Show方法，如果希望注册多个，则多个方法通过英文的逗号隔开，例如：s5.BindObject("/object", c, "Show,Order")
	// 注意：BindObject方法虽然第三个参数是不定参，但是源码中会判断输入的不定参数是否超过1个，如果超过1个则只会取第一个，另外如果传入空字符也会注册所有方法
	// HTTP请求 /object/show 时会调用Show方法
	s5.BindObject("/object", c, "Show")
	s5.SetPort(8205)
	s5.Start()

	// 2.1 绑定路由方法(绑定制定的路由到执行的对象方法执行)
	s6 := g.Server("BindObjectMethod")
	// HTTP请求/show 时会调用 Show方法
	s6.BindObjectMethod("/show", c, "Show")
	s6.SetPort(8206)
	s6.Start()

	// 3.1 RESTful对象注册(通常被用于API服务)
	c1 := new(Controller1)
	s7 := g.Server("BindObjectRest")
	// RESTFul方法注册，HTTP的Get请求将被映射到对象的Get方法（是什么请求就被映射到对象的同名方法）
	// 如果控制器并未定义对应的HTTP Method方法，那么该HTTP Method请求将收到404响应
	s7.BindObjectRest("/object", c1)
	s7.SetPort(8207)
	s7.Start()

	// 3.2 构造方法和析构方法
	c2 := new(Controller2)
	s8 := g.Server("Init-and-Shut")
	// Init
	// Hello
	// Shut
	s8.BindObject("/object", c2)
	s8.SetPort(8208)
	s8.Start()

	g.Wait()
}
