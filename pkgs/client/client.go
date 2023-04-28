package client

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/moqsien/xtray/pkgs/utils"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf/serial"
	_ "github.com/xtls/xray-core/main/confloader/external"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type ClientParams struct {
	RawUri string
	InPort int
}

type IOutbound interface {
	GetConfigStr(string) string
	GetRawUri() string
}

type XClient struct {
	*core.Instance
	RawUri  string
	Out     IOutbound
	ConfStr string
}

func NewXClient() *XClient {
	return &XClient{}
}

func (that *XClient) setOutbound(rawUri string) {
	that.RawUri = strings.TrimSpace(rawUri)
	if strings.HasPrefix(rawUri, "vmess://") {
		that.Out = &VmessOutbound{}
	} else if strings.HasPrefix(rawUri, "vless://") {
		that.Out = &VlessOutbound{}
	} else if strings.HasPrefix(rawUri, "ss://") {
		that.Out = &SSOutbound{}
	} else if strings.HasPrefix(rawUri, "ssr://") {
		that.Out = &SSROutbound{}
	} else if strings.HasPrefix(rawUri, "trojan://") {
		that.Out = &TrojanOutboud{}
	} else {
		fmt.Println("Unsupported vpn uri: ", rawUri)
		that.Out = nil
	}
}

func (that *XClient) Start(params *ClientParams) error {
	that.setOutbound(params.RawUri)
	if that.Out == nil {
		return errors.New("illegal Outbound, please check the uri")
	} else if params.InPort == 0 {
		return errors.New("illegal inbound port")
	}
	confStr := that.Out.GetConfigStr(that.RawUri)
	j := gjson.New(confStr)
	j.Set("inbounds.0.port", params.InPort)
	that.ConfStr = j.MustToJsonIndentString()

	if config, err := serial.DecodeJSONConfig(utils.StringToReader(that.ConfStr)); err == nil {
		var f *core.Config
		f, err = config.Build()
		if err != nil {
			fmt.Println("[Build config for Xray failed] ", err)
			return err
		}
		that.Instance, err = core.New(f)
		if err != nil {
			fmt.Println("[Init Xray Instance Failed] ", err)
			return err
		}
		that.Instance.Start()
	} else {
		fmt.Println("[Start Client Failed] ", err)
		return err
	}
	return nil
}

func (that *XClient) Close() {
	if that.Instance != nil {
		that.Instance.Close()
		that.Instance = nil
		runtime.GC()
	}
}
