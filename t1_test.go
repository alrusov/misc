package misc

import (
	"reflect"
	"runtime"
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

func TestNormalizeSlashes(t *testing.T) {
	type samples struct {
		in  string
		out string
	}
	smp := []samples{
		{"http://localhost", "http://localhost"},
		{"http://localhost/", "http://localhost"},
		{"http://localhost/////xxx/////yyy/zzz//", "http://localhost/xxx/yyy/zzz"},
		{"http:////localhost/////xxx///?u=https:////yyy/zzz//", "http://localhost/xxx/?u=https://yyy/zzz"},
	}

	for i, p := range smp {
		i++
		out := NormalizeSlashes(p.in)
		if out != p.out {
			t.Errorf(`Case %d failed: in "%s", out "%s", expected: "%s"`, i, p.in, out, p.out)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func TestAbsPath(t *testing.T) {
	type samples struct {
		in  string
		out string
	}
	var smp []samples

	switch runtime.GOOS {
	case "windows":
		smp = []samples{
			{`c:/qqq/www\eee`, `c:\qqq\www\eee`},
			{`qqq/www/eee`, appExecPath + `\qqq\www\eee`},
			{`\qqq\www\eee`, appExecPath + `\qqq\www\eee`},
			{`@qqq\www\eee`, appWorkDir + `\qqq\www\eee`},
		}
	case "linux":
		smp = []samples{
			{`/qqq/www/eee`, `/qqq/www/eee`},
			{`qqq/www/eee`, appExecPath + `/qqq/www/eee`},
			{`@qqq/www/eee`, appWorkDir + `/qqq/www/eee`},
		}
	}

	for i, p := range smp {
		i++
		out, err := AbsPath(p.in)
		if err != nil {
			t.Errorf(`Case %d failed: in "%s", got error "%s"`, i, p.in, err.Error())
			continue
		}
		if out != p.out {
			t.Errorf(`Case %d failed: in "%s", out "%s", expected: "%s"`, i, p.in, out, p.out)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
