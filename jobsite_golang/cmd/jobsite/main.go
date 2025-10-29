package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"jobsite/internal/extract"
	"jobsite/internal/fetch"
	"jobsite/internal/lock"
	"jobsite/internal/model"
	"jobsite/internal/normalize"
	"jobsite/internal/render"
	"jobsite/internal/search"
	"jobsite/internal/store"
)

// Version info (set via ldflags during build)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Parse command-line flags
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help")
	dbPathFlag := flag.String("db", "", "Database file path (default: data/jobs.sqlite)")
	outDirFlag := flag.String("out-dir", "", "Output directory (default: public)")
	lockFileFlag := flag.String("lock-file", "jobsite.lock", "Lock file path")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("jobsite version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Show help
	if *showHelp || len(flag.Args()) == 0 {
		fmt.Println("Jobsite - QA/SDET Job Parser")
		fmt.Println("\nUsage: jobsite [MODE] [FLAGS]")
		fmt.Println("\nModes:")
		fmt.Println("  daily   - Run daily job search and update")
		fmt.Println("  weekly  - Run weekly summary")
		fmt.Println("  seed    - Load seed data for testing")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Get mode (first positional argument)
	mode := "daily"
	if len(flag.Args()) > 0 {
		mode = flag.Args()[0]
	}

	serpKey := os.Getenv("SERPAPI_API_KEY")
	outDir := getenv("PUBLIC_DIR", "public")
	if *outDirFlag != "" {
		outDir = *outDirFlag
	}
	dbPath := getenv("DB_PATH", "data/jobs.sqlite")
	if *dbPathFlag != "" {
		dbPath = *dbPathFlag
	}
	siteTitle := getenv("SITE_TITLE", "QA/SDET Roles (Remote US + Wichita)")
	baseURL := getenv("BASE_URL", "https://jobs.example.com")

	// Acquire lock for daily/weekly runs
	var lck *lock.Lock
	if mode == "daily" || mode == "weekly" {
		var err error
		lck, err = lock.Acquire(*lockFileFlag)
		if err != nil {
			log.Fatalf("Failed to acquire lock: %v", err)
		}
		log.Printf("Lock acquired: %s", *lockFileFlag)
		defer func() {
			if err := lck.Release(); err != nil {
				log.Printf("Failed to release lock: %v", err)
			} else {
				log.Printf("Lock released: %s", *lockFileFlag)
			}
		}()
	}

	db, err := store.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "daily":
		queries := getQueries()
		log.Printf("Using %d search queries", len(queries))
		runDaily(db, serpKey, queries, outDir, siteTitle, baseURL)
	case "weekly":
		// TODO: implement a weekly HTML writer; for now reuse daily with last 7 days
		runDaily(db, serpKey, []string{}, outDir, siteTitle, baseURL)
	case "seed":
		loadSeed(db, outDir, siteTitle, baseURL)
	default:
		log.Fatalf("unknown command: %s", mode)
	}
}

