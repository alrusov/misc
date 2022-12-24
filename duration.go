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
	re         = `(\d+)(ns|u|us|ms|s|m|h|d|w|)`
	matchRE    = regexp.MustCompile(`^((?:\s*)` + re + `(?:\s*))+$`)
	intervalRE = regexp.MustCompile(re)

	timeUnitsDef = []struct {
		n string
		v int64
	}{
		{"w", int64(time.Hour * 24 * 7)},
		{"d", int64(time.Hour * 24)},
		{"h", int64(time.Hour)},
		{"m", int64(time.Minute)},
		{"s", int64(time.Second)},
		{"ms", int64(time.Millisecond)},
		{"us", int64(time.Microsecond)},
		{"u", int64(time.Microsecond)},
		{"ns", int64(time.Nanosecond)},
	}

	timeUnits map[string]int64
)

func init() {
	timeUnits = make(map[string]int64, len(timeUnitsDef))

	for _, v := range timeUnitsDef {
		timeUnits[v.n] = v.v
	}
}

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

	sign := int64(1)
	if interval[0] == '-' {
		sign = -1
		interval = interval[1:]
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

	return sign * val, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

func Duration2Interval(d time.Duration) string {
	return Int2Interval(int64(d))
}

func Int2Interval(d int64) string {
	p := make([]string, 0, 2*len(timeUnitsDef))

	for _, df := range timeUnitsDef {
		v1 := d / df.v
		v2 := d % df.v

		if v1 != 0 {
			p = append(p, strconv.FormatInt(v1, 10), df.n)
		}

		d = v2
	}

	if len(p) == 0 {
		return "0s"
	}

	return strings.Join(p, "")
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
