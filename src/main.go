package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/placer14/ob-crawler/cmd"
	"github.com/placer14/ob-crawler/crawler"
)

func main() {
	opts := crawler.CrawlOptions{}
	flag.StringVar(&opts.ApiHost, "api-host", "0.0.0.0", "host to use when connecting to OpenBazaar API")
	flag.IntVar(&opts.ApiPort, "api-port", 4002, "port to use when connecting to OpenBazaar API")
	flag.StringVar(&opts.AuthCookie, "auth-cookie", "", ".cookie content generated in OpenBazaar data path")
	flag.Parse()

	err := cmd.Crawl(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
