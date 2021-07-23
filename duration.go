package misc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

// empty as a seconds
var (
	re         = `(\d+)(ns|u|us|ms|s|m|h|d|)`
	matchRE    = regexp.MustCompile(`^((?:\s*)` + re + `(?:\s*))+$`)
	intervalRE = regexp.MustCompile(re)

	timeUnits = map[string]int64{
		"ns": int64(time.Nanosecond),
		"u":  int64(time.Microsecond),
		"us": int64(time.Microsecond),
		"ms": int64(time.Millisecond),
		"s":  int64(time.Second),
		"m":  int64(time.Minute),
		"h":  int64(time.Hour),
		"d":  int64(time.Hour * 24),
	}
)

//----------------------------------------------------------------------------------------------------------------------------//

// Interval2Duration --
func Interval2Duration(interval string) (time.Duration, error) {
	t, err := Interval2Int64(interval)
	return time.Duration(t), err
}

// Interval2Int64 --
func Interval2Int64(interval string) (int64, error) {
	interval = strings.TrimSpace(interval)

	if interval == "" {
		return 0, nil
	}

	if !matchRE.MatchString(interval) {
		return 0, fmt.Errorf(`bad interval "%s"`, interval)
	}

	v := intervalRE.FindAllStringSubmatch(strings.TrimSpace(interval), -1)
	if v == nil {
		return 0, fmt.Errorf(`bad interval "%s"`, interval)
	}

	val := int64(0)

	for _, vv := range v {
		n, err := strconv.ParseInt(vv[1], 10, 64)
		if err != nil {
			return 0, err
		}

		if vv[2] == "" {
			vv[2] = "s"
		}
		m, err := TimePrecisionDivider(vv[2], false)
		if err != nil {
			return 0, err
		}

		val += n * m
	}

	return val, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// TimePrecisionDivider --
func TimePrecisionDivider(precision string, upToSecond bool) (int64, error) {
	d, exists := timeUnits[precision]
	if exists && (!upToSecond || d <= int64(time.Second)) {
		return d, nil
	}

	return 1, fmt.Errorf(`unknown precision "%s"`, precision)
}

// CheckTimePrecision --
func CheckTimePrecision(precision string, upToSecond bool) bool {
	_, err := TimePrecisionDivider(precision, upToSecond)
	return err == nil
}

//----------------------------------------------------------------------------------------------------------------------------//
