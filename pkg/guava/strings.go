package guava

import (
	"os"
	"strconv"
	"strings"
)

func GetEnv(key string, fallVal string) string {
	ev := os.Getenv(key)
	if ev != "" {
		return ev
	}
	return fallVal
}

func GetStr(find string, fallVal string) string {
	if find != "" {
		return find
	}
	return fallVal
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func ToInt64(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic("convert type error")
	}
	return n
}

func Split(s, sep string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	str := strings.Trim(s, sep)
	return strings.Split(str, sep)
}

func ToInt(v string) int {
	if v == "" {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

func In(val interface{}, slice ...interface{}) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
