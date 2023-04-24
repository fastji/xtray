package client

/*
	"servers": [{
		"address": "127.0.0.1",
		"port": 1234,
		"password": "password",
		"email": "love@xray.com",
		"level": 0
	}]
*/
var TrojanStr string = `{
	"servers": [{
		"address": "127.0.0.1",
		"port": 1234,
		"password": "password",
		"email": "love@xray.com",
		"level": 0
	}]
}`

type TrojanOutboud struct {
	Address  string
	Port     int
	Password string
	Email    string
}
