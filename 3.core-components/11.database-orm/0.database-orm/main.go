package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"path/filepath"
	"runtime"
)

/*
数据库orm
https://goframe.org/pages/viewpage.action?pageId=1114686

orm获取数据库操作对象
1.获取当前系统名: runtime.GOOS
2.设置config文件路径
	2.1.如果是gf run main.go 则当前项目路径为main.go所在的路径
	2.2.如果是go run main.go 则当前项目路径为go.mod所在的路径
	上面这两条无论在linux环境还是在windows环境均适用

	g.Cfg().SetPath("../config")

3.获取数据库操作对象
	3.1 g.DB() 获取默认配置的数据库对象(单例模式,获取[database] -> [[database.default]])下的配置
		g.DB("name") 创建新的单例
	3.2 gdb.New() 原生New方法创建对象, 不同的gdb.New()创建的是两个不同的数据库对象,非单例
		gdb.New("name")
	3.3 gdb.Instance() 原生的单例管理方法创建数据库对象
		gdb.Instance("name")

	g.DB() 返回 DB实例
	gdb.New() 返回 (DB, error)
	gdb.Instance() 返回 (DB, error)
*/

type result struct {
	AdId       int `orm:"ad_id"`
	CustomerId int `orm:"customer_id"`
}

func main() {
	g.Log().Line(true).Info(filepath.Abs("."))
	osName := runtime.GOOS
	// 1.如果是gf run main.go 则当前项目路径为main.go所在的路径
	// 2.如果是go run main.go 则当前项目路径为go.mod所在的路径
	// 上面这两条无论在linux环境还是在windows环境均适用
	if osName == "windows" {
		g.Cfg().SetPath("./3.core-components/11.database-orm/config")
	} else if osName == "linux" {
		g.Log().Line(true).Info(filepath.Abs("../config"))
		g.Cfg().SetPath("../config")
	} else {
		g.Log().Line(true).Warning("os not supported")
		return
	}

	db := g.DB()
	// 这里不使用Fields过滤字段也是可以的,因为Scan会自动匹配res中结构体的属性
	dbModel := db.Model("ads").Fields("ad_id, customer_id").Where("ad_id<=?", 10)

	var res []*result
	err := dbModel.Scan(&res)
	if err != nil {
		g.Log().Line(true).Error(err)
		return
	}
	for _, v := range res {
		fmt.Println(*v)
	}
}
