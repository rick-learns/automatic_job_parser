package render

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"jobsite/internal/model"
)

type DailyPageData struct {
	SiteTitle string
	Day       string
	BaseURL   string
	Jobs      []model.Job
}

func WriteDaily(outDir string, siteTitle, baseURL string, jobs []model.Job) (string, error) {
	day := time.Now().Format("2006-01-02")
	dayDir := filepath.Join(outDir, day)
	_ = os.MkdirAll(dayDir, 0o755)
	// HTML generation removed - React frontend handles UI
	// if err := writeHTML(filepath.Join(dayDir, "index.html"), "templates/daily.html.tmpl", DailyPageData{
	// 	SiteTitle: siteTitle, Day: day, BaseURL: baseURL, Jobs: jobs,
	// }); err != nil {
	// 	return "", err
	// }
	if err := writeCSV(filepath.Join(dayDir, "jobs.csv"), jobs); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dayDir, "jobs.json"), jobs); err != nil {
		return "", err
	}
	latest := filepath.Join(outDir, "latest")
	_ = os.RemoveAll(latest)
	_ = copyDir(dayDir, latest)
	return dayDir, nil
}

func writeHTML(path, tplPath string, data any) error {
	t, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

func writeCSV(path string, jobs []model.Job) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"title", "company", "location", "salary_range", "url", "source", "discovered_date"})
	for _, j := range jobs {
		_ = w.Write([]string{j.Title, j.Company, j.Location, j.SalaryRaw, j.URL, j.Source, j.DiscoveredDate})
	}
	return nil
}

func writeJSON(path string, jobs []model.Job) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(jobs)
}
