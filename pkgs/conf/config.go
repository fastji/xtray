package conf

import (
	"path/filepath"

	"github.com/moqsien/xtray/pkgs/utils"
)

type Conf struct {
	FetcherUrl   string `json:"fetcher_url"`
	WorkDir      string `json:"work_dir"`
	RawProxyFile string `json:"raw_file"`
	PorxyFile    string `json:"proxy_file"`
}

var DefaultWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/proxy_files")

func NewConf() (conf *Conf) {
	conf = &Conf{}
	conf.WorkDir = DefaultWorkDir
	conf.FetcherUrl = "https://gitee.com/moqsien/test/raw/master/conf.txt"
	conf.RawProxyFile = filepath.Join(conf.WorkDir, "raw_proxy.json")
	conf.PorxyFile = filepath.Join(conf.WorkDir, "latest.json")
	return
}
