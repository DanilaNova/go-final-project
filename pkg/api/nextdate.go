package api

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

var ErrNoRepeatRule = errors.New("repeat rule is empty")
var ErrWrongRepeatRuleFormat = errors.New("repeat rule is in the wrong format")
var ErrUnknownRepeatRule = errors.New("unknown repeat rule")

var ErrWrongRepeatRuleFormatDays = fmt.Errorf(`%w, expected "d <number>"`, ErrWrongRepeatRuleFormat)
var ErrWrongRepeatRuleFormatYear = fmt.Errorf(`%w, expected "y"`, ErrWrongRepeatRuleFormat)
var ErrWrongRepeatRuleFormatWeek = fmt.Errorf(`%w, expected "w <numbers separated by a comma>"`, ErrWrongRepeatRuleFormat)
var ErrWrongRepeatRuleFormatMonth = fmt.Errorf(`%w, expected "m <numbers separated by a comma> [numbers separated by a comma]"`, ErrWrongRepeatRuleFormat)

var ErrNumberOutOfRange = errors.New("number is out of range")
var ErrNoNextDateSolution = errors.New("no solution for next date with given rule")

func dateIsAfter(date, now time.Time) bool {
	date_year, date_month, date_day := date.Date()
	now_year, now_month, now_day := now.Date()

	return date_year > now_year || (date_year == now_year && (date_month > now_month || (date_month == now_month && date_day > now_day)))
}

func getDaysInMonth(date time.Time) int {
	switch date.Month() {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if date.Year()%4 == 0 {
			return 29
		}
		return 28
	}
	panic("unreachable")
}

// Parses set of integers separated by a comma
func parseIntegers(str string, min int, max int) ([]int, error) {
	split := strings.Split(str, ",")
	numbers := make([]int, len(split))
	for i, numberStr := range split {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}
		if min > number || number > max {
			return nil, fmt.Errorf("%w (%d - %d)", ErrNumberOutOfRange, min, max)
		}
		numbers[i] = number
	}

	return numbers, nil
}

// rule "d <number>"
func ruleDay(now time.Time, date time.Time, params []string) (time.Time, error) {
	if len(params) != 1 {
		return time.Time{}, ErrWrongRepeatRuleFormatDays
	}
	days, err := strconv.Atoi(params[0])
	if err != nil {
		return time.Time{}, fmt.Errorf(`%w, number parsing error: %w`, ErrWrongRepeatRuleFormatDays, err)
	}
	if 1 > days || days > 400 {
		return time.Time{}, fmt.Errorf(`%w, number of days is outside of supported range (1 - 400)`, ErrWrongRepeatRuleFormatDays)
	}

	date = date.AddDate(0, 0, days)

	for !dateIsAfter(date, now) {
		date = date.AddDate(0, 0, days)
	}

	return date, nil
}

// rule "y"
func ruleYear(now time.Time, date time.Time, params []string) (time.Time, error) {
	if len(params) != 0 {
		return time.Time{}, ErrWrongRepeatRuleFormatYear
	}

	date = date.AddDate(1, 0, 0)
	if !date.After(now) {
		date = date.AddDate(now.Year()-date.Year(), 0, 0)
	}
	if !date.After(now) {
		date = date.AddDate(1, 0, 0)
	}

	return date, nil
}

// rule "w <numbers separated by a comma>"
func ruleWeekday(now time.Time, date time.Time, params []string) (time.Time, error) {
	if len(params) != 1 {
		return time.Time{}, ErrWrongRepeatRuleFormatWeek
	}

	weekdays_split := strings.Split(params[0], ",")
	weekdays := make([]int, len(weekdays_split))
	for i, weekdayStr := range weekdays_split {
		weekday, err := strconv.Atoi(weekdayStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("%w, number parsing error: %w", ErrWrongRepeatRuleFormatWeek, err)
		}
		if 1 > weekday || weekday > 7 {
			return time.Time{}, fmt.Errorf(`%w, weekday is outside of supported range (1 - 7)`, ErrWrongRepeatRuleFormatWeek)
		}
		weekdays[i] = weekday
	}

	slices.Sort(weekdays)

	date = date.AddDate(0, 0, 1)
	if !dateIsAfter(date, now) {
		date = now.AddDate(0, 0, 1)
	}

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	if pos, found := slices.BinarySearch(weekdays, weekday); !found {
		if pos == len(weekdays) {
			return date.AddDate(0, 0, 7-weekday+weekdays[0]), nil
		}
		return date.AddDate(0, 0, weekdays[pos]-weekday), nil
	}

	return date, nil
}

