package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"reflect"
	"strconv"
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
	InterfaceMap map[string]any

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
func Iface2Float(x any) (v float64, err error) {
	vv := reflect.ValueOf(x)

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64:
		v = vv.Float()
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = float64(vv.Int())
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = float64(vv.Uint())
		return

	case reflect.String:
		v, err = strconv.ParseFloat(vv.String(), 64)
		return

	case reflect.Slice:
		var s string
		s, err = bs2String(vv)
		if err != nil {
			return
		}
		v, err = strconv.ParseFloat(s, 64)
		return

	case reflect.Bool:
		if vv.Bool() {
			v = 1.
		} else {
			v = 0.
		}
		return

	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Int --
func Iface2Int(x any) (v int64, err error) {
	vv := reflect.ValueOf(x)

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64:
		v = int64(vv.Float())
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = vv.Int()
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = int64(vv.Uint())
		return

	case reflect.String:
		v, err = strconv.ParseInt(vv.String(), 10, 64)
		return

	case reflect.Slice:
		var s string
		s, err = bs2String(vv)
		if err != nil {
			return
		}
		v, err = strconv.ParseInt(s, 10, 64)
		return

	case reflect.Bool:
		if vv.Bool() {
			v = 1
		} else {
			v = 0
		}
		return

	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Uint --
func Iface2Uint(x any) (v uint64, err error) {
	vv := reflect.ValueOf(x)

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64:
		sv := vv.Float()
		if sv < 0 {
			err = fmt.Errorf("negative value %v", sv)
			return
		}
		v = uint64(sv)
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sv := vv.Int()
		if sv < 0 {
			err = fmt.Errorf("negative value %v", sv)
			return
		}
		v = uint64(sv)
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = vv.Uint()
		return

	case reflect.String:
		v, err = strconv.ParseUint(vv.String(), 10, 64)
		return

	case reflect.Slice:
		var s string
		s, err = bs2String(vv)
		if err != nil {
			return
		}
		v, err = strconv.ParseUint(s, 10, 64)
		return

	case reflect.Bool:
		if vv.Bool() {
			v = 1
		} else {
			v = 0
		}
		return

	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2String --
func Iface2String(x any) (v string, err error) {
	vv := reflect.ValueOf(x)

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64:
		f := vv.Float()
		// is math.Trunc(f) == f correct??
		if math.Trunc(f) == f && f >= float64(math.MinInt64) && f <= float64(math.MaxInt64) {
			v = strconv.FormatInt(int64(f), 10)
			return
		}
		v = strconv.FormatFloat(vv.Float(), 'g', -12, 64)
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(vv.Int(), 10)
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = strconv.FormatUint(vv.Uint(), 10)
		return

	case reflect.String:
		v = vv.String()
		return

	case reflect.Slice:
		var s string
		s, err = bs2String(vv)
		if err != nil {
			return
		}
		v = s
		return

	case reflect.Bool:
		v = strconv.FormatBool(vv.Bool())
		return

	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Bool --
func Iface2Bool(x any) (v bool, err error) {
	vv := reflect.ValueOf(x)

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Float32, reflect.Float64:
		v = int(vv.Float()) != 0
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = vv.Int() != 0
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = vv.Uint() != 0
		return

	case reflect.String:
		v, err = strconv.ParseBool(vv.String())
		return

	case reflect.Slice:
		var s string
		s, err = bs2String(vv)
		if err != nil {
			return
		}

		v, err = strconv.ParseBool(s)
		return

	case reflect.Bool:
		v = vv.Bool()
		return

	default:
		err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
		return
	}
}

// Iface2Time --
func Iface2Time(x any) (v time.Time, err error) {
	switch x := x.(type) {
	case time.Time:
		v = x
		return

	case string:
		v, err = ParseJSONtime(x)
		return

	case []byte:
		v, err = ParseJSONtime(UnsafeByteSlice2String(x))
		return

	default:
		var ts int64
		ts, err = Iface2Int(x)
		if err != nil {
			err = fmt.Errorf(`illegal type of the "%#v" - "%T", expected "%T"`, x, x, v)
			return
		}
		v = time.Unix(ts/int64(time.Second), ts%int64(time.Second))
		return
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func Iface2IfacePtr(src any, dstPtr any) (err error) {
	v := reflect.ValueOf(dstPtr)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf(`"%v" is not a pointer`, dstPtr)
	}

	e := v.Elem()
	var vv any

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

	case reflect.Struct:
		switch dv := dstPtr.(type) {
		case *time.Time:
			*dv, err = Iface2Time(src)
			if err != nil {
				return
			}

		}

	default:
		err = fmt.Errorf(`unsupported kind "%s"`, e.Kind())
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func BaseType(srcT reflect.Type) (t reflect.Type) {
	t = srcT
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool:
		t = reflect.TypeOf(false)

	case reflect.Int:
		t = reflect.TypeOf(int(0))
	case reflect.Int8:
		t = reflect.TypeOf(int8(0))
	case reflect.Int16:
		t = reflect.TypeOf(int16(0))
	case reflect.Int32:
		t = reflect.TypeOf(int32(0))
	case reflect.Int64:
		t = reflect.TypeOf(int64(0))

	case reflect.Uint:
		t = reflect.TypeOf(uint(0))
	case reflect.Uint8:
		t = reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		t = reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		t = reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		t = reflect.TypeOf(uint64(0))

	case reflect.Float32:
		t = reflect.TypeOf(float32(0))
	case reflect.Float64:
		t = reflect.TypeOf(float64(0))

	case reflect.String:
		t = reflect.TypeOf("")
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// MarshalBin --
func MarshalBin(src any) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err = encoder.Encode(src)
	return
}

// UnmarshalBin --
func UnmarshalBin(buf *bytes.Buffer, dst any) (err error) {
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(dst)
}

//----------------------------------------------------------------------------------------------------------------------------//

// v.Kind() is a slice - already checked
func bs2String(v reflect.Value) (s string, err error) {
	k := v.Type().Elem().Kind()
	if k != reflect.Uint8 {
		err = fmt.Errorf(`slice element %s is not %s`, k, reflect.Uint8)
		return
	}

	return string(v.Bytes()), nil
}

//----------------------------------------------------------------------------------------------------------------------------//
