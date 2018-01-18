package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/placer14/ob-crawler/crawler"
)

func main() {
	var opts = crawler.CrawlOptions{}

	flag.IntVar(&opts.ApiTimeout, "api-timeout", 30, "`time in seconds` to wait before abandoning a request")
	flag.IntVar(&opts.ApiPort, "api-port", 4002, "`port` to use when connecting to OpenBazaar API")
	flag.IntVar(&opts.WorkerPoolSize, "n", 20, "`number of concurrent crawlers` making API requests")
	flag.StringVar(&opts.ApiHost, "api-host", "api", "`host` to use when connecting to OpenBazaar API")
	flag.StringVar(&opts.AuthCookie, "auth-cookie", "", ".cookie `content` generated in OpenBazaar data path")
	flag.Parse()

	crawler := crawler.New(opts)
	err := crawler.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
