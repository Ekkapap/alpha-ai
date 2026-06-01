package main

// semantic.go — generate .graphify_analysis.json community labels.
// Priority: project-root/.env API key → agent fallback (return prompt to MCP caller).

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// APIKey holds the first found API key and its provider.
type APIKey struct {
	Provider string // "anthropic", "gemini", "openai"
	Key      string
}

// LoadAPIKey reads project-root/.env and returns the first found LLM key.
func LoadAPIKey(projectRoot string) *APIKey {
	envPath := filepath.Join(projectRoot, ".env")
	f, err := os.Open(envPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	vars := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		vars[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), `"'`)
	}

	for _, entry := range []struct{ env, provider string }{
		{"ANTHROPIC_API_KEY", "anthropic"},
		{"GEMINI_API_KEY", "gemini"},
		{"GOOGLE_API_KEY", "gemini"},
		{"OPENAI_API_KEY", "openai"},
	} {
		if v := vars[entry.env]; v != "" {
			return &APIKey{Provider: entry.provider, Key: v}
		}
	}
	return nil
}

// GenerateCommunityLabels calls the LLM API to name communities.
// Returns map[communityID]name.
func GenerateCommunityLabels(key *APIKey, communities map[string][]string, nodeLabels map[string]string) (map[string]string, error) {
	// Build a compact prompt
	var sb strings.Builder
	sb.WriteString("Name each community with 2-4 words describing what these code symbols do together.\n")
	sb.WriteString("Respond with JSON only: {\"0\": \"name\", ...}\n\n")
	for cid, members := range communities {
		// sample up to 8 members
		sample := members
		if len(sample) > 8 {
			sample = sample[:8]
		}
		labels := make([]string, 0, len(sample))
		for _, id := range sample {
			if lbl, ok := nodeLabels[id]; ok {
				labels = append(labels, lbl)
			} else {
				labels = append(labels, id)
			}
		}
		sb.WriteString(fmt.Sprintf("Community %s: %s\n", cid, strings.Join(labels, ", ")))
	}
	prompt := sb.String()

	switch key.Provider {
	case "anthropic":
		return callAnthropic(key.Key, prompt)
	case "gemini":
		return callGemini(key.Key, prompt)
	case "openai":
		return callOpenAI(key.Key, prompt)
	}
	return nil, fmt.Errorf("unknown provider: %s", key.Provider)
}

// AgentFallbackPrompt returns a prompt string for the calling agent to generate labels.
func AgentFallbackPrompt(communities map[string][]string, nodeLabels map[string]string) string {
	var sb strings.Builder
	sb.WriteString("No API key found. Please call set_labels with a JSON mapping of community IDs to names.\n")
	sb.WriteString("Example call: set_labels {\"0\": \"Auth Layer\", \"1\": \"Graph Core\"}\n\n")
	sb.WriteString("Communities:\n")
	for cid, members := range communities {
		sample := members
		if len(sample) > 6 {
			sample = sample[:6]
		}
		labels := make([]string, 0, len(sample))
		for _, id := range sample {
			if lbl, ok := nodeLabels[id]; ok {
				labels = append(labels, lbl)
			} else {
				labels = append(labels, id)
			}
		}
		sb.WriteString(fmt.Sprintf("  %s: %s\n", cid, strings.Join(labels, ", ")))
	}
	return sb.String()
}

// ── API callers ───────────────────────────────────────────────────────────────

func callAnthropic(key, prompt string) (map[string]string, error) {
	body := map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 512,
		"messages":   []map[string]any{{"role": "user", "content": prompt}},
	}
	return postJSON("https://api.anthropic.com/v1/messages", map[string]string{
		"x-api-key":         key,
		"anthropic-version": "2023-06-01",
		"content-type":      "application/json",
	}, body, func(data []byte) (map[string]string, error) {
		var resp struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(data, &resp); err != nil || len(resp.Content) == 0 {
			return nil, fmt.Errorf("parse anthropic response")
		}
		return parseJSONLabels(resp.Content[0].Text)
	})
}

func callGemini(key, prompt string) (map[string]string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + key
	body := map[string]any{
		"contents": []map[string]any{{"parts": []map[string]any{{"text": prompt}}}},
	}
	return postJSON(url, map[string]string{"content-type": "application/json"}, body, func(data []byte) (map[string]string, error) {
		var resp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}
		if err := json.Unmarshal(data, &resp); err != nil || len(resp.Candidates) == 0 {
			return nil, fmt.Errorf("parse gemini response")
		}
		if len(resp.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("empty gemini response")
		}
		return parseJSONLabels(resp.Candidates[0].Content.Parts[0].Text)
	})
}

func callOpenAI(key, prompt string) (map[string]string, error) {
	body := map[string]any{
		"model":      "gpt-4o-mini",
		"max_tokens": 512,
		"messages":   []map[string]any{{"role": "user", "content": prompt}},
	}
	return postJSON("https://api.openai.com/v1/chat/completions", map[string]string{
		"Authorization": "Bearer " + key,
		"content-type":  "application/json",
	}, body, func(data []byte) (map[string]string, error) {
		var resp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(data, &resp); err != nil || len(resp.Choices) == 0 {
			return nil, fmt.Errorf("parse openai response")
		}
		return parseJSONLabels(resp.Choices[0].Message.Content)
	})
}

func postJSON(url string, headers map[string]string, body any, parse func([]byte) (map[string]string, error)) (map[string]string, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API %d: %s", resp.StatusCode, string(data))
	}
	return parse(data)
}

// parseJSONLabels extracts {"0": "name"} from LLM response text (may have markdown fences).
func parseJSONLabels(text string) (map[string]string, error) {
	// Strip markdown fences
	text = strings.TrimSpace(text)
	if i := strings.Index(text, "{"); i >= 0 {
		text = text[i:]
	}
	if i := strings.LastIndex(text, "}"); i >= 0 {
		text = text[:i+1]
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse labels JSON: %w", err)
	}
	return result, nil
}
