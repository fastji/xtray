package client

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
)

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
	Security string
	Type     string
	Path     string
}

/*
trojan://b5e4e360-5946-470b-aad0-db98f50faa57@frontend.yijianlian.app:54430?security=tls&type=tcp&headerType=none#%F0%9F%87%BA%F0%9F%87%B8%20Relay%20%F0%9F%87%BA%F0%9F%87%B8%20United%20States%2011%20TG%3A%40SSRSUB
*/
func (that *TrojanOutboud) parse(rawUri string) {
	if strings.HasSuffix(rawUri, "trojan://") {
		if u, err := url.Parse(rawUri); err == nil {
			that.Address = u.Hostname()
			that.Port, _ = strconv.Atoi(u.Port())
			that.Password = u.User.Username()
			that.Security = u.Query().Get("security")
			that.Type = u.Query().Get("type")
			that.Path = u.Query().Get("path")
		}
	}
}

func (that *TrojanOutboud) GetConfigStr(rawUri string) (r string) {
	that.parse(rawUri)
	j := gjson.New(TrojanStr)
	j.Set("servers.0.address", that.Address)
	j.Set("servers.0.port", that.Port)
	j.Set("servers.0.password", that.Password)
	j.Set("servers.0.email", that.Email)
	serverStr := j.MustToJsonIndentString()
	streamStr := "{}"
	if that.Type != "" {
		j = gjson.New(StreamStr)
		j.Set("network", that.Type)
		j.Set("security", that.Security)
		j.Set("wsSettings.path", that.Path)
		streamStr = j.MustToJsonIndentString()
	}
	confStr := fmt.Sprintf(XrayConfStr, serverStr, streamStr)
	j = gjson.New(confStr)
	j.Set("outbounds.0.protocol", "trojan")
	return j.MustToJsonIndentString()
}
