package query

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
)

type IndexFunc func(context.Context, *model.Index) error

type Query struct {
	Storage   *storage.Gorm
	Codec     isr.Codec
	IndexFunc IndexFunc
}

func (q *Query) Get(ctx context.Context, index, key string) (model.Values, model.Schemas, error) {
	var (
		v model.Values
		s model.Schemas
	)

	return v, s, q.Storage.Transaction(q.txGet(ctx, index, key, &v, &s))
}

func (q *Query) Set(ctx context.Context, index, namespace string, o isr.Object) error {
	return q.Storage.Transaction(q.txSet(ctx, index, namespace, o))
}

func (q *Query) txSet(ctx context.Context, index, namespace string, o isr.Object) storage.Func {
	return func(tx storage.Gorm) error {
		var (
			idx = &model.Index{Name: index}
		)

		if i := strings.IndexRune(index, '='); i != -1 {
			idx.Name = index[:i]
			idx.Property = index[i+1:]
		}

		n, err := tx.Namespace(namespace)
		if err != nil {
			return err
		}

		ns, err := q.buildNamespaces(ctx, tx, idx, n)
		if err != nil {
			return err
		}

		sbase, err := q.buildSchemas(ctx, tx, ns, isr.Prefix)
		if err != nil {
			return err
		}

		var (
			schema = sbase.Fields().Object()
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
				return fmt.Errorf("the key %q (%#v) does not exist in schema", key, n.Value)
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
			v = model.MakeValues(n, f)
			s = model.MakeSchemas(n, f)
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

func (q *Query) txGet(ctx context.Context, index, key string, pv *model.Values, ps *model.Schemas) storage.Func {
	return func(tx storage.Gorm) error {
		var (
			idx = &model.Index{Name: index}
		)

		if i := strings.IndexRune(index, '='); i != -1 {
			idx.Name = index[:i]
			idx.Property = index[i+1:]
		}

		ns, err := q.buildNamespaces(ctx, tx, idx, nil)
		if err != nil {
			return err
		}

		v, err := q.buildValues(ctx, tx, ns, key)
		if err != nil {
			return err
		}

		s, err := q.buildSchemas(ctx, tx, ns, key)
		if err != nil {
			return err
		}

		*pv, *ps = v, s

		return nil
	}
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
	if err := q.indexFunc(ctx, g, idx); err != nil {
		return 0, err
	}

	var (
		m = idx.ValueIndex.Get()
		n = len(ns)
	)

	for i, ns := range ns {
		var prop interface{}

		if v, ok := m[ns.Name]; ok && v != "" {
			prop = types.MakeYAML(v).Value()
		}

		if err := ns.Property.Set(prop); err != nil {
			return 0, fmt.Errorf(
				"unable to set property %v for namespace %q indexed via %q: %w",
				prop, ns.Name, idx.Prefix(), err,
			)
		}

		if last != nil && last.Name == ns.Name {
			n = i + 1
			break
		}
	}

	return n, nil
}

func (q *Query) codec() isr.Codec {
	if q.Codec != nil {
		return q.Codec
	}
	return codec.Default
}

func (q *Query) buildValues(ctx context.Context, g storage.Gorm, ns model.Namespaces, key string) (model.Values, error) {
	var all model.Values

	for _, ns := range ns {
		v, err := g.ListValues(ns, key)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		all = append(all, v...)
	}

	return all, nil
}

func (q *Query) buildSchemas(ctx context.Context, g storage.Gorm, ns model.Namespaces, key string) (model.Schemas, error) {
	var all model.Schemas

	for _, ns := range ns {
		s, err := g.ListSchemas(ns, key)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		all = append(all, s...)
	}

	return all, nil
}

func (q *Query) indexFunc(ctx context.Context, g storage.Gorm, idx *model.Index) error {
	if q.IndexFunc != nil {
		return q.IndexFunc(ctx, idx)
	}

	return g.DB.Where("`name` = ? AND `property` = ?", idx.Name, idx.Property).First(&idx).Error
}
