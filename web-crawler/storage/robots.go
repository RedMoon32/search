package storage

import (
	"net/http"
	"sync"

	"github.com/temoto/robotstxt"
)

type RobotsTxtCache struct {
	cache map[string]*robotstxt.Group
	mu    sync.RWMutex
}

func NewRobotsStorage() *RobotsTxtCache {
	return &RobotsTxtCache{
		cache: make(map[string]*robotstxt.Group),
	}
}

func (r *RobotsTxtCache) Allowed(url string) bool {
	r.mu.RLock()
	group, exists := r.cache[url]
	r.mu.RUnlock()

	if !exists {
		resp, err := http.Get(url + "/robots.txt")
		if err != nil || resp.StatusCode != 200 {
			return true
		}

		robotsData, err := robotstxt.FromResponse(resp)
		if err != nil {
			return true
		}

		group = robotsData.FindGroup("*")

		r.mu.Lock()
		r.cache[url] = group
		r.mu.Unlock()
	}

	return group.Test(url)
}
