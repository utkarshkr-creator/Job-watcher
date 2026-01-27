package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ================== Y COMBINATOR JOBS ==================
// workatastartup.com - YC company job board

type YCJobResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	Remote      bool   `json:"remote"`
	URL         string `json:"url"`
}

func fetchYCJobs() ([]Job, error) {
	// YC Work at a Startup has a public API
	url := "https://www.workatastartup.com/api/jobs?query=software+engineer&page=1"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Fallback to scraping the HTML page
		return fetchYCJobsHTML()
	}

	var data struct {
		Jobs []YCJobResponse `json:"jobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fetchYCJobsHTML()
	}

	var jobs []Job
	for _, j := range data.Jobs {
		title := j.Title
		if j.CompanyName != "" {
			title = fmt.Sprintf("%s @ %s", j.Title, j.CompanyName)
		}
		if j.Remote {
			title += " (Remote)"
		} else if j.Location != "" {
			title += fmt.Sprintf(" (%s)", j.Location)
		}

		// Apply experience filter
		if !isEntryLevelJob(title) {
			continue
		}

		link := j.URL
		if link == "" {
			link = fmt.Sprintf("https://www.workatastartup.com/jobs/%d", j.ID)
		}

		jobs = append(jobs, Job{
			ID:     fmt.Sprintf("yc-%d", j.ID),
			Title:  title,
			Link:   link,
			Source: "YC Jobs",
		})
	}

	return jobs, nil
}

func fetchYCJobsHTML() ([]Job, error) {
	url := "https://www.workatastartup.com/jobs"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	var jobs []Job
	// Parse job links from HTML
	linkRegex := regexp.MustCompile(`href="(/jobs/\d+)"`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)

	seen := make(map[string]bool)
	for i, match := range matches {
		if len(match) < 2 || i >= 30 {
			break
		}
		path := match[1]
		if seen[path] {
			continue
		}
		seen[path] = true

		jobs = append(jobs, Job{
			ID:     "yc-" + strings.TrimPrefix(path, "/jobs/"),
			Title:  "Software Engineer @ YC Startup",
			Link:   "https://www.workatastartup.com" + path,
			Source: "YC Jobs",
		})
	}

	return jobs, nil
}

// ================== HACKER NEWS JOBS ==================
// Monthly "Who's Hiring?" thread

func fetchHNJobs() ([]Job, error) {
	// Get the latest "Who is hiring?" thread from HN
	// Search for the thread ID
	searchURL := "https://hn.algolia.com/api/v1/search?query=who%20is%20hiring&tags=ask_hn&hitsPerPage=5"

	req, _ := http.NewRequest("GET", searchURL, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult struct {
		Hits []struct {
			ObjectID  string `json:"objectID"`
			Title     string `json:"title"`
			CreatedAt string `json:"created_at"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	// Find the most recent "Who is hiring" thread
	var threadID string
	for _, hit := range searchResult.Hits {
		if strings.Contains(strings.ToLower(hit.Title), "who is hiring") {
			threadID = hit.ObjectID
			break
		}
	}

	if threadID == "" {
		return []Job{}, nil
	}

	// Get comments from the thread
	itemURL := fmt.Sprintf("https://hn.algolia.com/api/v1/items/%s", threadID)
	resp2, err := client.Get(itemURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var item struct {
		Children []struct {
			ID     int    `json:"id"`
			Text   string `json:"text"`
			Author string `json:"author"`
		} `json:"children"`
	}

	if err := json.NewDecoder(resp2.Body).Decode(&item); err != nil {
		return nil, err
	}

	var jobs []Job
	for i, child := range item.Children {
		if i >= 50 { // Limit to first 50 comments
			break
		}

		text := child.Text
		if len(text) < 50 {
			continue
		}

		// Extract company name from first line
		lines := strings.Split(text, "<p>")
		firstLine := lines[0]
		firstLine = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(firstLine, "")
		firstLine = strings.TrimSpace(firstLine)

		if len(firstLine) > 100 {
			firstLine = firstLine[:100] + "..."
		}

		// Skip if too senior
		if !isEntryLevelJob(text) {
			continue
		}

		jobs = append(jobs, Job{
			ID:     fmt.Sprintf("hn-%d", child.ID),
			Title:  firstLine,
			Link:   fmt.Sprintf("https://news.ycombinator.com/item?id=%d", child.ID),
			Source: "HN Jobs",
		})
	}

	return jobs, nil
}

// ================== REDDIT JOBS ==================
// r/cscareerquestions and r/forhire

func fetchRedditJobs() ([]Job, error) {
	subreddits := []string{
		"cscareerquestions",
		"forhire",
	}

	var allJobs []Job

	for _, sub := range subreddits {
		// Search for hiring posts
		url := fmt.Sprintf("https://www.reddit.com/r/%s/search.json?q=hiring+OR+job&sort=new&t=week&limit=25", sub)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "JobWatcher/1.0")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		var data struct {
			Data struct {
				Children []struct {
					Data struct {
						ID        string  `json:"id"`
						Title     string  `json:"title"`
						Permalink string  `json:"permalink"`
						Selftext  string  `json:"selftext"`
						Score     int     `json:"score"`
						Created   float64 `json:"created_utc"`
					} `json:"data"`
				} `json:"children"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for _, child := range data.Data.Children {
			post := child.Data
			title := post.Title

			// Filter for relevant posts
			lowerTitle := strings.ToLower(title)
			if !strings.Contains(lowerTitle, "hiring") &&
				!strings.Contains(lowerTitle, "job") &&
				!strings.Contains(lowerTitle, "[hiring]") {
				continue
			}

			// Skip if senior role
			if !isEntryLevelJob(title) && !isEntryLevelJob(post.Selftext) {
				continue
			}

			link := "https://www.reddit.com" + post.Permalink

			allJobs = append(allJobs, Job{
				ID:     "reddit-" + post.ID,
				Title:  title,
				Link:   link,
				Source: "Reddit",
			})
		}
	}

	return allJobs, nil
}

// ================== TRIPLEBYTE / KARAT ==================
// Note: These require login, so we scrape public listings

func fetchTriplebyteJobs() ([]Job, error) {
	// Triplebyte was acquired by Karat, limited public access
	url := "https://triplebyte.com/jobs"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	var jobs []Job
	linkRegex := regexp.MustCompile(`href="(/company/[^"]+/jobs/[^"]+)"`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)

	seen := make(map[string]bool)
	for i, match := range matches {
		if len(match) < 2 || i >= 20 {
			break
		}
		path := match[1]
		if seen[path] {
			continue
		}
		seen[path] = true

		jobs = append(jobs, Job{
			ID:     "triplebyte-" + fmt.Sprintf("%d", i),
			Title:  "Software Engineer",
			Link:   "https://triplebyte.com" + path,
			Source: "Triplebyte",
		})
	}

	return jobs, nil
}

// ================== HIRED.COM ==================
// Reverse job board - companies apply to you

func fetchHiredJobs() ([]Job, error) {
	// Hired.com requires signup, but we can scrape featured companies
	url := "https://hired.com/companies"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Hired is primarily a platform where YOU create a profile
	// Return empty as it's not traditional job scraping
	return []Job{}, nil
}
