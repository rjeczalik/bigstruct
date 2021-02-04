package codec

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/isr"
)

type Tar struct{}

var _ isr.Codec = (*Tar)(nil)

func (Tar) Encode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
		f isr.Fields
	)

	n.Children.ForEach(f.Append)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, f := range f {
		p, err := tobytes(f.Value)
		if err != nil {
			return &isr.Error{
				Type: "tar",
				Op:   "encode",
				Key:  f.Key,
				Err:  err,
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
			return &isr.Error{
				Type: "tar",
				Op:   "encode",
				Key:  f.Key,
				Err:  err,
			}
		}

		if _, err := tw.Write(p); err != nil {
			return &isr.Error{
				Type: "tar",
				Op:   "encode",
				Key:  f.Key,
				Err:  err,
			}
		}
	}

	if err := tw.Close(); err != nil {
		return &isr.Error{
			Type: "tar",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = buf.Bytes()
	n.Children = nil
	o[k] = n

	return nil
}

func (Tar) Decode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &isr.Error{
			Type: "tar",
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	var (
		tr = tar.NewReader(bytes.NewReader(p))
		f  isr.Fields
	)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &isr.Error{
				Type: "tar",
				Op:   "decode",
				Key:  key,
				Err:  err,
			}
		}

		key := cleanpath(hdr.Name)

		p, err := ioutil.ReadAll(tr)
		if err != nil {
			return &isr.Error{
				Type: "tar",
				Op:   "decode",
				Key:  key,
				Err:  err,
			}
		}

		f = append(f, isr.Field{
			Key:   key,
			Value: p,
		})
	}

	n.Value = nil
	n.Children = f.Object()

	o[k] = n

	return nil
}
