package client

/*
	"servers": [{
		"address": "127.0.0.1",
		"port": 1234,
		"users": [{
			"user": "test user",
			"pass": "test pass",
			"level": 0
		}]
	}]
*/
var SSStr string = `{
	"servers": [{
		"address": "127.0.0.1",
		"port": 1234,
		"users": [{
			"user": "test user",
			"pass": "test pass",
			"level": 0
		}]
	}]
}`

type SSOutbound struct {
	Address  string
	Port     int
	Username string
	Password string
}
