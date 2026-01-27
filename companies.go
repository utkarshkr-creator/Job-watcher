package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// CompanyCareer represents a company's career page configuration
type CompanyCareer struct {
	Name     string
	URL      string
	Selector string // CSS selector for job listings
	LinkAttr string // Attribute containing job link
}

// 60+ Big Tech Companies career pages - focused on India/Remote roles for freshers
var companyCareerPages = []CompanyCareer{
	// ========== Indian Unicorns & Startups ==========
	{Name: "Razorpay", URL: "https://razorpay.com/jobs/", Selector: "a[href*='/jobs/']", LinkAttr: "href"},
	{Name: "Zerodha", URL: "https://zerodha.com/careers/", Selector: "a[href*='careers']", LinkAttr: "href"},
	{Name: "PhonePe", URL: "https://www.phonepe.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Flipkart", URL: "https://www.flipkartcareers.com/#!/joblist", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Swiggy", URL: "https://careers.swiggy.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Zomato", URL: "https://www.zomato.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "CRED", URL: "https://careers.cred.club/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Meesho", URL: "https://careers.meesho.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Groww", URL: "https://groww.in/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Paytm", URL: "https://jobs.lever.co/paytm", Selector: "a.posting-title", LinkAttr: "href"},
	{Name: "Ola", URL: "https://www.olacabs.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Dunzo", URL: "https://www.dunzo.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Freshworks", URL: "https://www.freshworks.com/company/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Zoho", URL: "https://careers.zohocorp.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "InMobi", URL: "https://www.inmobi.com/company/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Postman", URL: "https://www.postman.com/company/careers/open-positions/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Dream11", URL: "https://www.dreamsports.group/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Udaan", URL: "https://careers.udaan.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Byju's", URL: "https://byjus.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Unacademy", URL: "https://unacademy.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "upGrad", URL: "https://www.upgrad.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Lenskart", URL: "https://www.lenskart.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Nykaa", URL: "https://careers.nykaa.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Cars24", URL: "https://www.cars24.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Zetwerk", URL: "https://www.zetwerk.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Vedantu", URL: "https://www.vedantu.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "ShareChat", URL: "https://sharechat.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Dailyhunt", URL: "https://www.dailyhunt.in/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Spinny", URL: "https://www.spinny.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Slice", URL: "https://www.sliceit.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Jupiter", URL: "https://jupiter.money/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Chargebee", URL: "https://www.chargebee.com/company/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "BrowserStack", URL: "https://www.browserstack.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Druva", URL: "https://www.druva.com/company/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "CleverTap", URL: "https://www.clevertap.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "MoEngage", URL: "https://www.moengage.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Hasura", URL: "https://hasura.io/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Polygon", URL: "https://polygon.technology/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "CoinDCX", URL: "https://coindcx.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "CoinSwitch", URL: "https://coinswitch.co/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Rapido", URL: "https://rapido.bike/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Urban Company", URL: "https://www.urbancompany.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Khatabook", URL: "https://khatabook.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "OkCredit", URL: "https://www.okcredit.in/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Acko", URL: "https://www.acko.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Digit Insurance", URL: "https://www.godigit.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "PolicyBazaar", URL: "https://www.policybazaar.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},

	// ========== Global Tech Companies with India offices ==========
	{Name: "Stripe", URL: "https://stripe.com/jobs/search?office_locations=Asia+Pacific--Bengaluru", Selector: "a[href*='/jobs/']", LinkAttr: "href"},
	{Name: "Notion", URL: "https://www.notion.so/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Figma", URL: "https://www.figma.com/careers/#job-openings", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Vercel", URL: "https://vercel.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Supabase", URL: "https://supabase.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "PlanetScale", URL: "https://planetscale.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Railway", URL: "https://railway.app/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Cloudflare", URL: "https://www.cloudflare.com/careers/jobs/?location=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Twilio", URL: "https://www.twilio.com/company/jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "MongoDB", URL: "https://www.mongodb.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Elastic", URL: "https://www.elastic.co/about/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "HashiCorp", URL: "https://www.hashicorp.com/careers/open-positions", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "GitLab", URL: "https://about.gitlab.com/jobs/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "GitHub", URL: "https://github.com/about/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Confluent", URL: "https://www.confluent.io/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Datadog", URL: "https://careers.datadoghq.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Snowflake", URL: "https://careers.snowflake.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Databricks", URL: "https://www.databricks.com/company/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Coinbase", URL: "https://www.coinbase.com/careers/positions", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Shopify", URL: "https://www.shopify.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Hubspot", URL: "https://www.hubspot.com/careers/jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Canva", URL: "https://www.canva.com/careers/jobs/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Airtable", URL: "https://airtable.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Asana", URL: "https://asana.com/jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Slack", URL: "https://slack.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Zoom", URL: "https://careers.zoom.us/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Dropbox", URL: "https://www.dropbox.com/jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Palantir", URL: "https://www.palantir.com/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Uber", URL: "https://www.uber.com/in/en/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Lyft", URL: "https://www.lyft.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Airbnb", URL: "https://careers.airbnb.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Spotify", URL: "https://www.lifeatspotify.com/jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Netflix", URL: "https://jobs.netflix.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Twitter/X", URL: "https://careers.twitter.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "LinkedIn", URL: "https://careers.linkedin.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Snap", URL: "https://careers.snap.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Pinterest", URL: "https://www.pinterestcareers.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Reddit", URL: "https://www.redditinc.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Discord", URL: "https://discord.com/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Roblox", URL: "https://careers.roblox.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Epic Games", URL: "https://www.epicgames.com/site/en-US/careers", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Unity", URL: "https://careers.unity.com/", Selector: "a[href*='job']", LinkAttr: "href"},

	// ========== Big Tech India ==========
	{Name: "Google India", URL: "https://careers.google.com/jobs/results/?location=India&q=software%20engineer", Selector: "a[href*='jobs']", LinkAttr: "href"},
	{Name: "Microsoft India", URL: "https://careers.microsoft.com/us/en/search-results?keywords=software%20engineer&location=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Amazon India", URL: "https://www.amazon.jobs/en/locations/india", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Meta India", URL: "https://www.metacareers.com/jobs?offices[0]=Bengaluru%2C%20India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Apple India", URL: "https://jobs.apple.com/en-in/search?location=india", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Adobe India", URL: "https://www.adobe.com/careers/india.html", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Oracle India", URL: "https://www.oracle.com/in/corporate/careers/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "SAP India", URL: "https://jobs.sap.com/search/?q=&locationsearch=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "IBM India", URL: "https://www.ibm.com/in-en/employment/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Intel India", URL: "https://jobs.intel.com/en/search-jobs/India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Nvidia India", URL: "https://nvidia.wd5.myworkdayjobs.com/NVIDIAExternalCareerSite?locationCountry=c4f78be1a8f14da0ab49ce1162348a5e", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Qualcomm India", URL: "https://careers.qualcomm.com/careers?location=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "VMware India", URL: "https://careers.vmware.com/location/india-jobs", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Cisco India", URL: "https://jobs.cisco.com/jobs/SearchJobs/India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "PayPal India", URL: "https://careers.pypl.com/home/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Atlassian", URL: "https://www.atlassian.com/company/careers/all-jobs?location=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Salesforce India", URL: "https://careers.salesforce.com/en/jobs/?country=India", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "ServiceNow", URL: "https://careers.servicenow.com/", Selector: "a[href*='job']", LinkAttr: "href"},
	{Name: "Workday", URL: "https://www.workday.com/en-us/company/careers/open-positions.html", Selector: "a[href*='job']", LinkAttr: "href"},
}

// Experience keywords to filter out (requires more than 2 years)
var seniorKeywords = []string{
	// Seniority titles
	"senior", "sr.", "sr ", "lead", "principal", "staff", "manager", "director",
	"head of", "vp ", "vice president", "architect",
	// Experience requirements > 2 years
	"3+", "4+", "5+", "6+", "7+", "8+", "10+",
	"3-5", "4-6", "5-7", "5-8", "6-8", "7-10", "8-10",
	"3 years", "4 years", "5 years", "6 years", "7 years", "8 years", "10 years",
	"3+ years", "4+ years", "5+ years",
	"three years", "four years", "five years",
}

// Keywords that indicate entry-level / 0-2 year roles (your experience range)
var entryLevelKeywords = []string{
	"fresher", "entry", "junior", "jr.", "jr ", "graduate", "new grad",
	"0-1", "0-2", "1-2", "1-3", "0-3", "2-3",
	"trainee", "associate", "intern", "campus",
	"entry level", "early career", "recent graduate",
}

// isEntryLevelJob checks if job is suitable for 1 year experience
func isEntryLevelJob(title string) bool {
	lower := strings.ToLower(title)

	// Check for senior keywords (exclude these - requires 3+ years)
	for _, kw := range seniorKeywords {
		if strings.Contains(lower, kw) {
			return false
		}
	}

	return true
}

// hasEntryLevelIndicator checks if job explicitly mentions entry-level
func hasEntryLevelIndicator(title string) bool {
	lower := strings.ToLower(title)
	for _, kw := range entryLevelKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// fetchCompanyJobs fetches jobs from a single company career page
func fetchCompanyJobs(company CompanyCareer) ([]Job, error) {
	req, err := http.NewRequest("GET", company.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	client := &http.Client{}
	resp, err := client.Do(req)
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
	seen := make(map[string]bool)
	jobIDRegex := regexp.MustCompile(`[\w-]+$`)

	doc.Find(company.Selector).Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr(company.LinkAttr)
		if !exists || link == "" {
			return
		}

		// Get title
		title := strings.TrimSpace(s.Text())
		if title == "" {
			title = s.Find("h2, h3, h4, span").First().Text()
			title = strings.TrimSpace(title)
		}

		if title == "" || len(title) > 150 {
			return
		}

		// Skip if already seen
		if seen[link] {
			return
		}
		seen[link] = true

		// EXPERIENCE FILTER: Skip senior/experienced roles
		if !isEntryLevelJob(title) {
			return
		}

		// Make absolute URL
		if !strings.HasPrefix(link, "http") {
			baseURL := company.URL
			if idx := strings.Index(baseURL, "//"); idx > 0 {
				if endIdx := strings.Index(baseURL[idx+2:], "/"); endIdx > 0 {
					baseURL = baseURL[:idx+2+endIdx]
				}
			}
			link = baseURL + link
		}

		// Extract job ID or use hash of URL for stability
		jobID := jobIDRegex.FindString(link)
		if jobID == "" {
			hash := sha256.Sum256([]byte(link))
			jobID = hex.EncodeToString(hash[:])[:12] // Use first 12 chars of hash
		}

		jobs = append(jobs, Job{
			ID:     fmt.Sprintf("%s-%s", strings.ToLower(strings.ReplaceAll(company.Name, " ", "")), jobID),
			Title:  fmt.Sprintf("%s @ %s", title, company.Name),
			Link:   link,
			Source: company.Name,
		})
	})

	return jobs, nil
}

// fetchAllCompanyJobsParallel fetches from all company career pages in parallel
func fetchAllCompanyJobsParallel() ([]Job, error) {
	var allJobs []Job
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrency to avoid rate limiting
	semaphore := make(chan struct{}, 10) // Max 10 concurrent requests

	for _, company := range companyCareerPages {
		wg.Add(1)
		go func(c CompanyCareer) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			jobs, err := fetchCompanyJobs(c)
			if err != nil {
				return
			}

			if len(jobs) > 0 {
				mu.Lock()
				allJobs = append(allJobs, jobs...)
				fmt.Printf("    %s: %d jobs\n", c.Name, len(jobs))
				mu.Unlock()
			}
		}(company)
	}

	wg.Wait()
	return allJobs, nil
}
