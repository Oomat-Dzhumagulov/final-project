package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat не указан")
	}

	t, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("dstart некоректен: %w", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное число: %s", parts[1])
		}
		if days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимый диапазон (1-400)")
		}

		t = NextOccur(t, now, days, 0, 0)
		return t.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат y")
		}
		t = NextOccur(t, now, 0, 0, 1)
		return t.Format(DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат w")
		}
		weekdays, err := parseWebDays(parts[1])
		if err != nil {
			return "", fmt.Errorf("ошибка при парсинге: %w", err)
		}
		startCheck := t
		if startCheck.Before(now) {
			startCheck = now
		}

		for i := 0; i < 366; i++ {
			if IsAfter(startCheck, now) {
				wd := int(startCheck.Weekday())
				if wd == 0 {
					wd = 7
				}

				if _, ok := weekdays[wd]; ok {
					return startCheck.Format(DateFormat), nil
				}
			}
			startCheck = startCheck.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат m")
		}

		daysMap, err := parseMonthDays(parts[1])
		if err != nil {
			return "", fmt.Errorf("ошибка при парсинге: %w", err)
		}

		monthsMap := make(map[int]bool)
		if len(parts) == 3 {
			mm, err := parseMonths(parts[2])
			if err != nil {
				return "", fmt.Errorf("ошибка при парсирге: %w", err)
			}
			monthsMap = mm
		}

		startCheck := t
		if startCheck.Before(now) {
			startCheck = now.AddDate(0, 0, 1)
		}

		for i := 0; i < 366*5; i++ {
			startCheck = startCheck.AddDate(0, 0, 1)

			if !IsAfter(startCheck, now) {
				continue
			}

			m := int(startCheck.Month())
			if len(monthsMap) > 0 {
				if _, ok := monthsMap[m]; !ok {
					continue
				}
			}

			d := startCheck.Day()

			lastDay := time.Date(startCheck.Year(), startCheck.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			found := false

			for reqDay := range daysMap {
				if reqDay > 0 {
					if reqDay == d {
						found = true
						break
					}
				} else {
					if lastDay+reqDay+1 == d {
						found = true
						break
					}
				}
			}
			if found {
				return startCheck.Format(DateFormat), nil
			}
		}
		return "", fmt.Errorf("дата для m не найдена")
	default:
		return "", fmt.Errorf("неизвестный формат")
	}
	return "", fmt.Errorf("дата не найдена")
}

func IsAfter(t1, t2 time.Time) bool {
	d1 := t1.Format(DateFormat)
	d2 := t2.Format(DateFormat)
	return d1 > d2
}

func NextOccur(t, now time.Time, days, months, years int) time.Time {
	t = t.AddDate(years, months, days)

	for !IsAfter(t, now) {
		t = t.AddDate(years, months, days)
	}
	return t
}

func parseWebDays(s string) (map[int]bool, error) {
	res := make(map[int]bool)
	parts := strings.Split(s, ",")
	for _, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 || v > 7 {
			return nil, fmt.Errorf("неверный день недели: %s", p)
		}
		res[v] = true
	}
	return res, nil
}

func parseMonthDays(s string) (map[int]bool, error) {
	res := make(map[int]bool)
	parts := strings.Split(s, ",")
	for _, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования: %s", p)
		}
		if v == 0 || v < -2 || v > 31 {
			return nil, fmt.Errorf("неверный день месяца: %d", v)
		}
		res[v] = true
	}
	return res, nil
}

func parseMonths(s string) (map[int]bool, error) {
	res := make(map[int]bool)
	parts := strings.Split(s, ",")
	for _, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 || v > 12 {
			return nil, fmt.Errorf("неверный месяц:n %s", p)
		}
		res[v] = true
	}
	return res, nil
}
