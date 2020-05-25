/*
Package misc implements a differents trivial functions
*/
package misc

import (
	"bytes"
	"compress/gzip"
	"io"
)

//----------------------------------------------------------------------------------------------------------------------------//

// GzipPack --  data will be truncated by ReadFrom!
func GzipPack(data io.Reader) (b *bytes.Buffer, err error) {
	tmp := new(bytes.Buffer)

	_, err = tmp.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	b = new(bytes.Buffer)
	w := gzip.NewWriter(b)

	_, err = w.Write(tmp.Bytes())
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// GzipUnpack -- data will be truncated by ReadFrom!
func GzipUnpack(data io.Reader) (b *bytes.Buffer, err error) {
	r, err := gzip.NewReader(data)

	if r != nil {
		defer r.Close()
	}

	if err != nil {
		return nil, err
	}

	b = new(bytes.Buffer)

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// GzipRepack -- data will be truncated by ReadFrom!
func GzipRepack(data io.Reader) (b *bytes.Buffer, err error) {
	u, err := GzipUnpack(data)
	if err != nil {
		return nil, err
	}

	return GzipPack(u)
}

//----------------------------------------------------------------------------------------------------------------------------//
