package client

/*
	"servers": [{
		"email": "love@xray.com",
		"address": "127.0.0.1",
		"port": 1234,
		"method": "加密方式",
		"password": "密码",
		"level": 0
	}]
*/
var SSRStr string = `{
	"servers": [{
		"email": "love@xray.com",
		"address": "127.0.0.1",
		"port": 1234,
		"method": "加密方式",
		"password": "密码",
		"level": 0
	}]
}`

type SSROutbound struct {
	Email    string
	Address  string
	Port     int
	Method   string
	Password string
}
