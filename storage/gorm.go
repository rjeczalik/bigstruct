package storage

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/rjeczalik/bigstruct/storage/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var tables = []interface{}{
	new(model.Namespace),
	new(model.Index),
	new(model.Value),
	new(model.Schema),
}

func newConfig(v url.Values) *gorm.Config {
	cfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	switch v.Get("debug") {
	case "0":
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	case "1":
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}

	delete(v, "debug")

	return cfg
}

type Func func(tx Gorm) error

type Gorm struct {
	DB *gorm.DB
}

var _ io.Closer = (*Gorm)(nil)

func (g Gorm) Close() error {
	db, err := g.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (g Gorm) Transaction(fn Func) error {
	return g.DB.Transaction(func(tx *gorm.DB) error {
		return fn(Gorm{DB: tx})
	})
}

func (g Gorm) Namespace(namespace string) (*model.Namespace, error) {
	name, prop, err := model.ParseNamespace(namespace)
	if err != nil {
		return nil, err
	}

	var ns model.Namespace

	if err := g.DB.Where("name = ?", name).First(&ns).Error; err != nil {
		return nil, err
	}

	if err := ns.Property.Set(prop); err != nil {
		return nil, fmt.Errorf("unable to set property %v for namespace %q: %w", prop, name, err)
	}

	return &ns, nil
}

func (g Gorm) UpsertNamespace(n *model.Namespace) error {
	return g.DB.
		Where("`name` = ? and `property` = ?", n.Name, n.Property).
		Assign(n).
		FirstOrCreate(n).
		Error
}

func (g Gorm) UpsertIndex(i *model.Index) error {
	return g.Transaction(g.txUpsertIndex(i))
}

func (g Gorm) txUpsertIndex(i *model.Index) Func {
	return func(tx Gorm) error {
		var (
			cur model.Index
		)

		switch err := tx.DB.Where("`name` = ? AND `property` = ?", i.Name, i.Property).First(&cur).Error; {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return tx.DB.Create(i).Error
		case err != nil:
			return err
		}

		var (
			db     = tx.DB.Model(new(model.Index)).Where("`id` = ?", cur.ID)
			update = new(model.Index)
		)

		if old := cur.ValueIndex; !cur.ValueIndex.Merge(i.ValueIndex.Get()).Equal(old) {
			update.ValueIndex = cur.ValueIndex
		}

		if old := cur.SchemaIndex; !cur.SchemaIndex.Merge(i.SchemaIndex.Get()).Equal(old) {
			update.SchemaIndex = cur.SchemaIndex
		}

		return db.Updates(update).Error
	}
}

func (g Gorm) UpsertValues(v model.Values) error {
	if len(v) == 0 {
		return nil
	}
	return g.Transaction(g.txUpsertValues(v))
}

func (g Gorm) txUpsertValues(v model.Values) Func {
	return func(tx Gorm) error {
		for _, v := range v {
			q := &model.Value{
				Key:               v.Key,
				NamespaceID:       v.NamespaceID,
				NamespaceProperty: v.NamespaceProperty,
			}

			err := tx.DB.Model(q).Where(q).Select("id", "value").Take(q).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}
			if v.RawValue == q.RawValue {
				v.ID = q.ID
				continue
			}

			if err := tx.DB.Model(q).Where("`id` = ?", q.ID).Delete(q).Error; err != nil {
				return err
			}
		}

		return tx.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(v).Error
	}
}

func (g Gorm) UpsertSchemas(s model.Schemas) error {
	if len(s) == 0 {
		return nil
	}
	return g.Transaction(g.txUpsertSchemas(s))
}

func (g Gorm) txUpsertSchemas(s model.Schemas) Func {
	return func(tx Gorm) error {
		for _, s := range s {
			q := &model.Schema{
				Key:               s.Key,
				NamespaceID:       s.NamespaceID,
				NamespaceProperty: s.NamespaceProperty,
			}

			err := tx.DB.Model(q).Where(q).Select("id").Take(&s.ID).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}
		}

		return tx.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(s).Error
	}
}

func (g Gorm) ListNamespaces() (model.Namespaces, error) {
	var n model.Namespaces

	return n, g.DB.
		Where("`priority` > -1").
		Order("`priority` ASC").
		Find(&n).Error
}

func (g Gorm) ListIndexes() (model.Indexes, error) {
	var i model.Indexes

	return i, g.DB.
		Order("`name`, `property` ASC").
		Find(&i).
		Error
}

func (g Gorm) ListSchemas(ns *model.Namespace, key string) (model.Schemas, error) {
	var (
		s  model.Schemas
		db = g.DB
	)

	if ns != nil {
		db = db.Where("`namespace_id` = ? AND `namespace_property` = ?", ns.ID, ns.Property)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

	// todo: s.Namespace.Property.Set(s.NamespaceProperty.Value())

	return s, db.
		Preload("Namespace").
		Order("`key` ASC").
		Find(&s).
		Error
}

func (g Gorm) ListValues(ns *model.Namespace, key string) (model.Values, error) {
	var (
		v  model.Values
		db = g.DB
	)

	if ns != nil {
		db = db.Where("`namespace_id` = ? AND `namespace_property` = ?", ns.ID, ns.Property)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

	// todo: v.Namespace.Property.Set(v.NamespaceProperty.Value())

	return v, db.
		Preload("Namespace").
		Order("`key` ASC").
		Find(&v).
		Error
}

func (g Gorm) DeleteNamespace(n *model.Namespace) error {
	return g.DB.
		Where(n).
		Delete((*model.Namespace)(nil)).
		Error
}

func (g Gorm) DeleteIndex(i *model.Index) error {
	return g.DB.
		Where(i).
		Delete((*model.Index)(nil)).
		Error
}
