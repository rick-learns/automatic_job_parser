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

	serpKey := os.Getenv("SERPER_API")
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
		runDaily(db, serpKey, []QueryConfig{}, outDir, siteTitle, baseURL)
	case "seed":
		loadSeed(db, outDir, siteTitle, baseURL)
	default:
		log.Fatalf("unknown command: %s", mode)
	}
}

func runDaily(db *store.DB, serpKey string, queries []QueryConfig, outDir, siteTitle, baseURL string) {
	newJobsCount := 0
	updatedJobsCount := 0
	seen := map[string]bool{}

	for i, cfg := range queries {
		log.Printf("Query %d/%d (Tier %d, %d pages): %s", i+1, len(queries), cfg.Tier, cfg.Pages, cfg.Query)
		
		// Fetch all pages for this query
		var allLinks []string
		for page := 0; page < cfg.Pages; page++ {
			start := page * 20
			links, err := search.SerpAPISearch(serpKey, cfg.Query, 20, start)
			if err != nil {
				log.Printf("search error on page %d: %v", page+1, err)
				continue
			}
			allLinks = append(allLinks, links...)
			log.Printf("Page %d/%d: Found %d links (total: %d)", page+1, cfg.Pages, len(links), len(allLinks))
		}
		
		log.Printf("Total %d unique links from query %d (Tier %d)", len(allLinks), i+1, cfg.Tier)
		
		for _, link := range allLinks {
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

type QueryConfig struct {
	Query string
	Pages int
	Tier  int
}

func defaultQueries() []QueryConfig {
	return []QueryConfig{
		// Optimized for 20-25 new jobs/day: 6 queries × 3 pages = 18 API calls
		// 1. Core SDET/QA Automation - highest signal
		{`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:apply.workable.com) ("Senior Quality Engineer" OR SDET OR "QA Automation" OR "Software Development Engineer in Test" OR "Test Automation Engineer") ("Remote" OR "United States" OR "US") -intern -internship -contract -temporary -freelance -agency`, 3, 1},
		// 2. Mobile QA - Appium, iOS, Android focus  
		{`(site:boards.greenhouse.io OR site:jobs.lever.co OR site:apply.workable.com) (Appium OR "Mobile QA" OR "iOS QA" OR "Android QA") (SDET OR "Quality Engineer" OR "QA Automation") ("Remote" OR "United States" OR "US") -intern -contract -temporary`, 3, 1},
		// 3. Senior/Staff titles - high-value roles
		{`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co) (intitle:Senior OR intitle:Staff OR intitle:Lead) (SDET OR "Quality Engineer" OR "QA Automation") ("Remote" OR "United States" OR "US") -intern -contract -temporary`, 3, 1},
		// 4. Playwright web automation
		{`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co OR site:apply.workable.com) (Playwright) (QA OR SDET OR "Test Automation" OR "Software Engineer in Test") ("Remote" OR "United States" OR "US") -intern -contract -temporary`, 3, 1},
		// 5. macOS/Desktop/Endpoint testing
		{`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co) (macOS OR "desktop client" OR "endpoint agent" OR "device management") (QA OR "Quality Engineer" OR SDET OR "Test Automation") ("Remote" OR "United States" OR "US") -intern -contract -temporary`, 3, 1},
		// 6. CI/CD & GitHub Actions automation
		{`(site:boards.greenhouse.io OR site:jobs.ashbyhq.com OR site:jobs.lever.co) ("QA Automation" OR SDET OR "Quality Engineer") ("CI/CD" OR "continuous integration" OR "GitHub Actions" OR "release readiness") ("Remote" OR "United States" OR "US") -intern -contract -temporary`, 3, 1},
	}
}

// getQueries returns queries from JOBSITE_QUERIES env var or defaults
func getQueries() []QueryConfig {
	envQueries := os.Getenv("JOBSITE_QUERIES")
	if envQueries != "" {
		// Split by comma and trim
		parts := strings.Split(envQueries, ",")
		queries := make([]QueryConfig, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				queries = append(queries, QueryConfig{Query: part, Pages: 3, Tier: 1})
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
