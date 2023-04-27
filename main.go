package main

import (
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/proxy"
)

func main() {
	config := conf.NewConf()
	fetcher := proxy.NewFetcher(config)
	fetcher.GetFile()
}
