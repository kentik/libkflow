package cmetrics

// Taken from https://github.com/influxdata/telegraf/blob/master/plugins/serializers/influx/escape.go

import (
	"strings"
)

const (
	influxEscapes = "\t\n\f\r ,="
)

var (
	influxEscaper = strings.NewReplacer(
		"\t", `\t`,
		"\n", `\n`,
		"\f", `\f`,
		"\r", `\r`,
		`,`, `\,`,
		` `, `\ `,
		`=`, `\=`,
	)
)

// Escape an influx tagkey, tagvalue, or fieldkey
func influxEscape(s string) string {
	if strings.ContainsAny(s, influxEscapes) {
		return influxEscaper.Replace(s)
	}
	return s
}
