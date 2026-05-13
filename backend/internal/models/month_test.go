package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMonth(t *testing.T) {
	testString := "November 2025"
	testMonth, err := NewMonth(testString)	
	
	assert.Nil(t, err)
	assert.Equal(t, "November", testMonth.MonthStr)
	assert.Equal(t, 11, testMonth.MonthInt)
	assert.Equal(t, "2025", testMonth.Year)

	test2String := "May 2025"
	testMonth, err = NewMonth(test2String)

	assert.Nil(t, err)
	assert.Equal(t, "May", testMonth.MonthStr)
	assert.Equal(t, 5, testMonth.MonthInt)
	assert.Equal(t, "2025", testMonth.Year)

	testLowercaseMonth := "november 2025"
	testMonth, err = NewMonth(testLowercaseMonth)
	
	assert.NotNil(t, err)
	assert.Nil(t, testMonth)

	testShortMonth := "Dec 2025"
	testMonth, err = NewMonth(testShortMonth)
	
	assert.NotNil(t, err)
	assert.Nil(t, testMonth)

	testFakeYear := "June 20456"
	testMonth, err = NewMonth(testFakeYear)

	assert.Nil(t, err)
	assert.Equal(t, "June", testMonth.MonthStr)
	assert.Equal(t, 6, testMonth.MonthInt)
	assert.Equal(t, "20456", testMonth.Year)

	testWordYear := "June TwoThousandTwentyFive"
	testMonth, err = NewMonth(testWordYear)

	assert.Nil(t, err)
	assert.Equal(t, "June", testMonth.MonthStr)
	assert.Equal(t, 6, testMonth.MonthInt)
	assert.Equal(t, "TwoThousandTwentyFive", testMonth.Year)

	
	testNoYear := "July"
	testMonth, err = NewMonth(testNoYear)

	assert.NotNil(t, err)
	assert.Nil(t, testMonth)
}
