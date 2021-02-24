package codec

import (
	"bytes"
	"context"
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

func (t Template) Encode(ctx context.Context, key string, s big.Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
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
	s[k] = n

	return nil
}

func (Template) Decode(context.Context, string, big.Struct) error {
	return errors.New("codec: template does not support decoding")
}

func (Template) GoString() string {
	return "codec.Template{}"
}
