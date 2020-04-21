package misc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/alrusov/bufpool"
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

//----------------------------------------------------------------------------------------------------------------------------//

// Iface2Float --
func Iface2Float(x interface{}) (v float64, err error) {
	switch x.(type) {
	case float32:
		v = float64(x.(float32))
		return
	case float64:
		v = x.(float64)
		return
	case int:
		v = float64(x.(int))
		return
	case int32:
		v = float64(x.(int32))
		return
	case int64:
		v = float64(x.(int64))
		return
	case string:
		v, err = strconv.ParseFloat(x.(string), 64)
		return
	default:
		err = fmt.Errorf(`Illegal type of "%#v" - "%T", expected "%T"`, x, x, float64(0))
		return
	}
}

// Iface2Int --
func Iface2Int(x interface{}) (v int64, err error) {
	switch x.(type) {
	case float32:
		v = int64(x.(float32))
		return
	case float64:
		v = int64(x.(float64))
		return
	case int:
		v = int64(x.(int))
		return
	case int32:
		v = int64(x.(int32))
		return
	case int64:
		v = x.(int64)
		return
	case string:
		v, err = strconv.ParseInt(x.(string), 10, 64)
		return
	default:
		err = fmt.Errorf(`Illegal type of "%#v" - "%T", expected "%T"`, x, x, int64(0))
		return
	}
}

// Iface2Uint --
func Iface2Uint(x interface{}) (v uint64, err error) {
	switch x.(type) {
	case float32:
		xx := x.(float32)
		if xx < 0 {
			err = fmt.Errorf("Negative value: %f", xx)
		}
		v = uint64(xx)
		return
	case float64:
		xx := x.(float64)
		if xx < 0 {
			err = fmt.Errorf("Negative value: %f", xx)
		}
		v = uint64(xx)
		return
	case int:
		v = uint64(x.(int))
		return
	case int32:
		v = uint64(x.(int32))
		return
	case int64:
		v = x.(uint64)
		return
	case string:
		v, err = strconv.ParseUint(x.(string), 10, 64)
		return
	default:
		err = fmt.Errorf(`Illegal type of "%#v" - "%T", expected "%T"`, x, x, int64(0))
		return
	}
}

// Iface2String --
func Iface2String(x interface{}) (v string, err error) {
	switch x.(type) {
	case float32:
		v = strconv.FormatFloat(float64(x.(float32)), 'g', 5, 64)
		return
	case float64:
		v = strconv.FormatFloat(x.(float64), 'g', 5, 64)
		return
	case int:
		v = strconv.FormatInt(int64(x.(int)), 10)
		return
	case int32:
		v = strconv.FormatInt(int64(x.(int32)), 10)
		return
	case int64:
		v = strconv.FormatInt(x.(int64), 10)
		return
	case string:
		v = x.(string)
		return
	default:
		err = fmt.Errorf(`Illegal type of "%#v" - "%T", expected "%T"`, x, x, "")
		return
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

// MarshalBin --
// Don't forget call bufpool.PutBuf(buf) in the calling function!
func MarshalBin(src interface{}) (buf *bytes.Buffer, err error) {
	buf = bufpool.GetBuf()
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
