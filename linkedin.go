package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LinkedIn search configuration
type LinkedInSearch struct {
	Keywords   string
	Location   string
	TimePosted string // r86400 = past 24h, r604800 = past week
}

var linkedinSearches = []LinkedInSearch{
	{Keywords: "software engineer", Location: "India", TimePosted: "r86400"},
	{Keywords: "backend developer", Location: "India", TimePosted: "r86400"},
	{Keywords: "full stack developer", Location: "India", TimePosted: "r86400"},
	{Keywords: "react developer", Location: "India", TimePosted: "r86400"},
	{Keywords: "node.js developer", Location: "India", TimePosted: "r86400"},
	{Keywords: "software engineer", Location: "", TimePosted: "r86400"}, // Remote
}

func fetchLinkedInJobs() ([]Job, error) {
	var allJobs []Job

	for _, search := range linkedinSearches {
		jobs, err := scrapeLinkedInSearch(search)
		if err != nil {
			fmt.Printf("  Warning: LinkedIn search failed for '%s': %v\n", search.Keywords, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	// Deduplicate
	seen := make(map[string]bool)
	var unique []Job
	for _, j := range allJobs {
		if !seen[j.ID] {
			seen[j.ID] = true
			unique = append(unique, j)
		}
	}

	return unique, nil
}

func scrapeLinkedInSearch(search LinkedInSearch) ([]Job, error) {
	// LinkedIn guest jobs API endpoint
	baseURL := "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

	// Build query parameters
	keywords := strings.ReplaceAll(search.Keywords, " ", "%20")
	location := strings.ReplaceAll(search.Location, " ", "%20")

	url := fmt.Sprintf("%s?keywords=%s&location=%s&f_TPR=%s&start=0",
		baseURL, keywords, location, search.TimePosted)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	var jobs []Job
	jobIDRegex := regexp.MustCompile(`/jobs/view/(\d+)`)

	// Parse job cards from LinkedIn HTML
	doc.Find("li, .base-card, .job-search-card").Each(func(i int, s *goquery.Selection) {
		// Find title
		title := s.Find(".base-search-card__title, h3, h4").First().Text()
		title = strings.TrimSpace(title)

		if title == "" {
			return
		}

		// Find link
		link, _ := s.Find("a.base-card__full-link, a[href*='/jobs/view/']").First().Attr("href")
		if link == "" {
			link, _ = s.Find("a").First().Attr("href")
		}

		if link == "" || !strings.Contains(link, "/jobs/") {
			return
		}

		// Extract job ID from URL
		var jobID string
		if matches := jobIDRegex.FindStringSubmatch(link); len(matches) > 1 {
			jobID = matches[1]
		} else {
			jobID = fmt.Sprintf("%d", i)
		}

		// Get company name
		company := s.Find(".base-search-card__subtitle, .job-search-card__company-name").First().Text()
		company = strings.TrimSpace(company)

		// Get location
		location := s.Find(".job-search-card__location, .base-search-card__metadata span").First().Text()
		location = strings.TrimSpace(location)

		// Combine title with company if available
		fullTitle := title
		if company != "" {
			fullTitle = fmt.Sprintf("%s @ %s", title, company)
		}
		if location != "" {
			fullTitle = fmt.Sprintf("%s (%s)", fullTitle, location)
		}

		jobs = append(jobs, Job{
			ID:     "linkedin-" + jobID,
			Title:  fullTitle,
			Link:   link,
			Source: "LinkedIn",
		})
	})

	return jobs, nil
}
