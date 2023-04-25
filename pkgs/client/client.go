package client

import (
	"strings"

	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/confloader/external"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type IOutbound interface {
	GetConfigStr(string) string
	GetRawUri() string
}

type XClient struct {
	*core.Instance
	RawUri string
	Out    IOutbound
}

func NewXClient(rawUri string) *XClient {
	var out IOutbound
	if strings.HasPrefix("vmess://", rawUri) {
		out = &VmessOutbound{}
	} else if strings.HasPrefix("vless://", rawUri) {
		out = &VlessOutbound{}
	} else if strings.HasPrefix("ss://", rawUri) {
		out = &SSOutbound{}
	} else if strings.HasPrefix("ssr://", rawUri) {
		out = &SSROutbound{}
	} else if strings.HasPrefix("trojan://", rawUri) {
		out = &TrojanOutboud{}
	} else {
		return nil
	}
	return &XClient{
		RawUri: rawUri,
		Out:    out,
	}
}