func runDaily(db *store.DB, serpKey string, queries []string, outDir, siteTitle, baseURL string) {
	newJobsCount := 0
	updatedJobsCount := 0
	seen := map[string]bool{}

	for i, q := range queries {
		log.Printf("Query %d/%d: %s", i+1, len(queries), q)
		links, err := search.SerpAPISearch(serpKey, q, 20)
		if err != nil {
			log.Printf("search error: %v", err)
			continue
		}
		log.Printf("Found %d links from query %d", len(links), i+1)
		for _, link := range links {
			canon := normalize.CanonicalURL(link)
			if seen[canon] {
				continue
			}
			seen[canon] = true

			html, err := fetch.Get(canon)
			if err != nil {
				log.Printf("fetch %s: %v", canon, err)
				continue
			}

			title, company, location, salary, posted := extract.FromHTML(html)
			min, max := normalize.SalaryToRangeUSD(salary)
			isRemote := normalize.IsRemoteUS(location, html)

			j := model.Job{
				URL: canon, Title: strings.TrimSpace(title), Company: strings.TrimSpace(company),
				Location: strings.TrimSpace(location), SalaryRaw: strings.TrimSpace(salary),
				SalaryMinUSD: min, SalaryMaxUSD: max,
				Source: sourceFromURL(canon), PostedDate: posted,
				DiscoveredDate: time.Now().UTC().Format("2006-01-02"),
				IsRemoteUS:     isRemote,
				Tags:           "appium,playwright,ci-cd,macos,ios,android",
			}
			stats, err := store.InsertJobWithStats(db, j)
			if err != nil {
				continue
			}
			if stats.Inserted > 0 {
				newJobsCount++
			} else if stats.Updated > 0 {
				updatedJobsCount++
			}
		}
	}

	jobs, err := store.LastNDays(db, 7)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New jobs inserted this run: %d", newJobsCount)
	log.Printf("Existing jobs updated this run: %d", updatedJobsCount)
	log.Printf("Total jobs in database (last 7 days): %d", len(jobs))
	dayDir, err := render.WriteDaily(outDir, siteTitle, baseURL, jobs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote:", dayDir)
}

func loadSeed(db *store.DB, outDir, siteTitle, baseURL string) {
	f, err := os.Open("data/seed.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var jobs []model.Job
	if err := dec.Decode(&jobs); err != nil {
		log.Fatal(err)
	}
	for _, j := range jobs {
		_ = store.InsertJob(db, j)
	}
	jobs7, err := store.LastNDays(db, 7)
	if err != nil {
		log.Fatal(err)
	}
	_, err = render.WriteDaily(outDir, siteTitle, baseURL, jobs7)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeded and rendered /public/latest")
	_ = exec.Command("bash", "-lc", "ls -la public/latest").Run()
}

func sourceFromURL(u string) string {
	switch {
	case strings.Contains(u, "greenhouse.io"):
		return "Greenhouse"
	case strings.Contains(u, "ashbyhq.com"):
		return "Ashby"
	case strings.Contains(u, "lever.co"):
		return "Lever"
	case strings.Contains(u, "myworkdayjobs"):
		return "Workday"
	case strings.Contains(u, "smartrecruiters"):
		return "SmartRecruiters"
	case strings.Contains(u, "workable.com"):
		return "Workable"
	default:
		return "ATS"
	}
}

func defaultQueries() []string {
	return []string{
		`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:myworkdayjobs.com OR site:jobs.smartrecruiters.com OR site:apply.workable.com) ("Senior Quality Engineer" OR SDET OR "QA Automation" OR "Software Development Engineer in Test" OR "Test Automation Engineer") (Appium OR Playwright OR "GitHub Actions" OR macOS OR iOS OR Android OR Golang) ("Remote" OR "United States" OR "US") -intern -internship -contract -temporary -freelance -agency`,
		`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:myworkdayjobs.com OR site:jobs.smartrecruiters.com OR site:apply.workable.com) (Appium OR "Mobile QA" OR "iOS QA" OR "Android QA" OR "mobile test automation") ("Senior" OR Lead OR Staff OR "Quality Engineer" OR SDET) ("Remote" OR "United States" OR "US") -intern -contract -temporary`,
		`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:myworkdayjobs.com OR site:jobs.smartrecruiters.com OR site:apply.workable.com) ("QA Automation" OR SDET OR "Quality Engineer") ("CI/CD" OR "continuous integration" OR "GitHub Actions" OR "release readiness" OR "risk-based testing") ("Remote" OR "United States" OR "US") -intern -contract -temporary`,
		`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:myworkdayjobs.com OR site:jobs.smartrecruiters.com OR site:apply.workable.com) (macOS OR "desktop client" OR "endpoint agent" OR "device management") (QA OR "Quality Engineer" OR SDET OR "Test Automation") ("Remote" OR "United States" OR "US") -intern -contract -temporary`,
		`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:myworkdayjobs.com OR site:jobs.smartrecruiters.com OR site:apply.workable.com) ("Quality Engineer" OR SDET OR "QA Automation" OR "Test Engineer") ("Wichita" OR "KS" OR "Kansas") -intern -internship -contract -temporary`,
	}
}

// getQueries returns queries from JOBSITE_QUERIES env var or defaults
func getQueries() []string {
	envQueries := os.Getenv("JOBSITE_QUERIES")
	if envQueries != "" {
		// Split by comma and trim
		parts := strings.Split(envQueries, ",")
		queries := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				queries = append(queries, part)
			}
		}
		if len(queries) > 0 {
			log.Printf("Using queries from JOBSITE_QUERIES env var")
			return queries
		}
	}
	return defaultQueries()
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
