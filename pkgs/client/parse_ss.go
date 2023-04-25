package client

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/moqsien/xtray/pkgs/utils"
)

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
	Raw      string
}

/*
ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTo3MjgyMjliOS0xNjRlLTQ1Y2ItYmZiMy04OTZiM2EwNTZhMTg=@node03.gde52px1vwf5q6301fxn.catapi.management:33907#%F0%9F%87%AC%F0%9F%87%A7%20Relay%20%F0%9F%87%AC%F0%9F%87%A7%20United%20Kingdom%2005%20TG%3A%40SSRSUB

ss://chacha20-ietf-poly1305:728229b9-164e-45cb-bfb3-896b3a056a18@node03.gde52px1vwf5q6301fxn.catapi.management:33907
*/
func (that *SSOutbound) parse(rawUri string) {
	that.Raw = rawUri
	if strings.Contains(rawUri, "ss://") {
		r := strings.ReplaceAll(rawUri, "ss://", "")
		uList := strings.Split(r, "@")
		if len(uList) == 2 {
			userInfo := utils.DecodeBase64(utils.NormalizeSSR(uList[0]))
			_uri := fmt.Sprintf("ss://%s@%s", userInfo, uList[1])
			if u, err := url.Parse(_uri); err == nil {
				that.Address = u.Hostname()
				that.Port, _ = strconv.Atoi(u.Port())
				that.Username = u.User.Username()
				that.Password, _ = u.User.Password()
			}
		}
	}
}

func (that *SSOutbound) GetConfigStr(rawUri string) (r string) {
	that.parse(rawUri)
	j := gjson.New(SSStr)
	j.Set("servers.0.address", that.Address)
	j.Set("servers.0.port", that.Port)
	j.Set("servers.0.users.0.user", that.Username)
	j.Set("servers.0.users.0.pass", that.Password)
	serverStr := j.MustToJsonIndentString()
	streamStr := "{}"
	confStr := fmt.Sprintf(XrayConfStr, serverStr, streamStr)
	j = gjson.New(confStr)
	j.Set("outbounds.0.protocol", "ss")
	return j.MustToJsonIndentString()
}

func (that *SSOutbound) GetRawUri() string {
	return that.Raw
}
