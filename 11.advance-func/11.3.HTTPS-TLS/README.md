# 一、下载安装certbot: https://certbot.eff.org/lets-encrypt/ubuntuxenial-nginx
## 1.1 Ubuntu16.04LTS 下载安装 snapd
sudo apt-get install snapd

## 1.2 将snapd更新到最新版本
sudo snap install core; sudo snap refresh core

## 1.3 下载安装certbot
sudo snap install --classic certbot

## 1.4 设置全局的certbot命令
sudo ln -s /snap/bin/certbot /usr/bin/certbot

# 二、配置Nginx
## 2.1 首先需要明确一点的是，如果不使用nginx进行请求转发，直接获取证书，使用：
certbot certonly --standalone -d 申请域名 --staple-ocsp -m 邮箱地址 --agree-tos

但是这会有一个问题，就是这样会默认占用80端口（如果系统中有nginx，则会产生端口冲突，因为端口只能被一个进程占用）;
由于做前后端开发难免使用nginx，所以推荐下面的方式证书的获取和安装

## 2.2 通过certbot自动配置nginx
* 首先，在/etc/nginx/sites-enabled目录下添加一个gf.conf文件，然后添加如下配置:
```text
server {
    server_name  gf.lcoderfit.com;

    charset utf-8;
    access_log /etc/nginx/sites-enabled/nohup.out;

    location / {
        proxy_pass https://localhost:8201;
    }
}
```
* 自动添加HTTPS配置到gf.conf 
sudo certbot --nginx
该命令会扫描所有的/etc/nginx/sites-enabled目录下的配置文件中配置的域名，然后列出一个域名列表如下：

![1631699604537](D:\PrivateProject\Gf-Tags\learn-gf\11.advance-func\11.3.HTTPS-TLS\域名列表.png)

* 选择数字9然后回车，就会自动生成证书和私钥，并且gf.conf中的配置也会发生变化（多了几行由certbot自动配置的部分）
Certificate is saved at: /etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem
Key is saved at:         /etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem

* 查看gf.conf
注意:这里巨坑,因为阿里云ECS需要开启安全组,而443端口默认是未开启的,访问:
https://ecs.console.aliyun.com/#/securityGroupDetail/region/cn-shanghai/groupId/sg-uf69ibpdra8ugo0gg3dt/detail/intranetIngress
```text
# listen可以写在下面的(这个是可以正确识别的,不要以为只能写server_name xxx上面)
server {
    server_name  gf.lcoderfit.com;

    charset utf-8;
    access_log /etc/nginx/sites-enabled/nohup.out;

    location / {
        # 注意：这里也可以使用 http://localhost:8200,将HTTPS请求转发到HTTP服务
        proxy_pass https://localhost:8201;
    }

    # 以下部分是由certbot自动生成
    # 这个443端口是可以修改的，HTTPS默认端口443，这样前端访问时不需要显式的写端口
    # 如果需要配置多个HTTPS的配置，则这里的端口可以换成其他端口，避免冲突
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}
server {
    # 这个表示当你请求的url中的域名部分为gf.lcoderfit.com时,会重定向到以https协议开头的url(301表示永久重定向)
    # $host遍历表示域名,$request_uri表示域名后的url部分
    # 这一部分可以自定义配置,例如访问http://gf.lcoderfit.com时,将请求转发到另一个端口
    if ($host = gf.lcoderfit.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    listen       80;
    server_name  gf.lcoderfit.com;
    return 404; # managed by Certbot
    # 可以将上面的return 404 和 if($host ...)的部分注释掉，然后将下面的注释打开，就可以实现HTTPS请求转发到https://locahost:8201
    # HTTP请求转发到http://localhost:8200; 
    # location / {
    #     proxy_pass http://localhost:8200
    # }
}
```
# 三、goframe后端程序
注意：由于服务器上80和443端口都被nginx占用了，所以后端接口的端口需要需改成其他端口，例如下面HTTP服务设置成8200端口，HTTPS服务设置成8201端口
查看端口开放情况:
    1. sudo apt-get install net-tools 安装netstat工具
    2. sudo netstat -tlpn (t表示tcp，l表示listening，p表示是否显示进程pid，如果不带n，则0.0.0.0会显示为*)
    3.下次访问不同，先看一下是否是端口未开放的问题(去阿里云安全组界面开放端口)
    4.
```text
package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Writeln("可以同时通过HTTPS和HTTP访问")
	})
    // 证书和密钥的路径
	s.EnableHTTPS(
		"/etc/letsencrypt/live/gf.lcoderfit.com/fullchain.pem",
		"/etc/letsencrypt/live/gf.lcoderfit.com/privkey.pem",
	)
	s.SetPort(8200)
	s.SetHTTPSPort(8201)
	s.Start()

	g.Wait()
}
```

# 四、证书续期
* certbot申请的证书有效期为3个月，所以需要定时续期，可以使用corntab命令
```text
1.Ubuntu16.04LTS下载cron工具
sudo apg-get install cron

2.启动cron服务
service cron start

3.crontab命令
crontab -e 编辑定时任务
crontab -l 显示定时任务
```

* crontab命令如下
格式类比英文时间格式，分钟 小时 一月中的某一天(1-31) 月份(1-12) 星期(0-6,0为星期天)
![1631707410686](D:\PrivateProject\Gf-Tags\learn-gf\11.advance-func\11.3.HTTPS-TLS\crontab命令.png)

* 续签certbot SSL证书的命令
```text
# certbot证书有效期为3个月，所以设置每一个月重新续期一次
0 3 1 * * service nginx stop && certbot renew --quiet --renew-hook "service nginx start"
# 为了防止certbot renew命令失败导致nginx无法重启，多设置一个保险
30 3 1 * * service nginx start
```

* 免费证书请参考
https://goframe.org/pages/viewpage.action?pageId=1114278

![1631707862873](D:\PrivateProject\Gf-Tags\learn-gf\11.advance-func\11.3.HTTPS-TLS\SSL免费证书.png)