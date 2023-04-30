package conf

import (
	"path/filepath"

	"github.com/moqsien/xtray/pkgs/utils"
)

/*
Verifier port range
*/
type VPortRange struct {
	Start int
	End   int
}

type Conf struct {
	FetcherUrl      string      `json:"fetcher_url"`
	WorkDir         string      `json:"work_dir"`
	RawProxyFile    string      `json:"raw_file"`
	PorxyFile       string      `json:"proxy_file"`
	PortRange       *VPortRange `json:"port_range"`
	Port            int         `json:"port"`
	TestUrl         string      `json:"test_url"`
	SwitchyOmegaUrl string      `json:"omega_url"`
	Timeout         int         `json:"timeout"`
}

var DefaultWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/proxy_files")

func NewConf() (conf *Conf) {
	conf = &Conf{}
	conf.WorkDir = DefaultWorkDir
	conf.FetcherUrl = "https://gitee.com/moqsien/test/raw/master/conf.txt"
	conf.RawProxyFile = filepath.Join(conf.WorkDir, "raw_proxy.json")
	conf.PorxyFile = filepath.Join(conf.WorkDir, "latest.json")
	conf.PortRange = &VPortRange{2020, 2150}
	conf.Port = 2019
	conf.TestUrl = "https://www.google.com"
	conf.SwitchyOmegaUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/switch-omega.zip"
	conf.Timeout = 3
	return
}
