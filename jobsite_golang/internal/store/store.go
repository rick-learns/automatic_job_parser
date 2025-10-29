package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"

	"jobsite/internal/model"
)

type DB struct{ *sql.DB }

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Set optimal SQLite pragmas for reliability and performance
	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA synchronous=NORMAL;`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS jobs (
  id INTEGER PRIMARY KEY,
  url TEXT UNIQUE,
  title TEXT, company TEXT, location TEXT,
  salary_raw TEXT, salary_min_usd INTEGER, salary_max_usd INTEGER,
  source TEXT, posted_date TEXT,
  discovered_date TEXT NOT NULL,
  is_remote_us INTEGER NOT NULL,
  tags TEXT
);
CREATE TABLE IF NOT EXISTS runs (
  run_id TEXT PRIMARY KEY,
  started_at_utc TEXT, finished_at_utc TEXT,
  query_count INTEGER, new_links INTEGER, pages_parsed INTEGER
);`)
	return &DB{db}, err
}

// InsertJob inserts or updates a job. On conflict, updates all fields except discovered_date
func InsertJob(db *DB, j model.Job) error {
	result, err := db.Exec(`INSERT INTO jobs
(url,title,company,location,salary_raw,salary_min_usd,salary_max_usd,source,posted_date,discovered_date,is_remote_us,tags)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(url) DO UPDATE SET
  title=excluded.title,
  company=excluded.company,
  location=excluded.location,
  salary_raw=excluded.salary_raw,
  salary_min_usd=excluded.salary_min_usd,
  salary_max_usd=excluded.salary_max_usd,
  source=excluded.source,
  posted_date=excluded.posted_date,
  is_remote_us=excluded.is_remote_us,
  tags=excluded.tags`,
		j.URL, j.Title, j.Company, j.Location, j.SalaryRaw, j.SalaryMinUSD, j.SalaryMaxUSD,
		j.Source, j.PostedDate, j.DiscoveredDate, boolToInt(j.IsRemoteUS), j.Tags)

	// Track if this was a new insert or update
	if err == nil {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 1 {
			// New insert
		} else {
			// Updated existing record
		}
	}
	return err
}

// InsertJobStats tracks upsert statistics
type InsertJobStats struct {
	Inserted int64
	Updated  int64
}

// InsertJobWithStats performs upsert and returns counts
func InsertJobWithStats(db *DB, j model.Job) (InsertJobStats, error) {
	result, err := db.Exec(`INSERT INTO jobs
(url,title,company,location,salary_raw,salary_min_usd,salary_max_usd,source,posted_date,discovered_date,is_remote_us,tags)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(url) DO UPDATE SET
  title=excluded.title,
  company=excluded.company,
  location=excluded.location,
  salary_raw=excluded.salary_raw,
  salary_min_usd=excluded.salary_min_usd,
  salary_max_usd=excluded.salary_max_usd,
  source=excluded.source,
  posted_date=excluded.posted_date,
  is_remote_us=excluded.is_remote_us,
  tags=excluded.tags`,
		j.URL, j.Title, j.Company, j.Location, j.SalaryRaw, j.SalaryMinUSD, j.SalaryMaxUSD,
		j.Source, j.PostedDate, j.DiscoveredDate, boolToInt(j.IsRemoteUS), j.Tags)

	var stats InsertJobStats
	if err == nil {
		lastID, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()

		if lastID > 0 {
			stats.Inserted = 1
		} else {
			stats.Updated = rowsAffected
		}
	}
	return stats, err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func LastNDays(db *DB, days int) ([]model.Job, error) {
	rows, err := db.Query(`SELECT url,title,company,location,salary_raw,salary_min_usd,salary_max_usd,source,posted_date,discovered_date,is_remote_us,tags
FROM jobs WHERE date(discovered_date) >= date(?, '-'||?||' day') ORDER BY discovered_date DESC`, time.Now().UTC().Format("2006-01-02"), days-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Job
	for rows.Next() {
		var j model.Job
		var min, max sql.NullInt64
		var remote int
		if err := rows.Scan(&j.URL, &j.Title, &j.Company, &j.Location, &j.SalaryRaw, &min, &max, &j.Source, &j.PostedDate, &j.DiscoveredDate, &remote, &j.Tags); err != nil {
			return nil, err
		}
		if min.Valid {
			v := int(min.Int64)
			j.SalaryMinUSD = &v
		}
		if max.Valid {
			v := int(max.Int64)
			j.SalaryMaxUSD = &v
		}
		j.IsRemoteUS = remote == 1
		out = append(out, j)
	}
	return out, nil
}
