package models

import (
	"errors"
	"strings"
)

type Month struct {
	MonthInt	int
	MonthStr	string
	Year		string
}

func NewMonth(inputString string) (*Month, error) {
	dateFields := strings.Fields(inputString)
	if len(dateFields) != 2 {
		return nil, errors.New("Split month inputString not length of 2: " + inputString)
	}
	
	monthInt := monthToInt(dateFields[0])
	if monthInt == 0 {
		return nil, errors.New("Invalid month name: " + dateFields[0])
	}

	return &Month{MonthStr: dateFields[0], MonthInt: monthInt, Year: dateFields[1]}, nil
}

var monthMap = map[string]int{
	"January": 1, "February": 2, "March": 3, "April": 4,
	"May": 5, "June": 6, "July": 7, "August": 8,
	"September": 9, "October": 10, "November": 11, "December": 12,
}

func monthToInt(monthString string) int {
	if val, ok := monthMap[monthString]; ok {
		return val
	}
	
	return 0
}


