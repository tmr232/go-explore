package main

import (
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
		builder.WriteString(strings.Repeat(" ", int(week[0].Weekday()+1)*3))
		for _, day := range week {
			builder.WriteString(fmt.Sprintf("%3d", day.Day()))
		}
		builder.WriteString(strings.Repeat(" ", (7-int(week[len(week)-1].Weekday()+1))*3))
		return builder.String()
	}
	return itertools.Map(formatWeek, byWeek(dates))
}

func byMonth(dates itertools.Iterator[time.Time]) itertools.Iterator[[]time.Time] {
	return itertools.ChunkBy(dates, time.Time.Month)
}

func monthTitle(date time.Time) string {
	name := date.Month().String()
	width := 24
	return fmt.Sprintf("%*s", -width, fmt.Sprintf("%*s", (width+len(name))/2, name))
}

func layoutMonths(months itertools.Iterator[[]time.Time]) itertools.Iterator[itertools.Iterator[string]] {
	return itertools.Map(
		func(month []time.Time) itertools.Iterator[string] {
			return itertools.Take(8, itertools.Chain(
				itertools.Literal(monthTitle(month[0])),
				formatWeeks(itertools.FromSlice(month)),
				itertools.Repeat(strings.Repeat(" ", 22)),
			))
		},
		months,
	)
}

func main() {
	months := layoutMonths(byMonth(datesOfYear(2022)))
	lines := itertools.Chunked(months, 3)
	for lines.Next() {
		interleaved := itertools.Interleave(lines.Value()...)
		for interleaved.Next() {
			for _, section := range interleaved.Value() {
				fmt.Print(section)
			}
			fmt.Println()
		}
	}
}
