package main

import (
	"fmt"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	s := "垃圾同*顺"
	// SubStrRune返回unicode字符串str从start开始，长度为length的新字符串。 参数length是可选的，它默认使用str的长度。
	rs := gstr.SubStrRune(s, 2)
	fmt.Println(rs)

	var (
		haystack = `goframe is very, very easy to use`
		needle   = `very`
		// PosI返回needle在haystack中第一次出现的位置，不区分大小写。 如果没有找到，则返回-1
		posI     = gstr.PosI(haystack, needle)
		// PosR返回needle在haystack中最后一次出现的位置，区分大小写。 如果没有找到，则返回-1
		posR     = gstr.PosR(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	//
	fileName := "垃圾.jpeg"
	res := gstr.SubStrRune(fileName, gstr.PosRRune(fileName, ".")+1, len(fileName))
	fmt.Println(res)

	// 正则
	sl, _ := gregex.MatchString(`^([0-9]+)(?i:([a-z]*))$`, "50mb")
	fmt.Println(sl)

	//
	path := "/a/b/c/20220124/"
	fmt.Println(getUrl(path, "ajsdkfajsdkf"))

	var a A
	c := new(C)
	a = c
	a.B()
}

func getUrl(path, fileName string) string {
	url := gstr.SubStr(path, gstr.Pos(path, "/c/")+1) + fileName
	return url
}

type A interface {
	B()
}

type C struct {

}

func (a *C) B() {
	fmt.Println("c.B()")
}