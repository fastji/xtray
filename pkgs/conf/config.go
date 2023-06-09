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
	GeoInfoUrl      string      `json:"geo_info_url"`
	Timeout         int         `json:"timeout"`
	VerifierCron    string      `json:"verifier_cron"`
	KeeperCron      string      `json:"keeper_cron"`
}

var DefaultWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/proxy_files")

func NewConf() (conf *Conf) {
	conf = &Conf{}
	conf.WorkDir = DefaultWorkDir
	conf.FetcherUrl = "https://gitee.com/moqsien/test/raw/master/conf.txt"
	conf.RawProxyFile = filepath.Join(conf.WorkDir, "raw_proxy.json")
	conf.PorxyFile = filepath.Join(conf.WorkDir, "latest.json")
	conf.PortRange = &VPortRange{2020, 2075}
	conf.Port = 2019
	conf.TestUrl = "https://www.google.com"
	conf.SwitchyOmegaUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/switch-omega.zip"
	conf.GeoInfoUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/geoinfo.zip"
	conf.Timeout = 3
	// "@every 1h30m10s" https://pkg.go.dev/github.com/robfig/cron
	conf.VerifierCron = "@every 2h"
	conf.KeeperCron = "@every 3m"
	return
}
