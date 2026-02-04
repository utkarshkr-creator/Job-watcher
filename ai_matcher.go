package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// AIConfig holds settings for Ollama
type AIConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Model     string `yaml:"model"`
	Threshold int    `yaml:"threshold"` // Score out of 100
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"` // json
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type MatchResult struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

var (
	resumeText string
	aiMutex    sync.Mutex
)

// loadResume loads resume.txt into memory
func loadResume() error {
	data, err := os.ReadFile("resume.txt")
	if err != nil {
		return err
	}
	resumeText = string(data)
	if len(resumeText) < 50 {
		return fmt.Errorf("resume.txt seems too short, please paste your resume content")
	}
	return nil
}

// scoreJobWithAI sends job title+resume to Ollama
func scoreJobWithAI(job Job, cfg AIConfig) (int, string, error) {
	// M3 Pro can handle parallel AI calls, so removed the lock
	// aiMutex.Lock()
	// defer aiMutex.Unlock()

	// Use Title and Source for context
	// Use Title and Source for context
	prompt := fmt.Sprintf(`Role: Hiring Manager. 
Task: Evaluate match for a Fresher/Entry-Level Candidate (0-2 YOE).

Candidate Resume:
%s

Job Title: %s
Company: %s

CRITICAL RULES:
1. IF Job Title contains "Senior", "Staff", "Lead", "Principal", "Architect", or requires >2 years experience: SCORE MUST BE 0.
2. IF Job is for "Intern", "New Grad", "Associate", "Junior", or "0-2 years": Score normally based on skill match.
3. IGNORE skill match if Rule #1 is violated.

Constraint: Return ONLY a JSON object with "score" (0-100) and "reason" (short string).
Example: {"score": 0, "reason": "Senior role (3+ years) not for fresher"}`, resumeText, job.Title, job.Source)

	reqBody := OllamaRequest{
		Model:  cfg.Model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, "", fmt.Errorf("ollama connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return 0, "", err
	}

	// Parse inner JSON from LLM
	var result MatchResult
	// Sometimes LLMs add markdown code blocks like ```json ... ```
	cleanResp := strings.TrimSpace(ollamaResp.Response)
	cleanResp = strings.TrimPrefix(cleanResp, "```json")
	cleanResp = strings.TrimSuffix(cleanResp, "```")

	if err := json.Unmarshal([]byte(cleanResp), &result); err != nil {
		// Fallback: simple text parsing if JSON fails
		return 0, cleanResp, nil
	}

	return result.Score, result.Reason, nil
}
