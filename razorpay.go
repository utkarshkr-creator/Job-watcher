package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchRazorpayJobs() ([]Job, error) {
	req, _ := http.NewRequest("GET", "https://razorpay.com/jobs/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	var jobs []Job

	// Try multiple selectors that Razorpay might use
	selectors := []string{
		"a[href*='/jobs/']",
		".job-card",
		".job-listing",
		"a[href*='lever.co']",
		"a[href*='greenhouse']",
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if !exists {
				return
			}

			// Skip non-job links
			if strings.Contains(link, "/jobs/") && !strings.HasSuffix(link, "/jobs/") {
				title := strings.TrimSpace(s.Text())
				if title == "" {
					title = s.Find("h3, h4, .title, .job-title").Text()
				}
				if title == "" {
					return
				}

				// Clean up title
				title = strings.TrimSpace(title)
				if len(title) > 100 {
					title = title[:100]
				}

				// Make absolute URL
				if !strings.HasPrefix(link, "http") {
					link = "https://razorpay.com" + link
				}

				jobs = append(jobs, Job{
					ID:     "razorpay-" + link,
					Title:  title,
					Link:   link,
					Source: "Razorpay",
				})
			}
		})
		if len(jobs) > 0 {
			break
		}
	}

	return jobs, nil
}
