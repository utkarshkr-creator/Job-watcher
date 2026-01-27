package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func generateStableHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])[:12]
}

// ListSource represents a public list (Google Sheet or GitHub README)
type ListSource struct {
	Name string
	URL  string
	Type string // "googlesheet" or "github"
}

// Known public hiring lists
var publicLists = []ListSource{
	{
		Name: "India New Grad Roles 2025 (GitHub)",
		URL:  "https://github.com/samiranghosh04/new-grad-tech-roles--india",
		Type: "github",
	},
	{
		Name: "Simplify.jobs GitHub List",
		URL:  "https://github.com/SimplifyJobs/New-Grad-Positions",
		Type: "github",
	},
	// Add your Google Sheets URLs here in config
}

func fetchSharedListJobs() ([]Job, error) {
	var allJobs []Job

	for _, source := range publicLists {
		var jobs []Job
		var err error

		if source.Type == "googlesheet" {
			jobs, err = scrapeGoogleSheet(source.URL)
		} else if source.Type == "github" {
			jobs, err = scrapeGitHubReadme(source.URL)
		}

		if err == nil {
			allJobs = append(allJobs, jobs...)
		} else {
			fmt.Printf("Error scraping %s: %v\n", source.Name, err)
		}
	}

	return allJobs, nil
}

// scrapeGoogleSheet downloads a public Google Sheet as CSV and identifies job rows
func scrapeGoogleSheet(sheetURL string) ([]Job, error) {
	// Convert /edit URL to /export?format=csv
	csvURL := sheetURL
	if strings.Contains(sheetURL, "/edit") {
		parts := strings.Split(sheetURL, "/edit")
		csvURL = parts[0] + "/export?format=csv"
	}

	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("empty sheet")
	}

	// Heuristic: Find columns for Company, Role, Link
	headers := rows[0]
	colMap := make(map[string]int)
	for i, h := range headers {
		lower := strings.ToLower(h)
		if strings.Contains(lower, "company") || strings.Contains(lower, "name") {
			colMap["company"] = i
		} else if strings.Contains(lower, "role") || strings.Contains(lower, "position") || strings.Contains(lower, "title") {
			colMap["role"] = i
		} else if strings.Contains(lower, "link") || strings.Contains(lower, "url") || strings.Contains(lower, "apply") {
			colMap["link"] = i
		} else if strings.Contains(lower, "location") {
			colMap["location"] = i
		}
	}

	// Validations
	if _, ok := colMap["company"]; !ok {
		return nil, fmt.Errorf("could not find Company column")
	}

	var jobs []Job
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		// Safer access
		getCol := func(key string) string {
			if idx, ok := colMap[key]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		company := getCol("company")
		role := getCol("role")
		link := getCol("link")
		location := getCol("location")

		if company == "" || (role == "" && link == "") {
			continue
		}

		// Defaults
		if role == "" {
			role = "Software Engineer"
		}

		title := fmt.Sprintf("%s @ %s", role, company)
		if location != "" {
			title += fmt.Sprintf(" (%s)", location)
		}

		// Experience filter
		if !isEntryLevelJob(title) {
			continue
		}

		// If no link, skip (unless we want just alerts)
		if link == "" {
			continue
		}

		// Normalize link
		if !strings.HasPrefix(link, "http") {
			continue
		}

		jobs = append(jobs, Job{
			ID:     fmt.Sprintf("sheet-%s", generateStableHash(company+role+link)),
			Title:  title,
			Link:   link,
			Source: "Shared List",
		})
	}

	return jobs, nil
}

// scrapeGitHubReadme looks for markdown tables in GitHub READMEs
func scrapeGitHubReadme(repoURL string) ([]Job, error) {
	resp, err := http.Get(repoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs []Job

	// Find all tables and iterate rows
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		// skip header often
		if s.Find("th").Length() > 0 {
			return
		}

		cols := s.Find("td")
		if cols.Length() < 3 {
			return
		}

		// Heuristic: usually Company | Role | Location | Application/Link
		// But layouts vary. We'll look for the first link.

		company := strings.TrimSpace(cols.Eq(0).Text())

		// Find link in any column
		var link string
		var role string

		s.Find("a").Each(func(j int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			text := strings.TrimSpace(a.Text())

			if exists && strings.HasPrefix(href, "http") &&
				!strings.Contains(href, "github.com") && // skip internal links often
				!strings.Contains(href, "linkedin.com/company") { // skip linkedin company pages
				link = href
				if role == "" {
					role = text
				}
			}
		})

		if link == "" {
			// Check if the role column has text but no link
			// Sometimes link is "Apply" button in last column
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				href, exists := a.Attr("href")
				if exists && strings.HasPrefix(href, "http") {
					link = href
				}
			})
		}

		if company == "" || link == "" {
			return
		}

		if role == "" || strings.ToLower(role) == "apply" {
			// Try to find text in 2nd column
			role = strings.TrimSpace(cols.Eq(1).Text())
		}

		if role == "" {
			role = "Software Engineer"
		}

		title := fmt.Sprintf("%s @ %s", role, company)

		// Filter
		if !isEntryLevelJob(title) {
			return
		}

		jobs = append(jobs, Job{
			ID:     fmt.Sprintf("github-list-%s", generateStableHash(company+role+link)),
			Title:  title,
			Link:   link,
			Source: "GitHub List",
		})
	})

	return jobs, nil
}
