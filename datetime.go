package misc

import (
	"fmt"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

// DateFormat -- standard format of the date
const DateFormat string = "02-01-2006"

// DateFormatRev -- reversed format of the date
const DateFormatRev string = "2006-01-02"

// TimeFormat -- format of the time
const TimeFormat string = "15:04:05"

// TimeFormatWithMS -- format of the time with milliseconds
const TimeFormatWithMS string = "15:04:05.000"

// DateTimeFormat -- format of the date and time
const DateTimeFormat string = DateFormat + " " + TimeFormat

// DateTimeFormatRev -- standard format of the date and time with reversed date
const DateTimeFormatRev string = DateFormatRev + " " + TimeFormat

// DateTimeFormatWithMS -- standard format of the date and time with milliseconds
const DateTimeFormatWithMS string = DateFormat + " " + TimeFormatWithMS

// DateTimeFormatRevWithMS -- standard format of the date and time with reversed date and milliseconds
const DateTimeFormatRevWithMS string = DateFormatRev + " " + TimeFormatWithMS

// DateTimeFormatJSON -- JSON format
const DateTimeFormatJSONWithoutZ string = DateFormatRev + "T" + TimeFormatWithMS
const DateTimeFormatJSON string = DateTimeFormatJSONWithoutZ + "Z"

// DateTimeFormatJSONTZ -- JSON format with TZ
const DateTimeFormatJSONTZ string = DateFormatRev + "T" + TimeFormatWithMS + DateTimeFormatTZ

// DateTimeFormatShortJSON -- Short JSON format
const DateTimeFormatShortJSON string = DateFormatRev + "T" + TimeFormat

// DateTimeFormatShortJSONTZ -- Short JSON format with TZ
const DateTimeFormatShortJSONTZ string = DateFormatRev + "T" + TimeFormat + DateTimeFormatTZ

// DateTimeFormatTZ --
const DateTimeFormatTZ = "Z07:00"

//----------------------------------------------------------------------------------------------------------------------------//

var jsonFormats = []string{
	DateTimeFormatJSONWithoutZ,
	DateTimeFormatJSON,
	DateTimeFormatJSONTZ,
	DateTimeFormatShortJSON,
	DateTimeFormatShortJSONTZ,
	DateFormatRev + "T" + TimeFormatWithMS + "-0700",
	DateFormatRev + "T" + TimeFormat + "-0700",
	DateFormatRev,
	DateFormatRev + "Z",
	DateFormatRev + "-0700",
}

// ParseJSONtime --
func ParseJSONtime(s string) (t time.Time, err error) {
	for _, f := range jsonFormats {
		t, err = time.Parse(f, s)
		if err == nil {
			return
		}
	}

	err = fmt.Errorf(`illegal time format`)
	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Time2JSON --
func Time2JSON(t time.Time) string {
	return t.UTC().Format(DateTimeFormatJSON)
}

// Time2JSONtz --
func Time2JSONtz(t time.Time) string {
	return t.Format(DateTimeFormatJSONTZ)
}

// UnixNano2JSON --
func UnixNano2JSON(ts int64) string {
	return Time2JSON(UnixNano2UTC(ts))
}

//----------------------------------------------------------------------------------------------------------------------------//

// UnixNano2UTC --
func UnixNano2UTC(ts int64) time.Time {
	return time.Unix(ts/int64(time.Second), ts%int64(time.Second)).UTC()
}

//----------------------------------------------------------------------------------------------------------------------------//

// NowUTC --
func NowUTC() time.Time {
	return time.Now().UTC()
}

// NowUnix --
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixNano --
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

//----------------------------------------------------------------------------------------------------------------------------//
