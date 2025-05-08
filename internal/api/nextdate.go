package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dStart string, repeat string) (string, error) {
	resRepeat := strings.Split(repeat, " ")
	cmd := resRepeat[0]

	switch cmd {
	case "d":
		if len(resRepeat) != 2 {
			return "", fmt.Errorf("не указан интервал в днях")
		}

		num, err := strconv.Atoi(resRepeat[1])
		if err != nil {
			return "", fmt.Errorf("ошибка преобразования числа: %v", err)
		}

		if num > 400 || num <= 0 {
			return "", fmt.Errorf("недопустимое значение интервала в днях")
		}

		return nextDay(now, dStart, num)

	case "y":
		return nextYear(now, dStart)

	case "w":
		if len(resRepeat) != 2 {
			return "", fmt.Errorf("не указаны дни недели")
		}

		return nextWeek(now, dStart, resRepeat[1])

	case "m":
		if len(resRepeat) < 2 {
			return "", fmt.Errorf("не указаны дни месяца")
		}

		months := ""

		if len(resRepeat) > 2 {
			months = resRepeat[2]
		}

		return nextMonth(now, dStart, resRepeat[1], months)

	default:
		return "", fmt.Errorf("недопустимый символ условия повторения задачи")
	}

}

func nextDay(now time.Time, dStart string, days int) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	for {
		startDate = startDate.AddDate(0, 0, days)
		if afterNow(startDate, now) {
			break
		}
	}

	return startDate.Format("20060102"), nil
}

func nextYear(now time.Time, dStart string) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	for {
		startDate = startDate.AddDate(1, 0, 0)
		if afterNow(startDate, now) {
			break
		}
	}

	return startDate.Format("20060102"), nil
}

func nextWeek(now time.Time, dStart string, daysStr string) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	dayStrs := strings.Split(daysStr, ",")
	weekdays := make([]int, 0, len(dayStrs))
	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("недопустимый день недели: %s", dayStr)
		}
		weekdays = append(weekdays, day)
	}

	current := startDate
	if !afterNow(startDate, now) {
		current = now.AddDate(0, 0, 1)
	}

	maxDate := now.AddDate(0, 2, 0)

	for {
		if afterNow(current, maxDate) {
			break
		}

		wd := int(current.Weekday())
		if wd == 0 {
			wd = 7
		}

		if containsInt(weekdays, wd) && !afterNow(now, current) {
			return current.Format("20060102"), nil
		}

		current = current.AddDate(0, 0, 1)
	}

	return "", nil
}

func nextMonth(now time.Time, dStart string, daysStr string, monthsStr string) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	days, err := parseDays(daysStr)
	if err != nil {
		return "", err
	}

	months, err := parseMonths(monthsStr)
	if err != nil {
		return "", err
	}

	current := now.AddDate(0, 0, 1)
	maxDate := now.AddDate(2, 0, 0)

	for {
		if afterNow(current, maxDate) {
			break
		}

		possibleDates := genPossDates(current.Year(), current.Month(), days)

		sort.Slice(possibleDates, func(i, j int) bool {
			return possibleDates[i].Before(possibleDates[j])
		})

		for _, date := range possibleDates {
			if (len(months) == 0 || containsInt(months, int(date.Month()))) &&
				afterNow(date, startDate) &&
				!afterNow(now, date) {
				return date.Format("20060102"), nil
			}
		}
		current = current.AddDate(0, 1, 0)
	}

	return "", nil
}

func parseDays(daysStr string) ([]int, error) {
	dayStrs := strings.Split(daysStr, ",")
	days := make([]int, 0, len(dayStrs))
	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return nil, fmt.Errorf("недопустимый день месяца: %s", dayStr)
		}
		days = append(days, day)
	}
	return days, nil
}

func parseMonths(monthsStr string) ([]int, error) {
	if monthsStr == "" {
		return nil, nil
	}
	months := make([]int, 0)
	for _, monthStr := range strings.Split(monthsStr, ",") {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			return nil, fmt.Errorf("недопустимый месяц: %s", monthStr)
		}
		months = append(months, month)
	}
	return months, nil
}

func genPossDates(year int, month time.Month, days []int) []time.Time {
	possibleDates := make([]time.Time, 0, len(days))
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for _, day := range days {
		var date time.Time
		switch day {
		case -1:
			date = time.Date(year, month, lastDay, 0, 0, 0, 0, time.UTC)
		case -2:
			if lastDay > 1 {
				date = time.Date(year, month, lastDay-1, 0, 0, 0, 0, time.UTC)
			}
		default:
			if day > 0 && day <= lastDay {
				date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			}
		}
		if !date.IsZero() {
			possibleDates = append(possibleDates, date)
		}
	}
	return possibleDates
}

func containsInt(slice []int, target int) bool {
	for _, m := range slice {
		if m == target {
			return true
		}
	}
	return false
}

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date.After(now)
}
