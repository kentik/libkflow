package api

import "strings"

func NormalizeName(name string) string {
	return strings.Replace(name, ".", "_", -1)
}
