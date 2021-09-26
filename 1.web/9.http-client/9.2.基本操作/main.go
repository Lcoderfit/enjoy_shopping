package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	if r, err := g.Client().Get("https://goframe.org"); err != nil {
		panic(err)
	} else {
		defer r.Close()
		fmt.Println(r.ReadAllString())
	}
}
