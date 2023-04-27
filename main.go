package main

import (
	"time"

	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/proxy"
)

func init() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
}

func main() {
	config := conf.NewConf()
	fetcher := proxy.NewFetcher(config)
	fetcher.GetFile()
}
