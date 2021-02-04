package query

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
)

type Query struct {
	Storage *storage.Gorm
	Codec   isr.Codec
}

func (q *Query) Get(ctx context.Context, index, key string) (model.Values, model.Schemas, error) {
	var (
		v model.Values
		s model.Schemas
	)

	return v, s, q.Storage.Transaction(q.txGet(ctx, index, key, &v, &s))
}

func (q *Query) txGet(ctx context.Context, index, key string, pv *model.Values, ps *model.Schemas) storage.Func {
	return func(tx storage.Gorm) error {
		var (
			idx = &model.Index{Name: index}
		)

		if i := strings.IndexRune(index, '/'); i != -1 {
			idx.Name = index[:i]
			idx.Property = index[i+1:]
		}

		ns, err := q.buildNamespaces(ctx, tx, idx)
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

func (q *Query) buildNamespaces(ctx context.Context, g storage.Gorm, idx *model.Index) (model.Namespaces, error) {
	ns, err := g.ListNamespaces()
	if err != nil {
		return nil, err
	}

	if err := q.indexNamespaces(ctx, g, ns, idx); err != nil {
		return nil, err
	}

	return ns, nil
}

func (q *Query) indexNamespaces(ctx context.Context, g storage.Gorm, ns model.Namespaces, idx *model.Index) error {
	if err := g.DB.Where("`name` = ? AND `property` = ?", idx.Name, idx.Property).First(&idx).Error; err != nil {
		return err
	}

	m := idx.ValueIndex.Get()

	for _, ns := range ns {
		var prop interface{}

		if v, ok := m[ns.Name]; ok && v != "" {
			prop = types.MakeYAML(v).Value()
		}

		if err := ns.Property.Set(prop); err != nil {
			return fmt.Errorf(
				"unable to set property %v for namespace %q indexed via %q: %w",
				prop, ns.Name, idx.Prefix(), err,
			)
		}
	}

	return nil
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
		v, err := g.ListValues(ns.ID, key)
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
		s, err := g.ListSchemas(ns.ID, key)
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
