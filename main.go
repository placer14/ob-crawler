package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/placer14/ob-crawler/crawler"
)

func main() {
	var opts = crawler.CrawlOptions{}

	flag.IntVar(&opts.MaxVisits, "max-visits", 0, "`maximum number of visits` the workers will make. 0 will allow the crawler to traverse the entire network")
	flag.IntVar(&opts.ApiTimeout, "api-timeout", 60, "`time in seconds` to wait before abandoning a request")
	flag.IntVar(&opts.ApiPort, "api-port", 4002, "`port` to use when connecting to OpenBazaar API")
	flag.IntVar(&opts.WorkerPoolSize, "n", 10, "`number of concurrent crawlers` making API requests")
	flag.StringVar(&opts.ApiHost, "api-host", "api", "`host` to use when connecting to OpenBazaar API")
	flag.StringVar(&opts.AuthCookie, "auth-cookie", "", ".cookie `content` generated in OpenBazaar data path")
	flag.Parse()

	crawler := crawler.New(opts)
	err := crawler.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Complete.\nFound %d listings across %d nodes.\n", crawler.ListingCount(), crawler.NodesVisited())
	os.Exit(0)
}
