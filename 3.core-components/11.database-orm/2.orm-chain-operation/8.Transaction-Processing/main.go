package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*

TODO：19 什么时候开启事务？？？  当涉及多次数据修改或插入操作时，且数据必须保持一致性，要么全成功，要么全失败，这就需要使用事务
        例如：添加用户接口，首先添加用户需要一次写入操作，然后权限表设置用户的权限需要一次，岗位表又需要写入一次,这三张表的数据需要保持一致性
        就要用到事务(因为一旦中途操作失败了,会进行数据回滚,如果不开启事务则会出现修改了一个表,但是其他几个表的数据没有修改的情况,数据的一致性就被破坏了)

什么是TX事务???
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model()
		group.GET("/tx", func(r *ghttp.Request){
			
		})

		group.GET("/raw-tx", func(r *ghttp.Request){

		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
