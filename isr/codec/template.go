package codec

import (
	"bytes"
	"path"
	"text/template"

	"github.com/rjeczalik/bigstruct/isr"
)

type Template struct {
	Data  interface{}
	Funcs map[string]interface{}
}

var _ isr.Codec = (*Template)(nil)

func (t Template) Encode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		// nothing to template, skip
		return nil
	}

	tmpl, err := template.New(key).Funcs(t.Funcs).Parse(string(p))
	if err != nil {
		return &isr.Error{
			Type: "template",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, t.Data); err != nil {
		return &isr.Error{
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

func (t Template) Decode(key string, o isr.Object) error {
	return nil
}

func (Template) GoString() string {
	return "codec.Template{}"
}
