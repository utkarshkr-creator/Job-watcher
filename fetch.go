package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Job struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source,omitempty"`
}

func fetchJobs() ([]Job, error) {
	req, _ := http.NewRequest("GET", "https://remoteok.com/api", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var jobs []Job
	for _, j := range raw {
		if j["id"] == nil || j["position"] == nil {
			continue
		}

		// Handle ID which can be string or number
		var id string
		switch v := j["id"].(type) {
		case float64:
			id = fmt.Sprintf("%.0f", v)
		case string:
			id = v
		default:
			continue
		}

		// Handle URL - API returns full URLs or relative paths
		urlPath, ok := j["url"].(string)
		if !ok {
			continue
		}

		link := urlPath
		if len(urlPath) > 0 && urlPath[0] == '/' {
			link = "https://remoteok.com" + urlPath
		}

		jobs = append(jobs, Job{
			ID:     "remoteok-" + id,
			Title:  j["position"].(string),
			Link:   link,
			Source: "RemoteOK",
		})
	}

	return jobs, nil
}
