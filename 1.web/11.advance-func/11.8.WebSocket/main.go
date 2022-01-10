package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

/*
1.将HTTP请求转换为websocket
	ws, err := r.WebSocket()

2.ws和wss
	2.1 支持HTTP的websocket ws://localhost:port/xxx
	2.2 支持HTTPS的WebSocket wss://localhost:port/xxx

3.读取数据和写入数据
	// 返回的msg为[]byte类型
	msgType, msg, err := ws.ReadMessage()
	ws.WriteMessage(msgType, msg)
*/

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("index_ws.html", nil)
	})
	s.BindHandler("/ws", func(r *ghttp.Request) {
		ws, err := r.WebSocket()
		if err != nil {
			glog.Error(err)
			r.Exit()
		}
		for {
			msgType, msg, err := ws.ReadMessage()
			g.Log().Println(msgType)

			if err != nil {
				return
			}
			if err = ws.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})
	s.SetPort(8300)
	s.Start()

	s1 := g.Server("s1")
	s1.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("index_wss.html", nil)
	})
	s1.BindHandler("/wss", func(r *ghttp.Request) {
		wss, err := r.WebSocket()
		if err != nil {
			glog.Line(true).Error(err)
			r.Exit()
		}
		for {
			msgType, msg, err := wss.ReadMessage()
			g.Log().Println(msgType)
			if err != nil {
				return
			}
			if err = wss.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})
	s1.EnableHTTPS("server/server.crt", "server/server.key")
	s1.SetHTTPSPort(8301)
	s1.SetPort(8302)
	s1.Start()

	g.Wait()
}