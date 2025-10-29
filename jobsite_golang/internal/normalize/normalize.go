package normalize

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func CanonicalURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Fragment = ""
	q := u.Query()
	for k := range q {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "utm_") || strings.HasPrefix(lk, "gh_src") || strings.HasPrefix(lk, "lever-source") {
			q.Del(k)
		}
	}
	u.RawQuery = q.Encode()
	u.Host = strings.ToLower(u.Host)
	return u.String()
}

var moneyRe = regexp.MustCompile(`\d{2,3}(?:,\d{3})?`)
var kRe = regexp.MustCompile(`(\d{2,3})\s*k`)

func SalaryToRangeUSD(s string) (minPtr, maxPtr *int) {
	if s == "" {
		return nil, nil
	}
	if k := kRe.FindAllStringSubmatch(strings.ToLower(s), -1); len(k) >= 1 {
		nums := []int{}
		for _, m := range k {
			n, _ := strconv.Atoi(m[1])
			nums = append(nums, n*1000)
		}
		if len(nums) == 1 {
			nums = append(nums, nums[0])
		}
		min, max := nums[0], nums[len(nums)-1]
		return &min, &max
	}
	if m := moneyRe.FindAllString(s, -1); len(m) >= 1 {
		nums := []int{}
		for _, a := range m {
			v, _ := strconv.Atoi(strings.ReplaceAll(a, ",", ""))
			nums = append(nums, v)
		}
		if len(nums) == 1 {
			nums = append(nums, nums[0])
		}
		min, max := nums[0], nums[len(nums)-1]
		return &min, &max
	}
	return nil, nil
}

func IsRemoteUS(loc string, pageText string) bool {
	l := strings.ToLower(loc + " " + pageText)
	return strings.Contains(l, "remote - us") ||
		strings.Contains(l, "remote (us)") ||
		(strings.Contains(l, "remote") && (strings.Contains(l, "united states") || strings.Contains(l, " us")))
}
