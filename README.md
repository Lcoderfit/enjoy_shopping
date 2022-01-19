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

5.注意:MySQL不能使用update t update a=1 and b=2 where 这种语法, and的语法没有用

6.save: Error 1406: Data too long for column 'name' at row 1, INSERT INTO `user`(`uid`,`name`) VALUES(2,'tianyi-save') 
ON DUPLICATE KEY UPDATE `uid`=VALUES(`uid`),`name`=VALUES(`name`)
```text
Save()插入的数据超过了字段的范围
```

7.主键或唯一索引重复
    -- 这个表示插入的数据总的name字段存在重复(name字段是一个唯一索引)
    ERROR 1062 (23000): Duplicate entry 'tianyi' for key 'name'
    -- 插入数据时存在主键重复,重复值为2
    ERROR 1062 (23000): Duplicate entry '2' for key 'PRIMARY'
```text
当存在主键或唯一索引重复时则更新:
    -- `name`=values(`name`)可以动态设置字段name的值,即如果存在主键或唯一索引冲突,更新name的值为values('2', 'tianyi')中name字段的值
    insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`=values(`uid`),`name`=values(`name`);

    -- 插入('2', 'tianyi')时,如果存在重复数据,则更新uid为2, name为'lh'
    insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`='2',`name`='lh';
    -- 注意,如果表中已经同时有uid为2和3的数据,则下面的语句仍会报错,因为插入uid为2的数据会导致主键冲突,此时会更新这条uid为2的数据为('3', 'lh');

    -- 但是uid=3的数据也是存在的,所以仍会: Duplicate entry '3' for key 'PRIMARY'
    insert into `user`(`uid`,`name`) values('2', 'tianyi') on duplicate key update `uid`='3',`name`='lh';
```

8.ERROR 1054 (42S22): Unknown column 'tianyi' in 'field list'
```text
-- 在插入数据时将值'tianyi'两边用``括起来，导致MySQL识别成了字段
insert ignore into `user`(`name`) values(`tianyi`);
```