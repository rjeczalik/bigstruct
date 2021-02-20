package codec

import (
	"bytes"
	"errors"
	"path"
	"text/template"

	"github.com/rjeczalik/bigstruct/big"
)

type Template struct {
	Data  interface{}
	Funcs map[string]interface{}
}

var _ big.Codec = (*Template)(nil)

func (t Template) Encode(key string, o big.Struct) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return nil // nothing to template, skip
	}

	tmpl, err := template.New(key).Funcs(t.Funcs).Parse(string(p))
	if err != nil {
		return &big.Error{
			Type: "template",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	var (
		buf bytes.Buffer
	)

	if err := tmpl.Execute(&buf, t.Data); err != nil {
		return &big.Error{
			Type: "template",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = buf.String()
	o[k] = n

	return nil
}

func (t Template) Decode(key string, o big.Struct) error {
	return errors.New("codec: template does not support decoding")
}

func (Template) GoString() string {
	return "codec.Template{}"
}
