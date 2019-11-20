package misc

import (
	"reflect"
	"testing"
)

//----------------------------------------------------------------------------------------------------------------------------//

func TestIMcoder(t *testing.T) {
	block := InterfaceMap{
		"int64":         int64(12345),
		"int32":         int32(54321),
		"string":        "123456789",
		"float64Slice1": []float64{3.1415926, 2.718281828},
		"float32Slice2": []float32{3.1415926, 2.718281828},
		"intSlice":      []int{1, 2, 3, 4, 5},
		"intSliceEmpty": []int(nil),
		//"intSliceEmpty2": []int{}, // dst will have value "[]int(nil)"" which is equivalent in most cases but reflect.DeepEqual will give an error
	}

	src := []InterfaceMap{
		block,
		block,
		block,
	}
	dst := []InterfaceMap{}

	buf, err := MarshalBin(src)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = UnmarshalBin(buf, &dst)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if !reflect.DeepEqual(src, dst) {
		t.Errorf("src(%#v) != dst(%#v)", src, dst)
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
