package main

import (
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

const (
	day = time.Hour * 24
)

const (
	fullerHuman = "Mon, 02 Jan 2006"
	fullHuman   = "Mon, 02 Jan"
	fuller      = "02-01-2006"
	full        = "02-01"
	weekday     = "Monday"
)

var (
	fullFormats = [...]string{fuller, full}
	relative    = [...]string{"Yesterday", "Today", "Tomorrow"}
	weekdays    = map[string]time.Weekday{
		time.Sunday.String():    time.Sunday,
		time.Monday.String():    time.Monday,
		time.Tuesday.String():   time.Tuesday,
		time.Wednesday.String(): time.Wednesday,
		time.Thursday.String():  time.Thursday,
		time.Friday.String():    time.Friday,
		time.Saturday.String():  time.Saturday,
	}
)

func parseDueDate(due string) (date time.Time, err error) {
	if len(due) < 3 {
		err = errors.New("due string too short")
		return
	}

	due = strings.ToLower(due)

	// relative date
	for i, rd := range relative {
		if strings.HasPrefix(strings.ToLower(rd), due) {
			t := time.Now()
			t = t.Add(day * time.Duration(i-1))

			date = now.New(t).EndOfDay()
			return
		}
	}

	// weekday
	for name, wd := range weekdays {
		if strings.HasPrefix(strings.ToLower(name), due) {
			t := time.Now()
			t = t.Add(day * time.Duration((wd-t.Weekday()+7)%7))

			date = now.New(t).EndOfDay()
			return
		}
	}

	// full format
	for _, format := range fullFormats {
		date, err = time.ParseInLocation(format, due, time.Local)
		if err == nil {
			if date.Year() == 0 {
				date = date.AddDate(time.Now().Year(), 0, 0)
			}
			break
		}
	}

	if err != nil {
		err = errors.Errorf("cannot parse date as any of %v", fullFormats[:])
		return
	}

	date = now.New(date).EndOfDay()
	return
}

func formatDueDate(t time.Time, human bool) string {
	today := now.EndOfDay().Truncate(time.Second)

	// relative date
	for i, s := range relative {
		r := today.Add(day * time.Duration(i-1))
		if t.Equal(r) {
			return s
		}
	}

	// weekday
	if now := time.Now(); t.After(now) && t.Before(now.Add(day*7)) {
		return t.Format(weekday)
	}

	// date with year
	if t.Year() != today.Year() {
		if human {
			return t.Format(fullerHuman)
		}

		return t.Format(fuller)
	}

	// date without year
	if human {
		return t.Format(fullHuman)
	}
	return t.Format(full)
}
