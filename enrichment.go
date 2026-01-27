package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// EnrichedJob holds extra data found for a job
type EnrichedJob struct {
	Recruiters    []RecruiterInfo
	CompanyInfo   string
	EnrichmentErr error
}

type RecruiterInfo struct {
	Name    string
	Profile string // URL
	Title   string // e.g. "Technical Recruiter"
}

// EnrichJob attempts to find recruiters and company info
// This is only called for high-scoring jobs to avoid rate limits
func EnrichJob(job Job) EnrichedJob {
	var ej EnrichedJob

	company := extractCompany(job)
	if company == "" {
		ej.EnrichmentErr = fmt.Errorf("could not extract company name")
		return ej
	}

	// 1. Find Recruiters (Fallback Strategy: DuckDuckGo/Google Search)
	// Query: site:linkedin.com/in "Technical Recruiter" "Company Name"
	recruiters, err := searchRecruiters(company)
	if err == nil {
		ej.Recruiters = recruiters
	} else {
		// Try a simpler search if first failed
		// fmt.Printf("Recruiter search failed: %v\n", err)
	}

	// 2. Company Info (Simple lookup)
	ej.CompanyInfo = fmt.Sprintf("Search: https://www.google.com/search?q=%s+funding+crunchbase", url.QueryEscape(company))

	return ej
}

func extractCompany(job Job) string {
	// Heuristic: "Role @ Company"
	if strings.Contains(job.Title, "@") {
		parts := strings.Split(job.Title, "@")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	// Fallback to Source if it's a company name (not "Indeed")
	if job.Source != "Indeed" && job.Source != "LinkedIn" && job.Source != "Naukri" {
		return job.Source
	}
	return ""
}

// searchRecruiters performs a search to find potential recruiters
// We use DuckDuckGo HTML which is easier to scrape than Google
func searchRecruiters(company string) ([]RecruiterInfo, error) {
	// Search query: site:linkedin.com/in "technical recruiter" company "India"
	query := fmt.Sprintf(`site:linkedin.com/in "technical recruiter" %s "India"`, company)
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var recruiters []RecruiterInfo
	doc.Find(".result__a").Each(func(i int, s *goquery.Selection) {
		if i >= 3 {
			return
		} // Limit to top 3

		title := s.Text()
		link, _ := s.Attr("href")

		// Parse name from Title usually "Name - Title | LinkedIn"
		parts := strings.Split(title, "-")
		name := strings.TrimSpace(parts[0])

		// Clean name
		if strings.Contains(name, "|") {
			name = strings.Split(name, "|")[0]
		}

		if link != "" && name != "" {
			recruiters = append(recruiters, RecruiterInfo{
				Name:    name,
				Profile: link,
				Title:   "Recruiter",
			})
		}
	})

	return recruiters, nil
}
