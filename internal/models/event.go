package models

import (
	"time"
)

type Event struct {
	Id				   int64
	Link               string
	Title              string
	Time               time.Time
	NoOfPerformers	   int64
}

//ToDo: get the actual time at time of scrape (playwright vs colly)
func NewEvent(title string, link string) *Event {
	return &Event{Title: title, Link: link, Time: time.Now(), NoOfPerformers: 0}
}

func (e *Event) AddOnePerformer() {
	e.NoOfPerformers += 1
}
