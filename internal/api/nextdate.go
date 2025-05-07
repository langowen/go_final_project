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

		return nextWeekday(now, dStart, resRepeat[1])

	case "m":
		if len(resRepeat) < 2 {
			return "", fmt.Errorf("не указаны дни месяца")
		}

		months := ""

		if len(resRepeat) > 2 {
			months = resRepeat[2]
		}

		return nextMonthDay(now, dStart, resRepeat[1], months)

	default:
		return "", fmt.Errorf("недопустимый символ условия повторения задачи")
	}

}

func nextDay(now time.Time, dStart string, days int) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	startDate = startDate.AddDate(0, 0, days)

	for startDate.Before(now) {
		startDate = startDate.AddDate(0, 0, days)
	}

	return startDate.Format("20060102"), nil
}

func nextYear(now time.Time, dStart string) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	startDate = startDate.AddDate(1, 0, 0)

	for startDate.Before(now) {
		startDate = startDate.AddDate(1, 0, 0)
	}

	return startDate.Format("20060102"), nil
}

func nextWeekday(now time.Time, dStart string, daysStr string) (string, error) {
	startDate, err := time.Parse("20060102", dStart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %v", err)
	}

	dayStrs := strings.Split(daysStr, ",")
	weekdays := make(map[int]bool, len(dayStrs))
	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("недопустимый день недели: %s", dayStr)
		}
		weekdays[day] = true
	}

	var current time.Time
	if startDate.After(now) {
		current = startDate
	} else {
		current = now.AddDate(0, 0, 1)
	}

	maxDate := current.AddDate(0, 1, 0)

	for current.Before(maxDate) {
		wd := int(current.Weekday())
		if wd == 0 {
			wd = 7
		}

		if weekdays[wd] && current.After(startDate) {
			return current.Format("20060102"), nil
		}

		current = current.AddDate(0, 0, 1)
	}

	return "", nil
}

func nextMonthDay(now time.Time, dStart string, daysStr string, monthsStr string) (string, error) {
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
	maxDate := current.AddDate(2, 0, 0)

	for current.Before(maxDate) {
		possibleDates := generatePossibleDates(current.Year(), current.Month(), days)
		sort.Slice(possibleDates, func(i, j int) bool {
			return possibleDates[i].Before(possibleDates[j])
		})

		for _, date := range possibleDates {
			if (months == nil || months[int(date.Month())]) && date.After(startDate) && !date.Before(now) {
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

func parseMonths(monthsStr string) (map[int]bool, error) {
	if monthsStr == "" {
		return nil, nil
	}
	months := make(map[int]bool)
	for _, monthStr := range strings.Split(monthsStr, ",") {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			return nil, fmt.Errorf("недопустимый месяц: %s", monthStr)
		}
		months[month] = true
	}
	return months, nil
}

func generatePossibleDates(year int, month time.Month, days []int) []time.Time {
	possibleDates := make([]time.Time, 0, len(days))
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for _, day := range days {
		var targetDate time.Time
		switch {
		case day == -1:
			targetDate = time.Date(year, month, lastDayOfMonth, 0, 0, 0, 0, time.UTC)
		case day == -2:
			targetDate = time.Date(year, month, lastDayOfMonth-1, 0, 0, 0, 0, time.UTC)
		case day > 0 && day <= lastDayOfMonth:
			targetDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		default:
			continue
		}
		possibleDates = append(possibleDates, targetDate)
	}
	return possibleDates
}
