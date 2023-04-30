package client

import "strings"

func ParseRawUri(rawUri string) string {
	var parser IOutbound
	if strings.HasPrefix(rawUri, "vmess://") {
		parser = &VmessOutbound{}
	} else if strings.HasPrefix(rawUri, "vless://") {
		parser = &VlessOutbound{}
	} else if strings.HasPrefix(rawUri, "ss://") {
		parser = &SSOutbound{}
	} else if strings.HasPrefix(rawUri, "ssr://") {
		parser = &SSROutbound{}
	} else if strings.HasPrefix(rawUri, "trojan://") {
		parser = &TrojanOutboud{}
	} else {
		return ""
	}
	parser.GetConfigStr(rawUri)
	return parser.GetString()
}
