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

var WorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/proxy_files")
