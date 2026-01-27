package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Instahyre API response structure
type InstahyreResponse struct {
	Jobs []struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Company    string `json:"company_name"`
		Location   string `json:"location"`
		Experience string `json:"experience"`
		Slug       string `json:"slug"`
	} `json:"results"`
}

// fetchInstahyreJobs fetches from Instahyre API
func fetchInstahyreJobs() ([]Job, error) {
	var allJobs []Job

	// Instahyre has a public API for job search
	searches := []string{
		"software-engineer",
		"backend-developer",
		"full-stack-developer",
		"frontend-developer",
	}

	for _, search := range searches {
		url := fmt.Sprintf("https://www.instahyre.com/api/v1/candidate/opportunities/?job_type=%s&experience=0-1", search)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		var data InstahyreResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for _, j := range data.Jobs {
			title := j.Title
			if j.Company != "" {
				title = fmt.Sprintf("%s @ %s", j.Title, j.Company)
			}
			if j.Location != "" {
				title = fmt.Sprintf("%s (%s)", title, j.Location)
			}

			allJobs = append(allJobs, Job{
				ID:     fmt.Sprintf("instahyre-%d", j.ID),
				Title:  title,
				Link:   fmt.Sprintf("https://www.instahyre.com/job/%s/", j.Slug),
				Source: "Instahyre",
			})
		}
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

// fetchHiristJobs fetches from Hirist (another India-focused job board)
func fetchHiristJobs() ([]Job, error) {
	var allJobs []Job

	urls := []string{
		"https://www.hirist.tech/jobs/software-engineer-fresher",
		"https://www.hirist.tech/jobs/backend-developer-fresher",
	}

	for _, url := range urls {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
		req.Header.Set("Accept", "text/html")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
	}

	return allJobs, nil
}

// fetchStackOverflowJobs fetches from Stack Overflow Jobs (via Indeed redirect)
func fetchStackOverflowJobs() ([]Job, error) {
	// Stack Overflow jobs was discontinued, jobs now redirect to other sources
	return []Job{}, nil
}

// fetchCutshortJobs fetches from Cutshort
func fetchCutshortJobs() ([]Job, error) {
	url := "https://cutshort.io/jobs?experience=0-2"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Cutshort uses heavy JS, so HTTP scraping won't work well
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return []Job{}, nil
}

// Internshala for internships and fresher jobs
func fetchInternshalaJobs() ([]Job, error) {
	url := "https://internshala.com/jobs/software-developer-jobs/"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs []Job

	// Parse HTML for job listings
	// Internhsala loads content via JS, so this may be limited
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return jobs, nil
}

// SimplifyJobs - API for entry-level jobs
type SimplifyJobsResponse struct {
	Jobs []struct {
		ID          string   `json:"id"`
		Title       string   `json:"title"`
		CompanyName string   `json:"company_name"`
		Locations   []string `json:"locations"`
		URL         string   `json:"url"`
	} `json:"jobs"`
}

func fetchSimplifyJobs() ([]Job, error) {
	// Simplify.jobs has a public API
	url := "https://api.simplify.jobs/v1/jobs?featured=true"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var data SimplifyJobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var jobs []Job
	for _, j := range data.Jobs {
		title := j.Title
		if j.CompanyName != "" {
			title = fmt.Sprintf("%s @ %s", j.Title, j.CompanyName)
		}
		if len(j.Locations) > 0 {
			title = fmt.Sprintf("%s (%s)", title, strings.Join(j.Locations, ", "))
		}

		jobs = append(jobs, Job{
			ID:     "simplify-" + j.ID,
			Title:  title,
			Link:   j.URL,
			Source: "Simplify",
		})
	}

	return jobs, nil
}
