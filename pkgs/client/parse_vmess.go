package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/moqsien/xtray/pkgs/utils"
)

/*
	"vnext": [{
		"address": "127.0.0.1",
		"port": 37192,
		"users": [{
			"id": "5783a3e7-e373-51cd-8642-c83782b807c5",
			"alterId": 0,
			"security": "auto",
			"level": 0
		}]
	}]
*/
var VmessStr string = `{
	"vnext": [{
		"address": "127.0.0.1",
		"port": 37192,
		"users": [{
			"id": "5783a3e7-e373-51cd-8642-c83782b807c5",
			"alterId": 0,
			"security": "none",
			"level": 0
		}]
	}]
}`

type VmessOutbound struct {
	Address  string `json:"address"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserId   string `json:"id"`
	Network  string `json:"network"`
	Security string `json:"security"`
	Path     string `json:"path"`
	Raw      string `json:"raw"`
}

/*
vmess://eyJ2IjogIjIiLCAicHMiOiAiZ2l0aHViLmNvbS9mcmVlZnEgLSBcdTdmOGVcdTU2ZmRDbG91ZEZsYXJlXHU1MTZjXHU1M2Y4Q0ROXHU4MjgyXHU3MGI5IDEiLCAiYWRkIjogIm1pY3Jvc29mdGRlYnVnLmNvbSIsICJwb3J0IjogIjgwIiwgImlkIjogIjEwMTdlZjZhLTY3ZDktNGJiMy1iNjY3LTBkNjdjMWVlNTU0NiIsICJhaWQiOiAiMCIsICJzY3kiOiAiYXV0byIsICJuZXQiOiAid3MiLCAidHlwZSI6ICJub25lIiwgImhvc3QiOiAidjEudXM5Lm1pY3Jvc29mdGRlYnVnLmNvbSIsICJwYXRoIjogIi9zZWN4IiwgInRscyI6ICIiLCAic25pIjogIiJ9

{"v": "2", "ps": "github.com/freefq - \u7f8e\u56fdCloudFlare\u516c\u53f8CDN\u8282\u70b9 1",
"add": "microsoftdebug.com", "port": "80", "id": "1017ef6a-67d9-4bb3-b667-0d67c1ee5546",
"aid": "0", "scy": "auto", "net": "ws", "type": "none", "host": "v1.us9.microsoftdebug.com",
"path": "/secx", "tls": "", "sni": ""}
*/

func (that *VmessOutbound) parse(rawUri string) {
	that.Raw = rawUri
	if strings.HasPrefix(rawUri, "vmess://") {
		rawUri = strings.ReplaceAll(rawUri, "vmess://", "")
	}
	rawUri = utils.DecodeBase64(rawUri)
	j := gjson.New(rawUri)
	that.Address = j.GetString("add")
	that.Host = j.GetString("host")
	that.Port, _ = strconv.Atoi(j.GetString("port"))
	that.UserId = j.GetString("id")
	that.Network = j.GetString("net")
	if that.Network == "" {
		that.Network = "tcp"
	}
	that.Security = j.GetString("tls")
	if that.Security == "" {
		that.Security = "none"
	}
	that.Path = j.GetString("path")
}

func (that *VmessOutbound) GetConfigStr(rawUri string) string {
	that.parse(rawUri)
	j := gjson.New(VlessStr)
	j.Set("vnext.0.address", that.Address)
	j.Set("vnext.0.port", that.Port)
	j.Set("vnext.0.users.0.id", that.UserId)
	vnextStr := j.MustToJsonIndentString()
	j = gjson.New(StreamStr)
	j.Set("network", that.Network)
	j.Set("security", that.Security)
	j.Set("wsSettings.path", that.Path)
	streamStr := j.MustToJsonIndentString()
	confStr := fmt.Sprintf(XrayConfStr, vnextStr, streamStr)
	j = gjson.New(confStr)
	j.Set("outbounds.0.protocol", "vmess")
	return j.MustToJsonIndentString()
}

func (that *VmessOutbound) GetRawUri() string {
	return that.Raw
}

func (that *VmessOutbound) GetString() string {
	return fmt.Sprintf("vmess://%s:%d", that.Address, that.Port)
}

func TestVmess(rawUri string) {
	v := &VmessOutbound{}
	fmt.Println(v.GetConfigStr(rawUri))
}
