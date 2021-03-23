package bigstruct

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
)

type Scope struct {
	Namespace *model.Namespace
	Schema    model.Schemas
	Value     model.Values
}

func NewScope(ns *model.Namespace, s big.Struct) *Scope {
	var (
		f = s.Fields()
	)

	return &Scope{
		Namespace: ns,
		Schema:    model.MakeSchemas(ns, f),
		Value:     model.MakeValues(ns, f),
	}
}

func (s *Scope) Fields() big.Fields {
	if s.Namespace.Meta().Schema {
		return s.fields()
	}
	return s.Value.Fields()
}

func (s *Scope) fields() big.Fields {
	return append(s.Schema.Fields(), s.Value.Fields()...)
}

func (s *Scope) Struct() big.Struct {
	return s.Fields().Struct()
}

func (s *Scope) Encode(ctx context.Context, codec big.Codec) error {
	if s.Namespace.Meta().Schema || len(s.Schema) == 0 || len(s.Value) == 0 {
		return nil // early return
	}

	var (
		o = s.fields().Struct()
		f big.Fields
	)

	if err := o.Encode(ctx, codec); err != nil {
		return err
	}

	f = o.Fields()

	s.Schema = model.MakeSchemas(s.Namespace, f)
	s.Value = model.MakeValues(s.Namespace, f)

	return nil
}

type Object struct {
	Index     *model.Index
	Namespace *model.Namespace
	Scopes    []Scope
}

func (obj *Object) Build(ctx context.Context, s big.Struct, c big.Codec) (*Scope, error) {
	var (
		schema = obj.Schemas().Fields().Struct()
	)

	if err := obj.validateSchema(ctx, s, schema); err != nil {
		return nil, fmt.Errorf("schema validation error: %w", err)
	}

	if err := obj.validateKeys(ctx, s, schema); err != nil {
		return nil, fmt.Errorf("key validation error: %w", err)
	}

	if err := obj.validateValues(ctx, s, schema, c); err != nil {
		return nil, fmt.Errorf("value validation error: %w", err)
	}

	return NewScope(obj.Namespace, s), nil
}

func (obj *Object) validateSchema(ctx context.Context, s, schema big.Struct) error {
	return s.Walk(func(key string, o big.Struct) error {
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
				if obj.Namespace.Meta().Schema {
					return fmt.Errorf("cannot override existing schema %q for %q key with %q type (%#v)", t, key, n.Type, n.Value)
				}
			}
		}

		return nil
	})
}

func (obj *Object) validateKeys(ctx context.Context, s, schema big.Struct) error {
	return s.ForEach(func(key string, o big.Struct) error {
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
}

func (obj *Object) validateValues(ctx context.Context, s, schema big.Struct, c big.Codec) error {
	var (
		scope = obj.scope(obj.Namespace)
		f     = s.Fields()
	)

	scope.Value = append(scope.Value, model.MakeValues(scope.Namespace, f)...)

	if err := scope.Encode(ctx, c); err != nil {
		return fmt.Errorf("error decoding %q namespace: %w", scope.Namespace.Ref(), err)
	}

	if err := schema.Merge(scope.Value.Fields().Struct()).Decode(ctx, c); err != nil {
		return fmt.Errorf("error decoding: %w", err)
	}

	return nil
}

func (obj *Object) scope(ns *model.Namespace) *Scope {
	for i, s := range obj.Scopes {
		if s.Namespace.Ref() == ns.Ref() {
			return &obj.Scopes[i]
		}
	}
	return nil
}

func (obj *Object) LoadSchema(tx storage.Gorm, prefix string) error {
	for i, scope := range obj.Scopes {
		s, err := tx.ListSchemas(scope.Namespace, prefix)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error loading schema for %q namespace: %w", scope.Namespace.Ref(), err)
		}

		obj.Scopes[i].Schema = append(obj.Scopes[i].Schema, s...)
	}

	return nil
}

func (obj *Object) LoadValue(tx storage.Gorm, prefix string) error {
	for i, scope := range obj.Scopes {
		v, err := tx.ListValues(scope.Namespace, prefix)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error loading value for %q namespace: %w", scope.Namespace.Ref(), err)
		}

		obj.Scopes[i].Value = append(obj.Scopes[i].Value, v...)
	}

	return nil
}

func (obj *Object) Encode(ctx context.Context, codec big.Codec) error {
	for i := range obj.Scopes {
		scope := &obj.Scopes[i]

		if err := scope.Encode(ctx, codec); err != nil {
			return fmt.Errorf("error decoding %q namespace: %w", scope.Namespace.Ref(), err)
		}
	}

	return nil
}

func (obj *Object) Schemas() model.Schemas {
	var schema model.Schemas

	for _, s := range obj.Scopes {
		if s.Namespace.Meta().Schema {
			schema = append(schema, s.Schema...)
		}
	}

	return schema
}

func (obj *Object) Values() model.Values {
	var value model.Values

	for _, s := range obj.Scopes {
		value = append(value, s.Value...)
	}

	return value
}

func (obj *Object) Namespaces() model.Namespaces {
	var namespaces model.Namespaces

	for _, s := range obj.Scopes {
		namespaces = append(namespaces, s.Namespace)
	}

	return namespaces
}

func (obj *Object) Fields() big.Fields {
	var f big.Fields

	for _, s := range obj.Scopes {
		f = append(f, s.Fields()...)
	}

	return f
}
