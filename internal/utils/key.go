package utils

import (
	"strings"
)

//NewKeyName Create consistent key name with given inputs
func NewKeyName(method, url, ip string) string {
	return strings.ToUpper(method + "_" + url + "_" + ip)
}
