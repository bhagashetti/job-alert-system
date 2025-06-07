package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Job struct {
	Title   string
	Company string
	URL     string
}

func ScrapeJobs(keyword string) ([]Job, error) {
	resp, err := http.Get("https://remoteok.com/api")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	var jobs []Job
	for _, item := range data[1:] { // Skip the first element (metadata)
		if title, ok := item["position"].(string); ok &&
			item["url"] != nil &&
			item["company"] != nil {

			if keyword == "" || containsIgnoreCase(title, keyword) {
				jobs = append(jobs, Job{
					Title:   title,
					Company: item["company"].(string),
					URL:     item["url"].(string),
				})
			}
		}
	}

	return jobs, nil
}

func containsIgnoreCase(text, substr string) bool {
	return len(substr) == 0 || (len(text) > 0 && (len(substr) > 0 &&
		stringContainsCI(text, substr)))
}

func stringContainsCI(a, b string) bool {
	return len(a) >= len(b) && (a == b ||
		len(a) > len(b) && (containsIgnoreCase(a[1:], b) || containsIgnoreCase(a, b[1:])))
}
