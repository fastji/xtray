package client

/*
	"vnext": [{
		"address": "example.com",
		"port": 443,
		"users": [{
			"id": "5783a3e7-e373-51cd-8642-c83782b807c5",
			"encryption": "none",
			"flow": "none",
			"level": 0
		}]
	}]
*/
var VlessStr string = `{
	"vnext": [{
		"address": "example.com",
		"port": 443,
		"users": [{
			"id": "5783a3e7-e373-51cd-8642-c83782b807c5",
			"encryption": "none",
			"flow": "none",
			"level": 0
		}]
	}]
}`

type VlessOutbound struct {
	Address  string
	Port     int
	UserId   string
	Security string
}

/*
vless://b1e41627-a3e9-4ebd-9c92-c366dd82b13f@xray.ibgfw.top:2083?encryption=none&security=tls&type=ws&host=&path=/wSXCvstU/#xray.ibgfw.top%3A2083
*/
func (that *VlessOutbound) parse(rawUri string) {}

func (that *VlessOutbound) GetConfigStr(rawUri string) string {
	that.parse(rawUri)
	return ""
}
