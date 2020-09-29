package misc

import "time"

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
const DateTimeFormatJSON string = DateFormatRev + "T" + TimeFormatWithMS + "Z"

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
	DateTimeFormatJSON,
	DateTimeFormatJSONTZ,
	DateTimeFormatShortJSON,
	DateTimeFormatShortJSONTZ,
	DateFormatRev + "T" + TimeFormatWithMS + "-0700",
	DateFormatRev + "T" + TimeFormat + "-0700",
}

// ParseJSONtime --
func ParseJSONtime(s string) (t time.Time, err error) {
	for _, f := range jsonFormats {
		t, err = time.Parse(f, s)
		if err == nil {
			return
		}
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//
