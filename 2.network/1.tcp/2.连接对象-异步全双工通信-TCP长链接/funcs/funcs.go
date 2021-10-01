package funcs

import (
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/gtcp"
	"learn-gf/2.network/1.tcp/2.连接对象-异步全双工通信-TCP长链接/types"
)

// 设置data参数为不定参的作用：可以设置当不传入data参数时的默认参数
func SendPkg(conn *gtcp.Conn, act string, data ...string) error {
	s := ""
	if len(data) > 0 {
		s = data[0]
	}

	msg, err := json.Marshal(types.Msg{
		Act:  act,
		Data: s,
	})
	if err != nil {
		panic(err)
	}
	return conn.SendPkg(msg)
}

func RecvPkg(conn *gtcp.Conn) (msg *types.Msg, err error) {
	data, err := conn.RecvPkg()
	if err != nil {
		return nil, err
	} else {
		msg = &types.Msg{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			return nil, gerror.Newf("invalid package structure:%s", err.Error())
		}
		return msg, err
	}
}
