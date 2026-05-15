package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	testTitle := "Opry 100"
	testLink := "www.fakeurl.com/opry-100"
	eventA := NewEvent(testTitle, testLink, time.Time{})

	assert.Equal(t, "Opry 100", eventA.Title)
	assert.Equal(t, "www.fakeurl.com/opry-100", eventA.Link)
	assert.Equal(t, int64(0), eventA.NoOfPerformers)
	assert.Equal(t, true, eventA.Upcoming)

	testTitle = "Charity Concert"
	testLink = "www.fakeurl.com/charity-concert"
	eventB := NewEvent(testTitle, testLink, time.Time{})

	assert.Equal(t, "Charity Concert", eventB.Title)
	assert.Equal(t, "www.fakeurl.com/charity-concert", eventB.Link)
	assert.Equal(t, int64(0), eventB.NoOfPerformers)
	assert.Equal(t, true, eventB.Upcoming)
}