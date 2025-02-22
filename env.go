package misc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------------//

const (
	DefaultEnvFile = ".env"
)

// LoadEnv --
func LoadEnv(fileName string) (e error) {
	if fileName == "" {
		fileName = DefaultEnvFile
	}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		if fileName != DefaultEnvFile {
			return err
		}

		return nil
	}

	n := 0
	msg := ""
	src := ""

	defer func() {
		f.Close()
		if msg != "" {
			e = fmt.Errorf("%s in line %d (%s)", msg, n, strings.TrimSpace(src))
		}
	}()

	fb := bufio.NewReader(f)

	eos := byte('\n')
	k := ""
	v := ""

	for {
		src, err = fb.ReadString(eos)
		if err != nil {
			return nil
		}

		n++

		s := strings.TrimSpace(
			strings.Split(src, "#")[0],
		)
		if s == "" {
			continue
		}

		if k == "" {
			if !strings.Contains(s, "=") {
				continue
			}

			sp := strings.SplitN(s, "=", 2)
			k = strings.TrimSpace(sp[0])
			v = strings.TrimSpace(sp[1])

			if len(k) == 0 {
				msg = "Empty name"
				return
			}

			ln := len(v)

			if ln >= 2 && (v[0] == '"' || v[0] == '\'') {
				if v[ln-1] != v[0] {
					msg = "Unclosed quotes"
					return
				}

				v = v[1 : ln-1]
				ln = len(v)
			}

			if ln > 0 &&
				(v[0] == '(' ||
					(len(v) > 1 && v[0] == '$' && v[1] == '(')) &&
				v[ln-1] != ')' {
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
