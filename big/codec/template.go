package codec

import (
	"bytes"
	"context"
	"path"
	"text/template"

	"github.com/rjeczalik/bigstruct/big"
)

var DefaultTemplate = make(Map)

var DefaultTemplateField = DefaultTemplate.
	RegisterMap("field", 50, make(Map))

var _ = DefaultTemplateField.
	Register("template", Template{})

var _ = DefaultTemplate.
	Register("template", Recursive{
		Codec: Template{},
	})

var templateKey = contextKey{"template-key"}

type Template struct {
	Data  interface{}
	Funcs map[string]interface{}
}

func TemplateWithContext(ctx context.Context, t Template) context.Context {
	return context.WithValue(ctx, templateKey, t)
}

func TemplateFromContext(ctx context.Context) (Template, bool) {
	if t, ok := ctx.Value(templateKey).(Template); ok {
		return t, true
	}

	return Template{}, false
}

var _ big.Codec = (*Template)(nil)

func (t Template) Encode(ctx context.Context, key string, s big.Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
	)

	if n.Value == nil {
		return nil // nothing to template, skip
	}

	if tmpl, ok := TemplateFromContext(ctx); ok {
		if tmpl.Data != nil {
			t.Data = tmpl.Data
		}
		if tmpl.Funcs != nil {
			t.Funcs = tmpl.Funcs
		}
	}

	p, err := tobytes(n.Value)
	if err != nil {
		return nil // nothing to template, skip
	}

	tmpl, err := template.New("encode").Funcs(t.Funcs).Parse(string(p))
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
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = buf.String()
	s[k] = n

	return nil
}

func (Template) Decode(context.Context, string, big.Struct) error {
	return nil // nop
}

func (Template) GoString() string {
	return "codec.Template{}"
}
