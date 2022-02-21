package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*ORM链式操作-悲观锁 & 乐观锁

一、悲观锁(Pessimistic Lock) 美 [ˌpesɪˈmɪstɪk]
	1.1 定义：就是每次取数据都悲观的认为别人会修改，所以每次拿数据时都会上锁，其他人要取数据就会被阻塞直到拿到锁，
		例如：行锁、表锁、读锁、写锁 =》 均是在操作之前先上锁

	1.2 使用 (N：只能在事务中使用)
		1.2.1 func (m *Model) LockShared() *Model  对应SQL: where xxx lock in share mode
			共享锁：不会阻塞其他事务读取被锁定的数据，但是会阻塞其他事务更新锁定的数据直到事务提交

			如果一个事务读取，然后上锁，另一个事务也读取，但是无法修改数据，直到前一个事务提交，之后才能修改，这样数据只被修改了一次，是正确的
			TODO: 如果一个事务读取，然后上锁，另外还有两个事务读取了数值x，一个试图x+3, 一个试图x+2, 是否会出现问题????
		1.2.2 func (m *Model) LockUpdate() *Model  对应SQL: where xx For Update
			For Update锁：其他事务对锁定数据的读取或写操作均会被阻塞
			如果两个事务均要读取x，事务A试图修改为x+2, 事务B此时是没办法读的，所以事务A将x更新为x+2, 事务B此时读取到的是x+2,然后再更新,正确



二、乐观锁 Optimistic Lock  美 [ˌɑːptɪˈmɪstɪk]
	2.1 定义： 每次读取数据，都乐观的认为数据不会被修改，所以不会上锁，但是更新的时候会判断在此期间数据有没有被修改，如果没有则更新，有则更新失败
		可以通过版本号来实现：为数据增加一个版本标识，取数的时候将该版本一起读出来，之后更新时将该版本号+1，将新的版本与此时数据表中数据当前版本进行对比，
		如果新版本>数据表当前版本，则更新，否则认为是过期数据

三、优缺点
	1.乐观锁：适用于写较少的情况，这样可以省去很多锁开销
	2.悲观锁：适用于写多的场景，	因为如果经常发生冲突，上层应用会不断重试，如果用乐观锁反而降低性能
*/

func main() {
	s := g.Server()
	g.Cfg().SetPath("../../config")
	s.Group("/", func(group *ghttp.RouterGroup) {
		m := g.DB().Model("user")
		group.GET("/lock", func(r *ghttp.Request) {
			// LockShared()
			res, err := m.Where("uid>?", 0).LockShared().All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("lock-shared: ", res)
			}

			// LockUpdate()
			res, err = m.WhereLike("name", "%l%").LockUpdate().All()
			if err != nil {
				g.Log().Line(true).Error(err.Error())
				r.Response.Writeln(err.Error())
			} else {
				r.Response.Writeln("lock-update: ", res)
			}
		})
	})
	s.SetPort(8200)
	s.Start()
	g.Wait()
}
