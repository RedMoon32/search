package crawler

import (
	"log"
	"net/url"
	"sync"

	"web-crawler/broker"
	"web-crawler/models"
	"web-crawler/utils"
)

type Parser interface {
	Parse(content string) (parsedData map[string]interface{}, newURLs []string)
	AddTask(task models.ParseTask)
}

type RobotsStorage interface {
	Allowed(url string) bool
}

type Crawler struct {
	broker        broker.Broker
	robotsStorage RobotsStorage
	parser        Parser
}

func NewCrawler(b broker.Broker, r RobotsStorage, p Parser) *Crawler {
	return &Crawler{
		broker:        b,
		robotsStorage: r,
		parser:        p,
	}
}

func (c *Crawler) Start() error {
	msgs, err := c.broker.Consume()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for msg := range msgs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c.processMessage(url)
		}(msg)
	}

	wg.Wait()
	return nil
}

func (c *Crawler) processMessage(rawURL string) {
	log.Println("Crawling URL:", rawURL)

	parsedURL, err := url.Parse(rawURL)
	url := parsedURL.Host

	if !c.robotsStorage.Allowed(url) {
		log.Printf("URL %s is disallowed by robots.txt", url)
		return
	}

	content, err := utils.FetchURL(url)
	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", url, err)
		return
	}

	// Send task to parsing channel
	c.parser.AddTask(models.NewParseTask(url, content))
}
