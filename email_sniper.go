package main

import (
	"fmt"
	"strings"
)

// EmailSniper handles guessing emails
type EmailSniper struct{}

func NewEmailSniper() *EmailSniper {
	return &EmailSniper{}
}

// GuessEmails generates likely email permutations for a recruiter
func (es *EmailSniper) GuessEmails(name, company string) []string {
	domain := es.guessDomain(company)
	if domain == "" {
		return []string{}
	}

	parts := strings.Fields(strings.ToLower(strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r == ' ' {
			return r
		}
		return -1
	}, name)))

	if len(parts) < 2 {
		return []string{} // need first and last name
	}

	first := parts[0]
	last := parts[len(parts)-1]

	f := string(first[0]) // first initial

	// Common corporate email patterns
	patterns := []string{
		fmt.Sprintf("%s.%s@%s", first, last, domain), // first.last@company.com
		fmt.Sprintf("%s%s@%s", first, last, domain),  // firstlast@company.com
		fmt.Sprintf("%s@%s", first, domain),          // first@company.com
		fmt.Sprintf("%s%s@%s", f, last, domain),      // flast@company.com (very common)
	}

	return patterns
}

// guessDomain tries to infer the domain from company name
func (es *EmailSniper) guessDomain(company string) string {
	// 1. Clean the company name
	clean := strings.ToLower(company)
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, "inc", "")
	clean = strings.ReplaceAll(clean, "ltd", "")
	clean = strings.ReplaceAll(clean, "pvt", "")
	clean = strings.TrimSpace(clean)

	// 2. Known overrides (optional)
	// if clean == "alphagroup" { return "alpha.com" }

	// 3. Simple heuristic: company.com
	// In production, you'd use a Clearbit/Hunter API or Google "company email format"
	return fmt.Sprintf("%s.com", clean)
}
