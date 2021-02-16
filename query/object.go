package query

import (
	"errors"
	"fmt"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
)

type Scope struct {
	Namespace *model.Namespace
	Schema    model.Schemas
	Value     model.Values
}

func (s *Scope) Object() isr.Object {
	return append(s.Schema.Fields(), s.Value.Fields()...).Object()
}

type Object struct {
	Index     *model.Index
	Namespace *model.Namespace
	Scopes    []Scope
}

func (obj *Object) LoadSchema(tx storage.Gorm, prefix string) error {
	for i, scope := range obj.Scopes {
		s, err := tx.ListSchemas(scope.Namespace, prefix)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed loading schema for %q namespace: %w", scope.Namespace.Ref(), err)
		}

		obj.Scopes[i].Schema = s
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
			return fmt.Errorf("failed loading value for %q namespace: %w", scope.Namespace.Ref(), err)
		}

		obj.Scopes[i].Value = v
	}

	return nil
}

func (obj *Object) Schemas() model.Schemas {
	var schema model.Schemas

	for _, s := range obj.Scopes {
		schema = append(schema, s.Schema...)
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
