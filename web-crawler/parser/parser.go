package parser

import (
	"log"
	"web-crawler/broker"
	"web-crawler/models"
	"web-crawler/utils"
)

type KVStore interface {
	Save(key string, value map[string]interface{}) error
	Load(key string) (map[string]interface{}, error)
}

type HTMLParser struct {
	taskQueue chan models.ParseTask
	kvStore   KVStore
}

func NewHTMLParser(kvStore KVStore) *HTMLParser {
	return &HTMLParser{
		taskQueue: make(chan models.ParseTask, 100),
		kvStore:   kvStore,
	}
}

func (p *HTMLParser) Parse(content string) (map[string]interface{}, []string) {
	parsedData := make(map[string]interface{})
	newURLs := []string{}
	// Implement HTML parsing logic here
	// For example, extracting title, metadata, and URLs
	return parsedData, newURLs
}

func (p *HTMLParser) StartParsingWorker(broker broker.Broker) {
	for task := range p.taskQueue {
		log.Printf("Parsing task: %s\n", task.URL)

		// Read file content
		content := task.Content
		if content == "" {
			log.Printf("Failed to parse content from url %s\n", task.URL)
			continue
		}

		// Parse content
		parsedData, newURLs := p.Parse(content)

		// Save metadata
		p.kvStore.Save(task.URL, map[string]interface{}{
			"parsed":    true,
			"timestamp": utils.CurrentTimestamp(),
			"metadata":  parsedData,
		})

		// Produce new URLs to Kafka
		for _, url := range newURLs {
			if err := broker.Produce(url); err != nil {
				log.Printf("Failed to produce URL %s: %v", url, err)
			}
		}
	}
}

func (p *HTMLParser) AddTask(task models.ParseTask) {
	p.taskQueue <- task
}
