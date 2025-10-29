package extract

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var salaryPats = []*regexp.Regexp{
	regexp.MustCompile(`\$\s?\d{2,3}(?:,\d{3})?\s*[-–]\s*\$\s?\d{2,3}(?:,\d{3})?`),
	regexp.MustCompile(`\$\s?\d{2,3}(?:,\d{3})?\s*(?:per year|annually|/year)`),
	regexp.MustCompile(`\$\s?\d{2,3}k(?:\s*[-–]\s*\$\s?\d{2,3}k)?`),
}

type jsonLD struct {
	Title              string `json:"title"`
	HiringOrganization struct {
		Name string `json:"name"`
	} `json:"hiringOrganization"`
	DatePosted string `json:"datePosted"`
}

func FromHTML(html string) (title, company, location, salary, datePosted string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}

	if h1 := strings.TrimSpace(doc.Find("h1").First().Text()); h1 != "" {
		title = h1
	} else {
		title = strings.TrimSpace(doc.Find("title").First().Text())
	}

	// JSON-LD
	doc.Find(`script[type="application/ld+json"]`).Each(func(i int, s *goquery.Selection) {
		var jd jsonLD
		if err := json.Unmarshal([]byte(strings.TrimSpace(s.Text())), &jd); err == nil {
			if jd.HiringOrganization.Name != "" && company == "" {
				company = jd.HiringOrganization.Name
			}
			if jd.Title != "" && title == "" {
				title = jd.Title
			}
			if jd.DatePosted != "" && datePosted == "" {
				datePosted = jd.DatePosted
			}
		}
	})

	if company == "" {
		if x := strings.TrimSpace(doc.Find(".company, .company-name, .app-title small, .posting-headline h3").First().Text()); x != "" {
			company = x
		}
	}

	if x := strings.TrimSpace(doc.Find(":contains('Location')").Next().Text()); x != "" {
		location = x
	}
	if location == "" {
		doc.Find("li, p, span, div").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			t := strings.TrimSpace(s.Text())
			if t == "" {
				return true
			}
			lt := strings.ToLower(t)
			if strings.Contains(lt, "remote") || strings.Contains(t, "Wichita") || strings.Contains(t, "United States") {
				location = t
				return false
			}
			return true
		})
	}

	full := doc.Text()
	for _, re := range salaryPats {
		if m := re.FindString(full); m != "" {
			salary = strings.TrimSpace(m)
			break
		}
	}

	return
}
