package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

var allowedHosts = map[string]bool{
	"boards.greenhouse.io":     true,
	"jobs.ashbyhq.com":         true,
	"jobs.lever.co":            true,
	"myworkdayjobs.com":        true,
	"jobs.smartrecruiters.com": true,
	"apply.workable.com":       true,
	"recruiting.adp.com":       true,
	"recruiting2.ultipro.com":  true,
	"jobs.jobvite.com":         true,
}
var allowedSuffixes = []string{".icims.com", ".bamboohr.com", ".recruitee.com", ".breezy.hr"}

func hostAllowed(u *url.URL) bool {
	if allowedHosts[u.Hostname()] {
		return true
	}
	h := u.Hostname()
	for _, suf := range allowedSuffixes {
		if len(h) >= len(suf) && h[len(h)-len(suf):] == suf {
			return true
		}
	}
	return false
}

type serperRequest struct {
	Query string `json:"q"`
	Num   int    `json:"num"`
	Start int    `json:"start"`
}

type serperResponse struct {
	Organic []struct {
		Link string `json:"link"`
	} `json:"organic"`
}

func SerpAPISearch(apiKey, q string, max int, start int) ([]string, error) {
	if apiKey == "" {
		return nil, errors.New("SERPER_API missing")
	}
	
	// Use Serper API instead of SerpAPI
	client := &http.Client{Timeout: 20 * time.Second}
	
	// Create JSON payload for Serper API
	payload := serperRequest{
		Query: q,
		Num:   20,
		Start: start,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", "https://google.serper.dev/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var sr serperResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}
	
	out := make([]string, 0, max)
	seen := map[string]bool{}
	for _, r := range sr.Organic {
		if r.Link == "" || seen[r.Link] {
			continue
		}
		u, err := url.Parse(r.Link)
		if err != nil {
			continue
		}
		if hostAllowed(u) {
			seen[r.Link] = true
			out = append(out, r.Link)
			if len(out) >= max {
				break
			}
		}
	}
	return out, nil
}
