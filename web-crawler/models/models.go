package models

type ParseTask struct {
	URL     string
	Content string
}

func NewParseTask(url string, content string) ParseTask {
	return ParseTask{
		URL:     url,
		Content: content,
	}
}
