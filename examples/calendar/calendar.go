/*
Based on https://github.com/ericniebler/range-v3/blob/master/example/calendar.cpp
*/
package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)
import "github.com/tmr232/go-explore/itertools"

func datesFrom(t time.Time) itertools.Iterator[time.Time] {
	return itertools.FromCallable(func() time.Time {
		current := t
		t = t.AddDate(0, 0, 1)
		return current
	})
}

func dates(startYear, stopYear int) itertools.Iterator[time.Time] {
	return itertools.TakeWhile(
		func(t time.Time) bool { return t.Year() < stopYear },
		datesFrom(time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)),
	)
}

func datesOfYear(year int) itertools.Iterator[time.Time] {
	return itertools.TakeWhile(
		func(t time.Time) bool { return t.Year() == year },
		datesFrom(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)),
	)
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
		builder.WriteString(strings.Repeat(" ", int(week[0].Weekday())*3))
		for _, day := range week {
			builder.WriteString(fmt.Sprintf("%3d", day.Day()))
		}
		builder.WriteString(strings.Repeat(" ", (7-int(week[len(week)-1].Weekday()+1))*3))
		builder.WriteString(" ")
		return builder.String()
	}
	return itertools.Map(formatWeek, byWeek(dates))
}

func byMonth(dates itertools.Iterator[time.Time]) itertools.Iterator[[]time.Time] {
	return itertools.ChunkBy(dates, time.Time.Month)
}

func monthTitle(date time.Time) string {
	name := date.Month().String()
	width := 22
	return fmt.Sprintf("%*s", -width, fmt.Sprintf("%*s", (width+len(name))/2, name))
}

func layoutMonths(months itertools.Iterator[[]time.Time]) itertools.Iterator[itertools.Iterator[string]] {
	return itertools.Map(
		func(month []time.Time) itertools.Iterator[string] {
			return itertools.Chain(
				itertools.Literal(monthTitle(month[0])),
				formatWeeks(itertools.FromSlice(month)),
			)
		},
		months,
	)
}

func main() {
	var startYear int
	var stop int
	var perLine int

	flag.IntVar(&startYear, "start", time.Now().Year(), "Year to start")
	flag.IntVar(&stop, "stop", time.Now().Year()+1, "Year to stop")
	flag.IntVar(&perLine, "per-line", 3, "Number months per line")

	flag.Parse()

	printCalendar(startYear, stop, perLine)
}

func printCalendar(startYear, stopYear, perLine int) {
	months := layoutMonths(byMonth(dates(startYear, stopYear)))
	lines := itertools.Chunked(months, perLine)
	for lines.Next() {
		interleaved := itertools.InterleaveLongest(strings.Repeat(" ", 22), lines.Value()...)
		for interleaved.Next() {
			for _, section := range interleaved.Value() {
				fmt.Print(section)
			}
			fmt.Println()
		}
	}
}
