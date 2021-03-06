package misc

import (
	"bytes"
	"reflect"
	"runtime"
	"testing"
	"time"
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
		{"^http://localhost:1234", "^http://localhost:1234"},
		{"^://///http://localhost:1234", "^://http://localhost:1234"},
		{"qqqwww/^://///http://localhost:1234", "qqqwww/^://http://localhost:1234"},
		{"qqqwww/////^://///http://localhost:1234", "qqqwww/^://http://localhost:1234"},
		{"^:http://localhost:1234", "^:http://localhost:1234"},
		{"http://localhost:1234", "http://localhost:1234"},
		{"http://localhost:1234/", "http://localhost:1234"},
		{"http://localhost", "http://localhost"},
		{"http://localhost/", "http://localhost"},
		{"http://localhost/////xxx/////yyy/zzz//", "http://localhost/xxx/yyy/zzz"},
		{"http://localhost/xxx/////yyy/zzz//", "http://localhost/xxx/yyy/zzz"},
		{"http:////localhost/////xxx///?u=https:////yyy/zzz//", "http://localhost/xxx/?u=https://yyy/zzz"},
		{"//localhost", "/localhost"},
		{"localhost///", "localhost"},
		{"//localhost///", "/localhost"},
		{"//localhost/", "/localhost"},
		{"/localhost/", "/localhost"},
		{"localhost/", "localhost"},
		{"localhost", "localhost"},
		{"//localhost/////xxx/////yyy/zzz//", "/localhost/xxx/yyy/zzz"},
		{"//localhost/xxx/////yyy/zzz//", "/localhost/xxx/yyy/zzz"},
		{"////localhost/////xxx///?u=https:////yyy/zzz//", "/localhost/xxx/?u=https://yyy/zzz"},
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

func TestGzip(t *testing.T) {
	s := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1234567890qwertiopasdfghjkl;'\zxcvbnm,./1234567890qwertiopasdfghjkl;`

	packed, err := GzipPack(bytes.NewReader([]byte(s)))
	if err != nil {
		t.Error(err)
		return
	}

	saved := bytes.NewBuffer(packed.Bytes())

	unpacked, err := GzipUnpack(packed)
	if err != nil {
		t.Error(err)
		return
	}

	if s != string(unpacked.Bytes()) {
		t.Errorf(`got "%s", expected "%s"`, unpacked, s)
		return
	}

	packed2, err := GzipRepack(saved)
	if err != nil {
		t.Error(err)
		return
	}

	unpacked2, err := GzipUnpack(packed2)
	if err != nil {
		t.Error(err)
		return
	}

	if s != string(unpacked2.Bytes()) {
		t.Errorf(`got "%s", expected "%s"`, unpacked2, s)
		return
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func TestParseJSONtime(t *testing.T) {
	src := []string{
		"2020-09-08T10:06:05.000+03:00",
		"2020-09-08T10:06:05.000+0300",
		"2020-09-08T10:06:05+03:00",
		"2020-09-08T10:06:05+0300",
		"2020-09-08T07:06:05.000Z",
		"2020-09-08T07:06:05.000",
		"2020-09-08T07:06:05Z",
		"2020-09-08T07:06:05",
	}

	expected := time.Date(2020, 9, 8, 7, 6, 5, 0, time.UTC).UnixNano()

	for _, s := range src {
		tt, err := ParseJSONtime(s)
		if err != nil {
			t.Errorf(`"%s": %s`, s, err.Error())
			continue
		}
		if tt.UnixNano() != expected {
			t.Errorf(`"%s": got %d, extected %d`, s, tt.UnixNano(), expected)
			continue
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func TestSha512Hash(t *testing.T) {
	if string(Sha512Hash([]byte("blah-blah-blah-1234567890!"))) != "55a7682312fa6f0ad3053e529d1d061f7c6e941145e3baae70696f2142c5cf6cc22259102b1986a007837e112444488028244a5f17b4d254b258b672e104c002" {
		t.Fatal("Doesn't work")
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
