package api

import (
	"time"
)

func recalculateReplaceDate(now time.Time, taskDate, taskRepeat string) (string, error) {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if taskDate == "" {
		return now.Format(DateFormat), nil
	} else {
		date, err := time.Parse(DateFormat, taskDate)
		if err != nil {
			return "", err
		}

		if date.Before(now) && taskRepeat != "" {
			return NextDate(now, taskDate, taskRepeat)
		} else if date.Before(now) {
			return now.Format(DateFormat), nil
		} else {
			return taskDate, nil
		}
	}
}
