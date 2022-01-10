# LearnGf

#### 介绍
gf学习用例

#### 软件架构
软件架构说明


#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx
4.一次性删除所有项目下的二进制文件, 注意，xargs是遍历前一个输入的多行，没遍历到一行就将文件路径传递给rm -rf命令
find ./ -name gf -o -name main -o -name main.exe* |xargs rm rf

5.在.gitignore中,*.exe只会过滤当前路径下的.exe后缀的文件
如果需要过滤所有子目录中带哟.exe后缀的文件,用 **表示:
```text
**/*.exe
**/*.exe~
**/main
```
#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)


#### 错误集合
1.Only one usage of each socket address (protocol/network address/port
) is normally permitted.
```text
每个套接字地址（协议/网络地址/端口）只使用一次)通常是允许的。
解决：这个一般是有多个程序占用同一端口，找到并杀掉多余的程序即可
```

2.This version of %1 is not compatible with the version of Windows you're runni
ng. Check your computer's system information and then contact the software publisher
```text
包名不是main包
```

3.reflect: call of reflect.Value.NumField on slice Value
```text
定义接受查询的返回数据的结构体时,错误使用成: var res *[]StructName 如果是存储列表,不应该定义成一个切片指针,而是: 
var res []*StructName
```

4.Scan called without calling Next
```text

```
