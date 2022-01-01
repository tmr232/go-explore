package main

import (
	"fmt"
	"strings"
	"time"
)
import "github.com/tmr232/go-explore/itertools"

type dateIterator itertools.Iterator[time.Time]

func datesFrom(t time.Time) itertools.Iterator[time.Time] {
	return itertools.FromCallable(func() time.Time {
		current := t
		t = t.AddDate(0, 0, 1)
		return current
	})
}

func byWeek(dates itertools.Iterator[time.Time]) itertools.Iterator[[]time.Time] {
	return itertools.ChunkBy(dates, func(t time.Time) int {
		_, week := t.AddDate(0, 0, +1).ISOWeek()
		return week
	})
}

func formatWeeks(dates itertools.Iterator[time.Time]) itertools.Iterator[string] {
	formatWeek := func(week []time.Time) string {
		builder := strings.Builder{}
		builder.WriteString(strings.Repeat(" ", int(week[0].Weekday()*3)))
		for _, day := range week {
			builder.WriteString(fmt.Sprintf("%3d", day.Day()))
		}
		return builder.String()
	}
	return itertools.Map(formatWeek, byWeek(dates))
}

func byMonth(dates itertools.Iterator[time.Time]) itertools.Iterator[[]time.Time] {
	return itertools.ChunkBy(dates, time.Time.Month)
}

func main() {
	dates := itertools.Take(10, datesFrom(time.Now()))
	for dates.Next() {
		fmt.Println(dates.Value())
	}

	weeks := byWeek(itertools.Take(10, datesFrom(time.Now())))
	for weeks.Next() {
		fmt.Println("Week!")
		for i, day := range weeks.Value() {
			fmt.Println(i, day)
		}
	}

	formatterWeeks := formatWeeks(itertools.Take(30, datesFrom(time.Now())))
	for formatterWeeks.Next() {
		fmt.Println(formatterWeeks.Value())
	}
}
