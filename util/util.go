package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func AppendStrs(s ...string) string {
	builder := strings.Builder{}

	for _, str := range s {
		builder.WriteString(str)
	}
	return builder.String()
}

func StrEquals(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	return a == b
}

func UrlEquals(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if strings.HasPrefix(a, "/") {
		a = a[1:]
	}
	if strings.HasPrefix(b, "/") {
		b = b[1:]
	}
	if strings.HasSuffix(a, "/") {
		a = a[:len(a)-1]
	}
	if strings.HasSuffix(b, "/") {
		b = b[:len(b)-1]
	}
	return a == b
}

func FormatStr(format string, v ...any) string {
	return fmt.Sprintf(format, v...)
}

func JoinStr(sep string, str ...string) string {
	return strings.Join(str, sep)
}

func ToString(v any) string {
	return fmt.Sprintf("%v", v)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func RemoveIfExists(path string) error {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return os.RemoveAll(path)
		} else {
			return os.Remove(path)
		}
	}
	return nil
}

func PrettyJsonMarshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
