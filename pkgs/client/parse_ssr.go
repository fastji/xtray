package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/moqsien/xtray/pkgs/utils"
)

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

/*
ssr://MTg4LjExOS42NS4xNDM6MTIzMzc6b3JpZ2luOnJjNDpwbGFpbjpiRzVqYmk1dmNtY2diamgwLz9vYmZzcGFyYW09JnJlbWFya3M9NUwtRTU3Mlg1cGF2VHcmZ3JvdXA9VEc1amJpNXZjbWM

ssr://188.119.65.143:12337:origin:rc4:plain:bG5jbi5vcmcgbjh0/?obfsparam=&remarks=5L-E572X5pavTw&group=TG5jbi5vcmc

plain:lncn.org n8t
remarks=俄罗斯
group=Lncn.org
*/
func (that *SSROutbound) parseEncryptMethod(str string) {
	vlist := strings.Split(str, "origin:")
	if len(vlist) == 2 {
		vlist = strings.Split(vlist[1], ":")
		if len(vlist) > 1 {
			that.Method = vlist[0]
		}
	}
}

func (that *SSROutbound) parsePassword(str string) {
	vlist := strings.Split(str, "plain:")
	if len(vlist) == 2 {
		vlist = strings.Split(vlist[1], "/")
		if len(vlist) > 1 {
			that.Method = utils.DecodeBase64(utils.NormalizeSSR(vlist[0]))
		}
	}
}

func (that *SSROutbound) parse(rawUri string) {
	if strings.HasSuffix(rawUri, "ssr://") {
		r := strings.ReplaceAll(rawUri, "ssr://", "")
		r = utils.DecodeBase64(utils.NormalizeSSR(r))
		vlist := strings.SplitN(r, ":", 3)
		if len(vlist) == 3 {
			that.Address = vlist[0]
			that.Port, _ = strconv.Atoi(vlist[1])
			that.parseEncryptMethod(r)
			that.parsePassword(r)
		}
	}
}

func (that *SSROutbound) GetConfigStr(rawUri string) (r string) {
	that.parse(rawUri)
	j := gjson.New(SSRStr)
	j.Set("servers.0.email", that.Email)
	j.Set("servers.0.address", that.Address)
	j.Set("servers.0.port", that.Port)
	j.Set("servers.0.method", that.Method)
	j.Set("servers.0.password", that.Password)
	serverStr := j.MustToJsonIndentString()
	streamStr := "{}"
	confStr := fmt.Sprintf(XrayConfStr, serverStr, streamStr)
	j = gjson.New(confStr)
	j.Set("outbounds.0.protocol", "ssr")
	return j.MustToJsonIndentString()
}
