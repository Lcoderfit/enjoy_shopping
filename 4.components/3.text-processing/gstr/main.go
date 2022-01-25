package main

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	s := "垃圾同*顺"
	rs := gstr.SubStrRune(s, 2)
	fmt.Println(rs)
}
