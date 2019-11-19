package misc

import (
	"fmt"
	"strconv"
)

//----------------------------------------------------------------------------------------------------------------------------//

// InterfaceMap --
type InterfaceMap map[string]interface{}

// StringMap --
type StringMap map[string]string

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

// GetFloat --
func (p InterfaceMap) GetFloat(name string) (v float64, err error) {
	x, exists := p[name]
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
func (p InterfaceMap) GetInt(name string) (v int64, err error) {
	x, exists := p[name]
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

// GetString --
func (p InterfaceMap) GetString(name string) (v string, err error) {
	x, exists := p[name]
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
