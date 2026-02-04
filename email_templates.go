package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GenerateColdEmail uses AI to draft a personalized email
func GenerateColdEmail(recruiter string, job Job, resume string, cfg AIConfig) string {
	if !cfg.Enabled || resume == "" {
		return ""
	}

	prompt := fmt.Sprintf(`Role: You are an expert copywriter for tech job applications.
Task: Write a short, high-impact cold email (max 100 words) to a recruiter.

Context:
- Recruiter Name: %s
- Job Title: %s
- Company: %s
- Candidate Experience:
%s

Instructions:
1. Subject line should be punchy (e.g., "Full-stack Engineer: [My Key Achievement]")
2. Body: Hook them immediately with a specific project/skill from my resume that matches their company/role.
3. Call to Action: "Available for a quick chat?"
4. Tone: Professional, confident, concise. No fluff.

Output Format:
Subject: [Subject]

[Body]`, recruiter, job.Title, job.Source, resume)

	reqBody := OllamaRequest{
		Model:  cfg.Model,
		Prompt: prompt,
		Stream: false,
		Format: "", // plain text
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error generating draft: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Sprintf("AI Error: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return ""
	}

	draft := strings.TrimSpace(ollamaResp.Response)
	return draft
}
