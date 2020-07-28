package misc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------------//

// LoadEnv --
func LoadEnv(fileName string) (e error) {
	if fileName == "" {
		fileName = ".env"
	}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	n := 0
	msg := ""

	defer func() {
		f.Close()
		if msg != "" {
			e = fmt.Errorf("%s in line %d", msg, n)
		}
	}()

	fb := bufio.NewReader(f)

	eos := byte('\n')
	k := ""
	v := ""

	for {
		s, err := fb.ReadString(eos)
		if err != nil {
			return nil
		}

		n++

		s = strings.TrimSpace(
			strings.Split(s, "#")[0],
		)
		if s == "" {
			continue
		}

		if k == "" {
			if strings.Index(s, "=") < 0 {
				msg = "Bad format"
				return
			}

			sp := strings.Split(s, "=")
			k = strings.TrimSpace(sp[0])
			v = strings.TrimSpace(sp[1])

			if len(k) == 0 {
				msg = "Empty name"
				return
			}

			ln := len(v)
			if ln > 0 && v[0] == '(' && v[ln-1] != ')' {
				continue
			}

			os.Setenv(k, v)
			k = ""
			continue
		}

		v += " " + s

		if s[len(s)-1] != ')' {
			continue
		}

		os.Setenv(k, v)
		k = ""
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
