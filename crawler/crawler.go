package crawler

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/placer14/ob-crawler/api"
)

// CrawlOptions is a struct of settings expected by the Crawler and used by
// cmd.Crawl when beginning a new crawl operation
type CrawlOptions struct {
	AuthCookie     string
	ApiHost        string
	ApiPort        int
	ApiTimeout     int
	WorkerPoolSize int
	MaxVisits      int
}

// Crawler is the engine which coordinates the workers when crawling the OB API
type Crawler struct {
	api            api.OpenBazaarlikeAPI
	queue          []string
	workerQueue    chan string
	workerPoolSize int
	workersActive  *sync.WaitGroup

	cacheMutex  *sync.Mutex
	lookupCache map[string]*nodeData

	listingCount         int
	maximumVisitsAllowed int
}

// New returns a configured Crawler based on the CrawlOptions provided in opt
func New(opt CrawlOptions) *Crawler {
	crawler := &Crawler{
		api: &api.Client{
			Host:       opt.ApiHost,
			Port:       opt.ApiPort,
			AuthCookie: opt.AuthCookie,
			Client: &http.Client{
				Timeout: (time.Duration(opt.ApiTimeout) * time.Second),
			},
		},
		cacheMutex:           &sync.Mutex{},
		workerPoolSize:       opt.WorkerPoolSize,
		workersActive:        &sync.WaitGroup{},
		maximumVisitsAllowed: opt.MaxVisits,
	}
	return crawler
}

// CountedListings returns the number of listed contracts found while crawling
// the target. This will return a safe value while Execute() is running

// Execute begins the crawl operation by connecting the API via the Host and
// Port provided in CrawlOptions in the argument to cmd.Crawl to retrieve the API's
// closest peers, walking outward to each peer counting the number of listings on
// each node found.
//
// Each node is only visited once regardless of successful count of listings.
func (c *Crawler) Execute() error {
	var err error

	fmt.Printf("Beginning crawl from %s...\n", c.api.HostPort())

	c.queue, err = c.api.GetPeers()
	if err != nil {
		return err
	}

	fmt.Printf("Found %d Seed Peers...\n", len(c.queue))
	c.workerQueue = make(chan string)
	c.lookupCache = make(map[string]*nodeData)

	for id := 1; id <= c.workerPoolSize; id++ {
		go c.startWorker(id)
	}
	c.assignJobs()
	return nil
}

// ListingCount returns the number of found listings while crawling the
// OpenBazaar network. It is safe to read this value while Execute is in progress
func (c *Crawler) ListingCount() (count int) {
	count = c.listingCount
	return
}

// NodesVisited returns the number of nodes which were attemped queried.
// It is safe to read this value while Execute is in progress
func (c *Crawler) NodesVisited() (count int) {
	count = len(c.lookupCache)
	return
}

func (c *Crawler) assignJobs() {
	for {
		if len(c.queue) == 0 {
			c.workersActive.Wait()
			if len(c.queue) == 0 {
				close(c.workerQueue)
				break
			}
		}

		c.cacheMutex.Lock()
		var nextHash = c.queue[0]
		c.queue = c.queue[1:]
		if c.maximumVisitsAllowed > 0 && len(c.lookupCache) >= c.maximumVisitsAllowed {
			c.cacheMutex.Unlock()
			continue
		}
		if _, exists := c.lookupCache[nextHash]; exists {
			c.cacheMutex.Unlock()
			continue
		}
		c.lookupCache[nextHash] = &nodeData{}
		c.cacheMutex.Unlock()
		c.workerQueue <- nextHash
		fmt.Printf("Assigning %s, %d visited, %d remaining, %d listings...\n", nextHash, c.NodesVisited(), len(c.queue), c.listingCount)
	}
}

type nodeData struct {
	peers        []string
	listingCount int
}

func (c *Crawler) startWorker(id int) {
	if _, ok := c.api.(*api.Client); ok {
		// Only delay when using the real Client
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

	for node := range c.workerQueue {
		c.workersActive.Add(1)
		peers, err := c.api.GetClosestPeers(node)
		if err != nil {
			peers = make([]string, 0)
			fmt.Printf("  worker %d: error (%s): %s\n", id, node, err)
		}
		count, err := c.api.GetListingsCount(node)
		if err != nil {
			fmt.Printf("  worker %d: error (%s): %s\n", id, node, err)
		}

		data := &nodeData{
			peers:        peers,
			listingCount: count,
		}
		c.cacheMutex.Lock()
		c.queue = append(c.queue, peers...)
		c.lookupCache[node] = data
		c.listingCount += count
		c.cacheMutex.Unlock()
		fmt.Printf("  worker %d: completed %s\n", id, node)
		c.workersActive.Done()
	}
}
