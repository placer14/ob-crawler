package crawler

import (
	"fmt"
	"net/http"

	"github.com/placer14/ob-crawler/api"
)

type CrawlOptions struct {
	AuthCookie string
	ApiHost    string
	ApiPort    int
}

type Crawler struct {
	api   *api.Client
	queue []string
}

func New(opt CrawlOptions) *Crawler {
	crawler := &Crawler{
		api: &api.Client{
			Host:       opt.ApiHost,
			Port:       opt.ApiPort,
			AuthCookie: opt.AuthCookie,
			Client:     &http.Client{},
		},
	}
	return crawler
}

func (c *Crawler) Execute() error {
	var err error

	fmt.Printf("Beginning crawl from %s:%d...\n", c.api.Host, c.api.Port)
	c.queue, err = c.api.GetPeers()
	if err != nil {
		return err
	}

	fmt.Printf("Found %d Peers: \n %+v\n", len(c.queue), c.queue)
	return nil
}
