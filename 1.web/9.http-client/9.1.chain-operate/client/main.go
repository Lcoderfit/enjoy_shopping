package main

import (
	"github.com/gogf/gf/frame/g"
	"time"
)

func main() {
	res := g.Client().Timeout(4 * time.Second).GetContent("http://localhost:8200")
	g.Log().Line(true).Println(res)

	res = g.Client().Timeout(3*time.Second).PostContent("http://localhost:8200", g.Map{
		"name": "john",
		"age":  10,
	})
	g.Log().Line(true).Println(res)
}
