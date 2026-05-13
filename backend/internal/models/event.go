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

func NewEvent(title string, link string, t time.Time) *Event {
	return &Event{Title: title, Link: link, Time: t, NoOfPerformers: 0}
}

func (e *Event) AddOnePerformer() {
	e.NoOfPerformers += 1
}
