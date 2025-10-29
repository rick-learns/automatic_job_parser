package search

import (
	"encoding/json"
	"errors"
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

type serpResp struct {
	OrganicResults []struct {
		Link string `json:"link"`
	} `json:"organic_results"`
}

func SerpAPISearch(apiKey, q string, max int) ([]string, error) {
	if apiKey == "" {
		return nil, errors.New("SERPAPI_API_KEY missing")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	values := url.Values{
		"engine":  {"google"},
		"q":       {q},
		"num":     {"20"},
		"hl":      {"en"},
		"gl":      {"us"},
		"api_key": {apiKey},
	}
	req, _ := http.NewRequest("GET", "https://serpapi.com/search.json?"+values.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var sr serpResp
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	out := make([]string, 0, max)
	seen := map[string]bool{}
	for _, r := range sr.OrganicResults {
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
