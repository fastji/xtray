package client

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gogf/gf/encoding/gjson"
)

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
	Address    string
	Port       int
	UserId     string
	Security   string
	Encryption string
	Type       string
	Path       string
	Raw        string
}

/*
vless://b1e41627-a3e9-4ebd-9c92-c366dd82b13f@xray.ibgfw.top:2083?encryption=none&security=tls&type=ws&host=&path=/wSXCvstU/#xray.ibgfw.top%3A2083
*/
func (that *VlessOutbound) parse(rawUri string) {
	that.Raw = rawUri
	if r, err := url.Parse(rawUri); err == nil && r.Scheme == "vless" {
		that.Address = r.Hostname()
		that.Port, _ = strconv.Atoi(r.Port())
		that.UserId = r.User.Username()
		that.Security = r.Query().Get("security")
		that.Encryption = r.Query().Get("encryption")
		that.Type = r.Query().Get("type")
		that.Path = r.Query().Get("path")
	}
}

func (that *VlessOutbound) GetConfigStr(rawUri string) string {
	that.parse(rawUri)
	j := gjson.New(VlessStr)
	j.Set("vnext.0.address", that.Address)
	j.Set("vnext.0.port", that.Port)
	j.Set("vnext.0.users.0.id", that.UserId)
	j.Set("vnext.0.users.0.encryption", that.Encryption)
	vnextStr := j.MustToJsonIndentString()
	j = gjson.New(StreamStr)
	j.Set("network", that.Type)
	j.Set("security", that.Security)
	j.Set("wsSettings.path", that.Path)
	streamStr := j.MustToJsonIndentString()
	confStr := fmt.Sprintf(XrayConfStr, vnextStr, streamStr)
	j = gjson.New(confStr)
	j.Set("outbounds.0.protocol", "vless")
	return j.MustToJsonIndentString()
}

func (that *VlessOutbound) GetRawUri() string {
	return that.Raw
}
