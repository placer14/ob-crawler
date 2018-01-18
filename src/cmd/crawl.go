package cmd

import (
	"github.com/placer14/ob-crawler/crawler"
)

func Crawl(opts crawler.CrawlOptions) error {
	crawler := crawler.New(opts)
	return crawler.Execute()
}