// rule "m <numbers separated by a comma> [numbers separated by a comma]"
func ruleMonth(now time.Time, date time.Time, params []string) (time.Time, error) {
	if 1 > len(params) || len(params) > 2 {
		return time.Time{}, ErrWrongRepeatRuleFormatMonth
	}

	monthdays, err := parseIntegers(params[0], -2, 31)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w, days parsing error: %w", ErrWrongRepeatRuleFormatMonth, err)
	}

	var months []int
	if len(params) == 2 {
		months, err = parseIntegers(params[1], 1, 12)
		if err != nil {
			return time.Time{}, fmt.Errorf("%w, months parsing error: %w", ErrWrongRepeatRuleFormatMonth, err)
		}
	}

	slices.Sort(monthdays)
	slices.Sort(months)

	var monthdays_forward []int
	var monthdays_reverse = monthdays
	for i, v := range monthdays {
		if v >= 0 {
			monthdays_reverse = monthdays[:i]
			monthdays_forward = monthdays[i:]
			break
		}
	}

	date = date.AddDate(0, 0, 1)
	if !dateIsAfter(date, now) {
		date = now.AddDate(0, 0, 1)
	}
	end := date.AddDate(4, 0, 0) // Limit searching timespan to 4 years

	for !dateIsAfter(date, end) {
		month := int(date.Month())

		if pos, found := slices.BinarySearch(months, month); len(months) != 0 && !found {
			if pos == len(months) {
				date = date.AddDate(0, 12-month+months[0], 1-date.Day())
			} else {
				date = date.AddDate(0, months[pos]-month, 1-date.Day())
			}
		}

		daysInMonth := getDaysInMonth(date)

		i1, i2 := 0, 0

		monthdays_forward_limit := len(monthdays_forward)
		for monthdays_forward_limit > 0 && monthdays_forward[monthdays_forward_limit-1] > daysInMonth {
			monthdays_forward_limit--
		}

		for i1 < monthdays_forward_limit && i2 < len(monthdays_reverse) {
			v1, v2 := monthdays_forward[i1], monthdays_reverse[i2]+daysInMonth+1

			var selected_day int
			switch cmp.Compare(v1, v2) {
			case -1:
				i1++
				selected_day = v1
			case 0:
				i1++
				fallthrough
			case 1:
				i2++
				selected_day = v2
			}

			if selected_day >= date.Day() {
				return date.AddDate(0, 0, selected_day-date.Day()), nil
			}
		}

		for ; i1 < monthdays_forward_limit; i1++ {
			v := monthdays_forward[i1]
			if v >= date.Day() {
				return date.AddDate(0, 0, v-date.Day()), nil
			}
		}
		for ; i2 < len(monthdays_reverse); i2++ {
			v := daysInMonth + 1 + monthdays_reverse[i2]
			if v >= date.Day() {
				return date.AddDate(0, 0, v-date.Day()), nil
			}
		}

		// Start of the next month
		date = date.AddDate(0, 1, 1-date.Day())
	}

	return time.Time{}, ErrNoNextDateSolution
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrNoRepeatRule
	}

	date, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		return "", err
	}

	repeat_params := strings.Split(repeat, " ")

	switch repeat_params[0] {
	case "d":
		date, err = ruleDay(now, date, repeat_params[1:])
	case "y":
		date, err = ruleYear(now, date, repeat_params[1:])
	case "w":
		date, err = ruleWeekday(now, date, repeat_params[1:])
	case "m":
		date, err = ruleMonth(now, date, repeat_params[1:])
	default:
		err = fmt.Errorf(`%w (%s)`, ErrUnknownRepeatRule, repeat_params[0])
	}
	if err != nil {
		return "", err
	}

	return date.Format(db.DateFormat), nil
}

func nextDateHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.Header().Set("Allow", "GET")
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response.Header().Set("Content-Type", "text/plain")
	now, err := time.Parse(db.DateFormat, request.FormValue("now"))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_, err := response.Write([]byte(err.Error()))
		if err != nil {
			log.Println("ERROR: could not write error response: ", err)
		}
		return
	}
	date := request.FormValue("date")
	repeat := request.FormValue("repeat")

	date, err = NextDate(now, date, repeat)
	if err != nil {
		log.Println("ERROR: cannot solve next date: ", err)
		response.WriteHeader(http.StatusInternalServerError)
		_, err = response.Write([]byte(err.Error()))
		if err != nil {
			log.Println("ERROR: could not write error response: ", err)
		}
		return
	}

	response.WriteHeader(http.StatusOK)
	_, err = response.Write([]byte(date))
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
	}
}
