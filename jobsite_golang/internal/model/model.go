package model

type Job struct {
	URL            string `json:"url"`
	Title          string `json:"title"`
	Company        string `json:"company"`
	Location       string `json:"location"`
	SalaryRaw      string `json:"salary_raw"`
	SalaryMinUSD   *int   `json:"salary_min_usd"`
	SalaryMaxUSD   *int   `json:"salary_max_usd"`
	Source         string `json:"source"`
	PostedDate     string `json:"posted_date"`
	DiscoveredDate string `json:"discovered_date"`
	IsRemoteUS     bool   `json:"is_remote_us"`
	Tags           string `json:"tags"`
}
