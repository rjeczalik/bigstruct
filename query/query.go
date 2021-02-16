package query

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
)

type (
	IndexFunc func(context.Context, *model.Index) error
	CodecFunc func(context.Context, string, *model.Index) (isr.Codec, error)
)

type Query struct {
	Storage *storage.Gorm
	Codec   isr.Codec

	DynamicIndex IndexFunc // StaticIndex by default
	CustomCodec  CodecFunc // no custom codec support by default
}

func (q *Query) Object(ctx context.Context, index, namespace string) (*Object, error) {
	var obj Object

	return &obj, q.Storage.Transaction(q.txObject(ctx, index, namespace, &obj))
}

func (q *Query) txObject(ctx context.Context, index, namespace string, out *Object) storage.Func {
	return func(tx storage.Gorm) (err error) {
		obj := Object{
			Index: new(model.Index),
		}

		if err = obj.Index.SetRef(index); err != nil {
			return err
		}

		if namespace != "" {
			obj.Namespace, err = tx.Namespace(namespace)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		ns, err := q.buildNamespaces(ctx, tx, obj.Index, obj.Namespace)
		if err != nil {
			return err
		}

		for _, n := range ns {
			obj.Scopes = append(obj.Scopes, Scope{
				Namespace: n,
			})
		}

		*out = obj

		return nil
	}
}

func (q *Query) Get(ctx context.Context, index, namespace, key string) (isr.Object, error) {
	var obj isr.Object
	return obj, q.Storage.Transaction(q.txGet(ctx, index, namespace, key, &obj))
}

func (q *Query) Set(ctx context.Context, index, namespace string, o isr.Object) error {
	return q.Storage.Transaction(q.txSet(ctx, index, namespace, o))
}

func (q *Query) txSet(ctx context.Context, index, namespace string, o isr.Object) storage.Func {
	return func(tx storage.Gorm) (err error) {
		var obj Object

		if err = q.txObject(ctx, index, namespace, &obj)(tx); err != nil {
			return err
		}

		if err = obj.LoadSchema(tx, isr.Prefix); err != nil {
			return err
		}

		var (
			schema = obj.Schemas().Fields().Object()
		)

		// validate schema

		err = o.Walk(func(key string, o isr.Object) error {
			var (
				k = path.Base(key)
				n = o[k]
			)

			if n.Type != "" {
				switch t := schema.TypeAt(key); {
				case t == "":
					// ok
				case t == n.Type:
					n.Type = "" // strip
					o[k] = n
				default:
					return fmt.Errorf("cannot override existing schema %q for %q key with %q type (%#v)", t, key, n.Type, n.Value)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		var (
			value = o.Copy().Raw()
		)

		// validate key

		err = value.ForEach(func(key string, o isr.Object) error {
			var (
				d, k = path.Split(key)
				n    = o[k]
			)

			if n.Value == nil {
				return nil
			}

			if _, ok := schema.At(d)[k]; !ok {
				return fmt.Errorf("the key %q (%T) does not exist in schema", key, n.Value)
			}

			return nil
		})

		if err != nil {
			return err
		}

		// validate value

		if err := schema.Merge(value).Decode(q.codec()); err != nil {
			return err
		}

		var (
			f = o.Fields()
			v = model.MakeValues(obj.Namespace, f)
			s = model.MakeSchemas(obj.Namespace, f)
		)

		if err := tx.UpsertSchemas(s); err != nil {
			return err
		}

		if err := tx.UpsertValues(v); err != nil {
			return err
		}

		return nil
	}
}

func (q *Query) txGet(ctx context.Context, index, namespace, key string, out *isr.Object) storage.Func {
	return func(tx storage.Gorm) (err error) {
		var obj Object

		if err = q.txObject(ctx, index, namespace, &obj)(tx); err != nil {
			return err
		}

		if err = obj.LoadSchema(tx, key); err != nil {
			return err
		}

		if err = obj.LoadValue(tx, key); err != nil {
			return err
		}

		if *out, err = q.buildObject(ctx, &obj); err != nil {
			return err
		}

		return nil
	}
}

func (q *Query) buildObject(ctx context.Context, obj *Object) (isr.Object, error) {
	var fields isr.Fields

	for _, s := range obj.Scopes {
		var (
			m = s.Namespace.Meta()
			o = s.Object()
		)

		if m.CustomCodec != "" {
			codec, err := q.customCodec(ctx, m.CustomCodec, obj.Index)
			if err != nil {
				return nil, fmt.Errorf(
					"failed loading %q codec for %q namespace: %w",
					m.CustomCodec, s.Namespace.Ref(), err,
				)
			}

			if err := o.Encode(codec); err != nil {
				return nil, fmt.Errorf(
					"failed encoding %q namespace values with %q codec: %w",
					s.Namespace.Ref(), m.CustomCodec, err,
				)
			}
		}

		fields = append(fields, o.Fields()...)
	}

	return fields.Object().ShakeTypes(), nil
}

func (q *Query) buildNamespaces(ctx context.Context, g storage.Gorm, idx *model.Index, last *model.Namespace) (model.Namespaces, error) {
	ns, err := g.ListNamespaces()
	if err != nil {
		return nil, err
	}

	n, err := q.indexNamespaces(ctx, g, ns, idx, last)
	if err != nil {
		return nil, err
	}

	return ns[:n], nil
}

func (q *Query) indexNamespaces(ctx context.Context, g storage.Gorm, ns model.Namespaces, idx *model.Index, last *model.Namespace) (int, error) {
	if err := q.index(ctx, g, idx); err != nil {
		return 0, err
	}

	var (
		m = idx.Index.Map()
		n = len(ns)
	)

	for i, ns := range ns {
		var prop string

		if v, ok := m[ns.Name]; ok && v != "" {
			prop = fmt.Sprint(v)
		}

		if err := ns.SetProperty(prop); err != nil {
			return 0, fmt.Errorf(
				"unable to set property %v for namespace %q indexed via %q: %w",
				prop, ns.Name, idx.Ref(), err,
			)
		}

		if last != nil && last.Name == ns.Name {
			n = i + 1
			break
		}
	}

	return n, nil
}

func (q *Query) index(ctx context.Context, tx storage.Gorm, idx *model.Index) error {
	if q.DynamicIndex != nil {
		return q.DynamicIndex(ctx, idx)
	}

	return q.StaticIndex(ctx, tx, idx)
}

func (q *Query) StaticIndex(ctx context.Context, tx storage.Gorm, idx *model.Index) error {
	return tx.DB.Where("`name` = ? AND `property` = ?", idx.Name, idx.Property).First(&idx).Error
}

func (q *Query) customCodec(ctx context.Context, typ string, idx *model.Index) (isr.Codec, error) {
	if q.CustomCodec != nil {
		return q.CustomCodec(ctx, typ, idx)
	}
	return nil, fmt.Errorf("custom codec %q not supported", typ)
}

func (q *Query) codec() isr.Codec {
	if q.Codec != nil {
		return q.Codec
	}
	return codec.Default
}
