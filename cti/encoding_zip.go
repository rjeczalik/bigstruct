package cti

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type ZipEncoder struct{}

var _ Encoder = (*ZipEncoder)(nil)

func (ZipEncoder) String() string { return "zip" }

func (ZipEncoder) FileExt() []string { return []string{"zip"} }

func (te ZipEncoder) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
		f Fields
	)

	n.Children.ForEach(f.Append)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, f := range f {
		p, err := tobytes(f.Value)
		if err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      f.Key,
				Err:      err,
			}
		}

		if len(p) == 0 {
			continue // skip, nothing to encode
		}

		// make the path relative in the archive
		name := filepath.FromSlash(strings.TrimLeft(f.Key, `/\`))

		zf, err := zw.Create(name)
		if err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      f.Key,
				Err:      err,
			}
		}

		if _, err := zf.Write(p); err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      f.Key,
				Err:      err,
			}
		}
	}

	if err := zw.Close(); err != nil {
		return &EncodingError{
			Encoding: te.String(),
			Key:      key,
			Err:      err,
		}
	}

	n.Value = buf.Bytes()
	n.Children = nil
	o[k] = n

	return nil
}

func (te ZipEncoder) Decode(key string, o Object) error {
	var (
		k   = path.Base(key)
		n   = o[k]
		enc Encoding
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &EncodingError{
			Encoding: te.String(),
			Key:      key,
			Err:      err,
		}
	}

	if len(enc) > 1 {
		enc = enc[:len(enc)]
	}

	zr, err := zip.NewReader(bytes.NewReader(p), int64(len(p)))
	if err != nil {
		return &EncodingError{
			Encoding: te.String(),
			Key:      key,
			Err:      err,
		}
	}

	var f Fields

	for _, zf := range zr.File {
		key := cleanpath(zf.Name)

		rc, err := zf.Open()
		if err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      key,
				Err:      err,
			}
		}

		p, err := ioutil.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      key,
				Err:      err,
			}
		}

		f = append(f, Field{
			Key:      key,
			Encoding: enc.Copy(),
			Value:    p,
		})
	}

	n.Value = nil
	n.Children = f.Object()

	o[k] = n

	return nil
}
