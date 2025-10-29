# Jobsite (Go)

Single-binary job finder for QA/SDET roles (Remote US + Wichita, KS). Crawls ATS postings found via X-ray queries, extracts fields, stores in SQLite, and renders a static HTML site + CSV/JSON. Includes a **seed mode** so you can see the UI immediately.

## Quick start

```bash
cd jobsite_golang
go mod download
go build -o jobsite ./cmd/jobsite

# Preview with seed data
./jobsite seed

# Serve the static site (public/latest) however you like
# Example (python): python3 -m http.server -d public/latest 8080
```

Files are emitted to `public/YYYY-MM-DD/` and mirrored to `public/latest/`.

## Daily run (real search)
Set your SerpAPI key:
```bash
export SERPAPI_API_KEY=YOUR_KEY
./jobsite daily
```

## Cron (Linux)
```cron
12 8 * * * /opt/jobsite/jobsite daily
55 23 * * 0 /opt/jobsite/jobsite weekly
```

## Config (.env.example)
- `SERPAPI_API_KEY`: Google search via SerpAPI
- `PUBLIC_DIR`: output folder (default `public`)
- `DB_PATH`: SQLite path (default `data/jobs.sqlite`)
- `SITE_TITLE`: page title
- `BASE_URL`: the subdomain you'll host on
