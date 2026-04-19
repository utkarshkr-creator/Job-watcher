// Test script to check individual job sources
// Run with: go run . --test-sources

package main

import (
	"fmt"
	"os"
)

func testSources() {
	fmt.Println("🧪 Testing Job Sources")
	fmt.Println("========================================================")

	// Load config
	cfg = loadConfig()
	if len(cfg.Keywords) > 0 {
		initKeywords(cfg)
	}

	// Test each source individually
	sources := map[string]func() ([]Job, error){
		"RemoteOK":     fetchJobs,
		"Razorpay":     fetchRazorpayJobs,
		"Wellfound":    fetchWellfoundJobs,
		"Indeed":       func() ([]Job, error) { return fetchIndeedJobs(cfg.IndeedRSS) },
		"LinkedIn":     fetchLinkedInJobs,
		"Naukri":       fetchNaukriJobs,
		"Instahyre":    fetchInstahyreJobs,
		"YC Jobs":      fetchYCJobs,
		"HN Jobs":      fetchHNJobs,
		"Reddit":       fetchRedditJobs,
		"Triplebyte":   fetchTriplebyteJobs,
		"Shared Lists": fetchSharedListJobs,
	}

	results := make(map[string]int)
	errors := make(map[string]error)

	for name, fetchFunc := range sources {
		// Check if source is enabled
		sourceKey := getSourceKey(name)
		if sourceKey != "" && !cfg.Sources[sourceKey] {
			fmt.Printf("⏭️  %s: DISABLED in config\n", name)
			continue
		}

		fmt.Printf("🔍 Testing %s...\n", name)
		jobs, err := fetchFunc()
		
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			errors[name] = err
			results[name] = 0
		} else {
			fmt.Printf("   ✅ Found %d jobs\n", len(jobs))
			results[name] = len(jobs)
			
			// Show first 3 jobs as sample
			if len(jobs) > 0 {
				fmt.Println("   Sample jobs:")
				for i, job := range jobs {
					if i >= 3 {
						break
					}
					fmt.Printf("     - %s\n", job.Title)
				}
			}
		}
		fmt.Println()
	}
	
	fmt.Println("\n========================================================")
	fmt.Println("📊 Summary")
	
	totalJobs := 0
	workingSources := 0
	failedSources := 0
	
	for name, count := range results {
		if count > 0 {
			workingSources++
			totalJobs += count
			fmt.Printf("✅ %s: %d jobs\n", name, count)
		} else if _, hasError := errors[name]; hasError {
			failedSources++
			fmt.Printf("❌ %s: FAILED\n", name)
		} else {
			fmt.Printf("⚠️  %s: 0 jobs (may need JavaScript or be rate-limited)\n", name)
		}
	}
	
	fmt.Printf("\nTotal: %d jobs from %d sources\n", totalJobs, workingSources)
	if failedSources > 0 {
		fmt.Printf("Failed: %d sources\n", failedSources)
	}
	
	// Exit with error if all sources failed
	if workingSources == 0 {
		fmt.Println("\n⚠️  WARNING: No sources returned jobs!")
		fmt.Println("This could indicate:")
		fmt.Println("  - Network connectivity issues")
		fmt.Println("  - API changes on job sites")
		fmt.Println("  - Rate limiting / IP blocking")
		fmt.Println("  - Sites requiring JavaScript")
		os.Exit(1)
	}
}

func getSourceKey(name string) string {
	mapping := map[string]string{
		"RemoteOK":     "remoteok",
		"Razorpay":     "razorpay",
		"Wellfound":    "wellfound",
		"Indeed":       "indeed",
		"LinkedIn":     "linkedin",
		"Naukri":       "naukri",
		"Instahyre":    "instahyre",
		"YC Jobs":      "ycjobs",
		"HN Jobs":      "hnjobs",
		"Reddit":       "reddit",
		"Triplebyte":   "triplebyte",
		"Shared Lists": "sharedlists",
	}
	return mapping[name]
}
