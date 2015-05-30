package date

import (
	"log"
	"testing"
)

func TestGetHourMinuteStr(t *testing.T) {
	dateStr := GetDateStr()
	yearStr := GetYearStr()
	hourStr, minStr := GetHourMinuteStr()

	log.Println(dateStr, yearStr, hourStr, minStr)
}
