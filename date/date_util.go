package date

import (
	"fmt"
	"time"
)

func GetYearStr() string {
	now := time.Now()
	year, _, _ := now.UTC().Date()

	return fmt.Sprintf("%04d", year)
}

func GetDateStr() string {
	now := time.Now()
	year, mon, day := now.UTC().Date()

	return fmt.Sprintf("%04d%02d%02d", year, mon, day)
}

func GetHourMinuteStr() (string, string) {
	hour := time.Now().Hour()
	minute := time.Now().Minute()

	return fmt.Sprintf("%02d", hour), fmt.Sprintf("%02d", minute)
}
