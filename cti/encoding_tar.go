package cti

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type TarEncoder struct{}

var _ Encoder = (*TarEncoder)(nil)

func (TarEncoder) String() string { return "tar" }

func (TarEncoder) FileExt() []string { return []string{"tar"} }

func (te TarEncoder) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
		f Fields
	)

	n.Children.ForEach(f.Append)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

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

		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(p)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      f.Key,
				Err:      err,
			}
		}

		if _, err := tw.Write(p); err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      f.Key,
				Err:      err,
			}
		}
	}

	if err := tw.Close(); err != nil {
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

func (te TarEncoder) Decode(key string, o Object) error {
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
		enc = enc[:len(enc)-1]
	}

	var (
		tr = tar.NewReader(bytes.NewReader(p))
		f  Fields
	)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &EncodingError{
				Encoding: te.String(),
				Key:      key,
				Err:      err,
			}
		}

		key := cleanpath(hdr.Name)

		p, err := ioutil.ReadAll(tr)
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
