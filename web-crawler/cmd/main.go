package main

import (
	"log"
	"web-crawler/broker"
	"web-crawler/crawler"
	"web-crawler/parser"
	"web-crawler/storage"
	"web-crawler/utils"
)

func main() {
	config := utils.LoadConfig()

	kafkaBroker := broker.NewKafkaBroker(config.Kafka.Brokers, config.Kafka.Topic)
	//robotsStorage := storage.NewRobotsStorage()
	robotsStorage := storage.NewFakeRobotStorage()

	kvStore := storage.NewKVStore(config.KVStore.Path)
	//fileStore := storage.NewFileStore(config.Storage.Path)
	parser := parser.NewHTMLParser(kvStore)
	go func() {
		parser.StartParsingWorker(kafkaBroker)
	}()

	crawler := crawler.NewCrawler(kafkaBroker, robotsStorage, parser)

	// Start crawling
	if err := crawler.Start(); err != nil {
		log.Fatalf("Failed to start crawler: %v", err)
	}
}
