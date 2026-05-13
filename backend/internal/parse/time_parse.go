package parse

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var monthNames = map[string]time.Month{
	"january": time.January, "february": time.February,
	"march": time.March, "april": time.April,
	"may": time.May, "june": time.June,
	"july": time.July, "august": time.August,
	"september": time.September, "october": time.October,
	"november": time.November, "december": time.December,
}

// matches "Month D[D]" optionally preceded by a day-of-week like "Tuesday, "
var textMonthDayPattern = regexp.MustCompile(`(?i)(?:(?:Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday),?\s+)?(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2})`)
var textTimePattern = regexp.MustCompile(`(?i)(\d{1,2}):(\d{2})\s*(AM|PM)`)

func ParseTimeFromLink(link string) time.Time {
	idx := strings.LastIndex(link, "/")
	if idx == -1 || len(link) <= idx+1 || link[idx+1] < '0' || link[idx+1] > '9' {
		return time.Time{}
	}
	slug := link[idx+1:]

	if len(slug) < 10 {
		return time.Time{}
	}
	date, err := time.Parse(time.DateOnly, slug[:10])
	if err != nil {
		return time.Time{}
	}

	atIdx := strings.LastIndex(slug, "-at-")
	if atIdx == -1 {
		return time.Time{}
	}
	tokens := strings.Split(slug[atIdx+4:], "-")

	var hour, min int
	var ampm string
	switch len(tokens) {
	case 2:
		hour, _ = strconv.Atoi(tokens[0])
		ampm = tokens[1]
	case 3:
		hour, _ = strconv.Atoi(tokens[0])
		min, _ = strconv.Atoi(tokens[1])
		ampm = tokens[2]
	default:
		return time.Time{}
	}

	if ampm == "pm" && hour != 12 {
		hour += 12
	} else if ampm == "am" && hour == 12 {
		hour = 0
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, min, 0, 0, time.UTC)
}

func ParseDateTimeFromText(text string) time.Time {
	text = strings.Join(strings.Fields(text), " ")

	dateMatch := textMonthDayPattern.FindStringSubmatch(text)
	timeMatch := textTimePattern.FindStringSubmatch(text)
	if dateMatch == nil || timeMatch == nil {
		return time.Time{}
	}

	month, ok := monthNames[strings.ToLower(dateMatch[1])]
	if !ok {
		return time.Time{}
	}

	day, _ := strconv.Atoi(dateMatch[2])
	hour, _ := strconv.Atoi(timeMatch[1])
	min, _ := strconv.Atoi(timeMatch[2])
	ampm := strings.ToUpper(timeMatch[3])

	if ampm == "PM" && hour != 12 {
		hour += 12
	} else if ampm == "AM" && hour == 12 {
		hour = 0
	}

	now := time.Now().UTC()
	t := time.Date(now.Year(), month, day, hour, min, 0, 0, time.UTC)
	if t.Before(now) {
		t = time.Date(now.Year()+1, month, day, hour, min, 0, 0, time.UTC)
	}
	return t
}
