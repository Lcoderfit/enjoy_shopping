package main

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

/*
gmd5 md5加密
https://pkg.go.dev/github.com/gogf/gf/crypto/gmd5

func Encrypt(data interface{}) (encrypt string, err error)
func EncryptBytes(data []byte) (encrypt string, err error)
func EncryptFile(path string) (encrypt string, err error)
func EncryptString(data string) (encrypt string, err error)
func MustEncrypt(data interface{}) string
func MustEncryptBytes(data []byte) string
func MustEncryptFile(path string) string
func MustEncryptString(data string) string

1.Must开头的均返回string，Encrypt开头的均返回(string, error)
2.  Encrypt 对任意类型的数据进行md5加密
	EncryptString 对string类型的数据进行加密
	EncryptBytes 对[]byte{}类型的数据进行加密
	EncryptFile 对文件内容进行加密,传入的是一个文件路径
*/

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		encryptStr, err := gmd5.EncryptFile("./a.txt")
		if err != nil {
			g.Log().Line(true).Println(err)
		}
		r.Response.Writeln(encryptStr)
	})

	s.BindHandler("/godmin", func(r *ghttp.Request) {
		password := "lufeit14513"
		salt := "O73OdcQd43"
		encryptPass := EncryptPassword(password, salt)
		g.Log().Line(true).Info(encryptPass)
		r.Response.Writeln(encryptPass)
	})
	s.SetPort(8200)
	s.Start()

	g.Wait()
}

func EncryptPassword(password, salt string) string {
	return gmd5.MustEncryptString(gmd5.MustEncryptString(password) + gmd5.MustEncryptString(salt))
}
