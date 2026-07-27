package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	n := r.FormValue("now")
	now, err := time.Parse(DateFormat, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", fmt.Errorf("repeat param can not be empty")
	}

	if len(dstart) == 0 {
		return "", fmt.Errorf("start date can not be empty")
	}

	tStart, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	repeatSlice := strings.Split(repeat, " ")

	var newDate time.Time

	switch repeatSlice[0] {
	case "y":
		newDate = tStart.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
	case "d":
		if len(repeatSlice) != 2 {
			return "", fmt.Errorf("invalid arguments count")
		}
		days, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", fmt.Errorf("invalid repeat param")
		}

		newDate = tStart.AddDate(0, 0, days)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}
	case "w":
		if len(repeatSlice) != 2 {
			return "", fmt.Errorf("invalid arguments count")
		}

		weekdays, err := extractNumbers(repeatSlice[1])
		if err != nil {
			return "", nil
		}

		newDate = now             // Т.к. нас интересуют следующие числа, то считаем сразу от текущей даты
		for i := 1; i <= 7; i++ { // Достаточно проитерировать только 7 дней
			newDate = newDate.AddDate(0, 0, 1)
			if slices.Contains(weekdays, weekdayToNum(newDate.Weekday())) {
				break
			}
		}
	case "m":
		argc := len(repeatSlice)
		if argc < 2 {
			return "", fmt.Errorf("invalid arguments count")
		}

		months := make([]int, 0)
		if argc == 3 {
			months, err = extractNumbers(repeatSlice[2])
			if err != nil {
				return "", err
			}
		}

		newDate = tStart
		for {
			y, m, d := newDate.Date()
			lastDay := daysInMonth(m, y)
			days, err := extractNumbers(repeatSlice[1])
			if err != nil {
				return "", err
			}
			replaceNegativeDays(days, lastDay)
			slices.Sort(days)

			if slices.Contains(months, int(m)) || len(months) == 0 {
				for _, v := range days {
					if v == d && newDate.After(now) {
						return newDate.Format(DateFormat), nil
					}
				}
			}
			newDate = newDate.AddDate(0, 0, 1)
		}
	default:
		return "", fmt.Errorf("invalid repeat param")
	}

	return newDate.Format(DateFormat), nil
}

func weekdayToNum(w time.Weekday) int {
	return (int(w)+6)%7 + 1
}

func extractNumbers(src string) ([]int, error) {
	strNums := strings.Split(src, ",")

	nums := make([]int, 0)
	for _, v := range strNums {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		nums = append(nums, num)
	}

	return nums, nil
}

func daysInMonth(m time.Month, year int) int {
	t := time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}

func replaceNegativeDays(d []int, lastDay int) {
	for i, v := range d {
		switch v {
		case -1:
			d[i] = lastDay
		case -2:
			d[i] = lastDay - 1
		}
	}
}
