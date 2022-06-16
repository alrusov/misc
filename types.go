package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	// CtxString --
	CtxString string

	// CtxInt64 --
	CtxInt64 int64

	// CtxUint64 --
	CtxUint64 uint64
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	// InterfaceMap --
	InterfaceMap map[string]interface{}

	// ByteSliceMap --
	ByteSliceMap map[string][]byte

	// StringMap --
	StringMap map[string]string

	// BoolMap --
	BoolMap map[string]bool

	// IntMap --
	IntMap map[string]int

	// Int64Map --
	Int64Map map[string]int64

	// UintMap --
	UintMap map[string]uint

	// Uint64Map --
	Uint64Map map[string]uint64

	// Float64Map --
	Float64Map map[string]float64
)

//----------------------------------------------------------------------------------------------------------------------------//

// GetFloat --
func (m InterfaceMap) GetFloat(name string) (v float64, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2Float(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

// GetInt --
func (m InterfaceMap) GetInt(name string) (v int64, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2Int(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

// GetUint --
func (m InterfaceMap) GetUint(name string) (v uint64, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2Uint(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

// GetString --
func (m InterfaceMap) GetString(name string) (v string, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2String(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

// GetBool --
func (m InterfaceMap) GetBool(name string) (v bool, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2Bool(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

// GetTime --
func (m InterfaceMap) GetTime(name string) (v time.Time, err error) {
	x, exists := m[name]
	if !exists {
		err = fmt.Errorf(`%s: parameter not found`, name)
		return
	}

	v, err = Iface2Time(x)
	if err != nil {
		err = fmt.Errorf("%s: %s", name, err.Error())
		return
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Iface2Float --
func Iface2Float(x interface{}) (v float64, err error) {
	switch x := x.(type) {
	case float32:
		v = float64(x)
		return
	case float64:
		v = x
		return
	case int:
		v = float64(x)
		return
	case int8:
		v = float64(x)
		return
	case int32:
		v = float64(x)
		return
	case int64:
		v = float64(x)
		return
	case uint:
		v = float64(x)
		return
	case uint8:
		v = float64(x)
		return
	case uint32:
		v = float64(x)
		return
	case uint64:
		v = float64(x)
		return
	case string:
		v, err = strconv.ParseFloat(x, 64)
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Int --
func Iface2Int(x interface{}) (v int64, err error) {
	switch x := x.(type) {
	case float32:
		v = int64(x)
		return
	case float64:
		v = int64(x)
		return
	case int:
		v = int64(x)
		return
	case int8:
		v = int64(x)
		return
	case int32:
		v = int64(x)
		return
	case int64:
		v = x
		return
	case uint:
		v = int64(x)
		return
	case uint8:
		v = int64(x)
		return
	case uint32:
		v = int64(x)
		return
	case uint64:
		v = int64(x)
		return
	case string:
		v, err = strconv.ParseInt(x, 10, 64)
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Uint --
func Iface2Uint(x interface{}) (v uint64, err error) {
	switch x := x.(type) {
	case float32:
		xx := x
		if xx < 0 {
			err = fmt.Errorf("negative value: %f", xx)
		}
		v = uint64(xx)
		return
	case float64:
		xx := x
		if xx < 0 {
			err = fmt.Errorf("negative value: %f", xx)
		}
		v = uint64(xx)
		return
	case int:
		v = uint64(x)
		return
	case int8:
		v = uint64(x)
		return
	case int32:
		v = uint64(x)
		return
	case int64:
		v = uint64(x)
		return
	case uint:
		v = uint64(x)
		return
	case uint8:
		v = uint64(x)
		return
	case uint32:
		v = uint64(x)
		return
	case uint64:
		v = x
		return
	case string:
		v, err = strconv.ParseUint(x, 10, 64)
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2String --
func Iface2String(x interface{}) (v string, err error) {
	switch x := x.(type) {
	case float32:
		v = strconv.FormatFloat(float64(x), 'g', 5, 64)
		return
	case float64:
		v = strconv.FormatFloat(x, 'g', 5, 64)
		return
	case int:
		v = strconv.FormatInt(int64(x), 10)
		return
	case int8:
		v = strconv.FormatInt(int64(x), 10)
		return
	case int32:
		v = strconv.FormatInt(int64(x), 10)
		return
	case int64:
		v = strconv.FormatInt(x, 10)
		return
	case string:
		v = x
		return
	case uint:
		v = strconv.FormatUint(uint64(x), 10)
		return
	case uint8:
		v = strconv.FormatUint(uint64(x), 10)
		return
	case uint32:
		v = strconv.FormatUint(uint64(x), 10)
		return
	case uint64:
		v = strconv.FormatUint(x, 10)
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Bool --
func Iface2Bool(x interface{}) (v bool, err error) {
	switch x := x.(type) {
	case bool:
		v = x
		return
	case float32:
		v = int64(x) != 0
		return
	case float64:
		v = int64(x) != 0
		return
	case int:
		v = x != 0
		return
	case int8:
		v = x != 0
		return
	case int32:
		v = x != 0
		return
	case int64:
		v = x != 0
		return
	case uint:
		v = x != 0
		return
	case uint8:
		v = x != 0
		return
	case uint32:
		v = x != 0
		return
	case uint64:
		v = x != 0
		return
	case string:
		v = false
		switch strings.ToLower(x) {
		case "true", "t", "1":
			v = true
		}
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Time --
func Iface2Time(x interface{}) (v time.Time, err error) {
	switch x := x.(type) {
	case time.Time:
		v = x
		return
	case float32:
		v = UnixNano2UTC(int64(x))
		return
	case float64:
		v = UnixNano2UTC(int64(x))
		return
	case int:
		v = UnixNano2UTC(int64(x))
		return
	case int8:
		v = UnixNano2UTC(int64(x))
		return
	case int32:
		v = UnixNano2UTC(int64(x))
		return
	case int64:
		v = UnixNano2UTC(x)
		return
	case uint:
		v = UnixNano2UTC(int64(x))
		return
	case uint8:
		v = UnixNano2UTC(int64(x))
		return
	case uint32:
		v = UnixNano2UTC(int64(x))
		return
	case uint64:
		v = UnixNano2UTC(int64(x))
		return
	case string:
		v, err = ParseJSONtime(x)
		return
	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func Iface2IfacePtr(src interface{}, dstPtr interface{}) (err error) {
	v := reflect.ValueOf(dstPtr)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf(`"%v" is not a pointer`, dstPtr)
	}

	e := v.Elem()
	var vv interface{}

	switch e.Kind() {
	case reflect.Bool:
		vv, err = Iface2Bool(src)
		if err != nil {
			return
		}
		e.SetBool(vv.(bool))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vv, err = Iface2Int(src)
		if err != nil {
			return
		}
		e.SetInt(vv.(int64))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vv, err = Iface2Uint(src)
		if err != nil {
			return
		}
		e.SetUint(vv.(uint64))

	case reflect.Float32, reflect.Float64:
		vv, err = Iface2Float(src)
		if err != nil {
			return
		}
		e.SetFloat(vv.(float64))

	case reflect.String:
		vv, err = Iface2String(src)
		if err != nil {
			return
		}
		e.SetString(vv.(string))

	default:
		err = fmt.Errorf(`unsupported kind "%s"`, e.Kind())
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// MarshalBin --
func MarshalBin(src interface{}) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err = encoder.Encode(src)
	return
}

// UnmarshalBin --
func UnmarshalBin(buf *bytes.Buffer, dst interface{}) (err error) {
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(dst)
}

//----------------------------------------------------------------------------------------------------------------------------//
