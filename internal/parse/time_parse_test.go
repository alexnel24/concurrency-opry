package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeFromLink_HourOnlyPM(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-06-02-grand-ole-opry-opry-100-at-7-pm")
	expected := time.Date(2026, time.June, 2, 19, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestParseTimeFromLink_HourAndMinutesPM(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-06-02-grand-ole-opry-opry-100-at-9-30-pm")
	expected := time.Date(2026, time.June, 2, 21, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestParseTimeFromLink_12AM(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-07-03-grand-ole-opry-opry-100-at-12-am")
	expected := time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestParseTimeFromLink_HourOnlyAM(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-11-27-grand-ole-opry-opry-100-at-7-am")
	expected := time.Date(2026, time.November, 27, 7, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestParseTimeFromLink_12PM(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-06-06-opry-at-the-ryman-opry-100-at-12-pm")
	expected := time.Date(2026, time.June, 6, 12, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, result)
}

func TestParseTimeFromLink_UndatedURL(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/grand-ole-opry-opry-100-1074602")
	assert.True(t, result.IsZero())
}

func TestParseTimeFromLink_NoSlash(t *testing.T) {
	result := ParseTimeFromLink("grand-ole-opry")
	assert.True(t, result.IsZero())
}

func TestParseTimeFromLink_ValidDateNoTime(t *testing.T) {
	result := ParseTimeFromLink("https://www.opry.com/show/2026-06-02-grand-ole-opry")
	assert.True(t, result.IsZero())
}

func TestParseDateTimeFromText_StandardPM(t *testing.T) {
	result := ParseDateTimeFromText("Date & Time Tuesday, June 25 7:00 PM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.June, 25, 19, 0, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.June, 25, 19, 0, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_MinutesAndPM(t *testing.T) {
	result := ParseDateTimeFromText("Tuesday, June 2 9:30 PM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.June, 2, 21, 30, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.June, 2, 21, 30, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_AM(t *testing.T) {
	result := ParseDateTimeFromText("Monday, March 5 7:00 AM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.March, 5, 7, 0, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.March, 5, 7, 0, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_12PM(t *testing.T) {
	result := ParseDateTimeFromText("Friday, August 1 12:00 PM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.August, 1, 12, 0, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.August, 1, 12, 0, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_12AM(t *testing.T) {
	result := ParseDateTimeFromText("Saturday, October 4 12:00 AM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.October, 4, 0, 0, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.October, 4, 0, 0, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_WhitespaceNormalization(t *testing.T) {
	result := ParseDateTimeFromText("Date  &  Time\n  June   25\n  7:00 PM")
	now := time.Now().UTC()
	expected := time.Date(now.Year(), time.June, 25, 19, 0, 0, 0, time.UTC)
	if expected.Before(now) {
		expected = time.Date(now.Year()+1, time.June, 25, 19, 0, 0, 0, time.UTC)
	}
	assert.Equal(t, expected, result)
}

func TestParseDateTimeFromText_NoDate(t *testing.T) {
	result := ParseDateTimeFromText("7:00 PM")
	assert.True(t, result.IsZero())
}

func TestParseDateTimeFromText_NoTime(t *testing.T) {
	result := ParseDateTimeFromText("Tuesday, May 19")
	assert.True(t, result.IsZero())
}

func TestParseDateTimeFromText_EmptyString(t *testing.T) {
	result := ParseDateTimeFromText("")
	assert.True(t, result.IsZero())
}
