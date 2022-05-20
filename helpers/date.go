package helpers

import "time"

const DateFormatDefault = "2006-01-02 15:04:05"

func FormatDefaultDate(t time.Time) string {
	return FormatDate(t, DateFormatDefault)
}

func FormatDate(t time.Time, f string) string {
	return t.Local().Format(f)
}
