package api

import (
	"time"
)

func recalculateReplaceDate(now time.Time, taskDate, taskRepeat string) (string, error) {
	now = now.Truncate(time.Hour * 24)
	if taskDate == "" {
		return now.Format(DateFormat), nil
	} else {
		date, err := time.Parse(DateFormat, taskDate)
		if err != nil {
			return "", err
		}

		if date.Before(now) && taskRepeat != "" {
			return NextDate(now, taskDate, taskRepeat)
		} else {
			return now.Format(DateFormat), nil
		}
	}
}
