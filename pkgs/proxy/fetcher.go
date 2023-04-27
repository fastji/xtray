package proxy

import (
	"os"

	"github.com/gocolly/colly/v2"
	futil "github.com/moqsien/free/pkgs/utils"
	"github.com/moqsien/xtray/pkgs/conf"
)

type Fetcher struct {
	collector *colly.Collector
	conf      *conf.Conf
}

func NewFetcher(c *conf.Conf) *Fetcher {
	return &Fetcher{
		collector: colly.NewCollector(),
		conf:      c,
	}
}

func (that *Fetcher) GetFile() {
	that.collector.OnResponse(func(r *colly.Response) {
		if result, err := futil.DefaultCrypt.AesDecrypt(r.Body); err == nil {
			os.WriteFile(that.conf.RawProxyFile, result, os.ModePerm)
		}
	})
	that.collector.Visit(that.conf.FetcherUrl)
}
