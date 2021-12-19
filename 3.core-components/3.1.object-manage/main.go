package main

/*对象管理
一、耦合与内聚
	1.耦合：软件模块与模块之间的关联程度
	2.内聚：与耦合相对,指一个模块内部各个成分之间的相关程度的度量
二、goframe内置常用数据类型
	2.1 Var/Ctx -> gvar.Var/context.Context
		gvar.Var是内置的数据类型，可以调用各种类型的转换方法
	2.2 四大类: MapAB/ListAB/SliceA/ArrayA
		2.2.1 MapAB A表示的是map[A]B类型, g.Map表示的map[string]interface{}类型
		2.2.2 ListAB表示的是[]MapAB类型，g.List表示的是[]map[string]interface{}
		2.2.3 Array/Slice/ArrayAny/SliceAny/ 表示的是[]interface{}类型
		2.2.4 SliceStr/ArrayStr表示的是[]string类型
		2.2.5 SliceInt/ArrayInt []int类型

		type (
			Var = gvar.Var        // Var is a universal variable interface, like generics.
			Ctx = context.Context // Ctx is alias of frequently-used context.Context.
		)

		type (
			Map        = map[string]interface{}      // Map is alias of frequently-used map type map[string]interface{}.
			MapAnyAny  = map[interface{}]interface{} // MapAnyAny is alias of frequently-used map type map[interface{}]interface{}.
			MapAnyStr  = map[interface{}]string      // MapAnyStr is alias of frequently-used map type map[interface{}]string.
			MapAnyInt  = map[interface{}]int         // MapAnyInt is alias of frequently-used map type map[interface{}]int.
			MapStrAny  = map[string]interface{}      // MapStrAny is alias of frequently-used map type map[string]interface{}.
			MapStrStr  = map[string]string           // MapStrStr is alias of frequently-used map type map[string]string.
			MapStrInt  = map[string]int              // MapStrInt is alias of frequently-used map type map[string]int.
			MapIntAny  = map[int]interface{}         // MapIntAny is alias of frequently-used map type map[int]interface{}.
			MapIntStr  = map[int]string              // MapIntStr is alias of frequently-used map type map[int]string.
			MapIntInt  = map[int]int                 // MapIntInt is alias of frequently-used map type map[int]int.
			MapAnyBool = map[interface{}]bool        // MapAnyBool is alias of frequently-used map type map[interface{}]bool.
			MapStrBool = map[string]bool             // MapStrBool is alias of frequently-used map type map[string]bool.
			MapIntBool = map[int]bool                // MapIntBool is alias of frequently-used map type map[int]bool.
		)

		type (
			List        = []Map        // List is alias of frequently-used slice type []Map.
			ListAnyAny  = []MapAnyAny  // ListAnyAny is alias of frequently-used slice type []MapAnyAny.
			ListAnyStr  = []MapAnyStr  // ListAnyStr is alias of frequently-used slice type []MapAnyStr.
			ListAnyInt  = []MapAnyInt  // ListAnyInt is alias of frequently-used slice type []MapAnyInt.
			ListStrAny  = []MapStrAny  // ListStrAny is alias of frequently-used slice type []MapStrAny.
			ListStrStr  = []MapStrStr  // ListStrStr is alias of frequently-used slice type []MapStrStr.
			ListStrInt  = []MapStrInt  // ListStrInt is alias of frequently-used slice type []MapStrInt.
			ListIntAny  = []MapIntAny  // ListIntAny is alias of frequently-used slice type []MapIntAny.
			ListIntStr  = []MapIntStr  // ListIntStr is alias of frequently-used slice type []MapIntStr.
			ListIntInt  = []MapIntInt  // ListIntInt is alias of frequently-used slice type []MapIntInt.
			ListAnyBool = []MapAnyBool // ListAnyBool is alias of frequently-used slice type []MapAnyBool.
			ListStrBool = []MapStrBool // ListStrBool is alias of frequently-used slice type []MapStrBool.
			ListIntBool = []MapIntBool // ListIntBool is alias of frequently-used slice type []MapIntBool.
		)

		type (
			Slice    = []interface{} // Slice is alias of frequently-used slice type []interface{}.
			SliceAny = []interface{} // SliceAny is alias of frequently-used slice type []interface{}.
			SliceStr = []string      // SliceStr is alias of frequently-used slice type []string.
			SliceInt = []int         // SliceInt is alias of frequently-used slice type []int.
		)

		type (
			Array    = []interface{} // Array is alias of frequently-used slice type []interface{}.
			ArrayAny = []interface{} // ArrayAny is alias of frequently-used slice type []interface{}.
			ArrayStr = []string      // ArrayStr is alias of frequently-used slice type []string.
			ArrayInt = []int         // ArrayInt is alias of frequently-used slice type []int.
		)

三、常用对象
	3.1 常用对象通过单例模式进行管理，通过设置不同的单例名称获取不同的单例对象(例如g.Server("a"))
	3.2各种常用对象
		1.客户端对象: g.Client().Do
		2.校验对象: g.Validator().xxx
		3.配置管理对象: g.Cfg(),默认读取config.toml/yaml/yml/json/ini/xml并缓存，配置文件被修改时会自动刷新缓存
			配置管理对象会自动使用单例名称对文件进行检索: g.Cfg("redis") 则会默认检索 redis.toml/yaml/yml/json/init/xml
			当文件不存在时，才会检索config.toml
		4.日志管理对象：g.Log() 会读取默认配置文件中的logger配置项，且只会初始化一次日志对象
		5.g.View()
		6.读取默认配置文件中的server配置项: g.Server()
		7.g.TcpServer()
		8.g.UdpServer()
		9.ORM对象 g.DB()/g.Model() 默认读取配置文件中的database配置项,
		10.redis客户端对象: g.Redis() 默认读取配置文件中的redis配置项
		11.资源管理对象: g.Res()
		12.国际化管理对象: g.I18n()

*/

func main() {

}
