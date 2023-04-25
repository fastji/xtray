package conf

type Conf struct {
	FetcherUrl   string `json:"fetcher_url"`
	WorkDir      string `json:"work_dir"`
	RawProxyFile string `json:"raw_file"`
	PorxyFile    string `json:"proxy_file"`
}
