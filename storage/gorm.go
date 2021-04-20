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
	new(model.Overlay),
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

type Scope func(tx Gorm) Gorm

type Gorm struct {
	DB     *gorm.DB
	Scopes []Scope
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
		return fn(g.WithDB(tx))
	})
}

func (g Gorm) WithDB(db *gorm.DB) Gorm {
	return Gorm{
		DB:     db,
		Scopes: g.Scopes,
	}
}

func (g Gorm) WithScopes(scopes ...Scope) Gorm {
	return Gorm{
		DB:     g.DB,
		Scopes: append(g.Scopes, scopes...),
	}
}

func (g Gorm) db() *gorm.DB {
	for _, scope := range g.Scopes {
		g = scope(g)
	}

	return g.DB
}

func (g Gorm) Overlay(overlay string) (*model.Overlay, error) {
	name, prop, err := model.ParseRef(overlay)
	if err != nil {
		return nil, err
	}

	var o model.Overlay

	if err := g.db().Where("name = ?", name).First(&o).Error; err != nil {
		return nil, err
	}

	if err := o.SetProperty(prop); err != nil {
		return nil, fmt.Errorf("unable to set property %q for overlay %q: %w", prop, name, err)
	}

	return &o, nil
}

func (g Gorm) UpsertOverlay(n *model.Overlay) error {
	return g.Transaction(g.txUpsertOverlay(n))
}

func (g Gorm) txUpsertOverlay(n *model.Overlay) Func {
	return func(tx Gorm) error {
		var (
			cur model.Overlay
		)

		switch err := tx.db().Where("`name` = ?", n.Name).First(&cur).Error; {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return tx.DB.Create(n).Error
		case err != nil:
			return err
		}

		var (
			db     = tx.db().Model(&cur).Where("`id` = ?", cur.ID)
			update = &model.Overlay{
				Priority: n.Priority,
			}
		)

		if cur.Metadata.Update(n.Metadata) {
			update.Metadata = cur.Metadata
		}

		return db.Updates(update).Error
	}
}

func (g Gorm) UpsertIndex(i *model.Index) error {
	return g.Transaction(g.txUpsertIndex(i))
}

func (g Gorm) txUpsertIndex(i *model.Index) Func {
	return func(tx Gorm) error {
		var (
			cur model.Index
		)

		switch err := tx.db().Where("`name` = ? AND `property` = ?", i.Name, i.Property).First(&cur).Error; {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return tx.DB.Create(i).Error
		case err != nil:
			return err
		}

		var (
			db     = tx.db().Model(new(model.Index)).Where("`id` = ?", cur.ID)
			update = new(model.Index)
		)

		if cur.Index.Update(i.Index) {
			update.Index = cur.Index
		}

		if cur.Metadata.Update(i.Metadata) {
			update.Metadata = cur.Metadata
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
			if v.Overlay != nil && v.OverlayID == 0 {
				v.OverlayID = v.Overlay.ID
			}

			q := &model.Value{
				Key:             v.Key,
				OverlayID:       v.OverlayID,
				OverlayProperty: v.OverlayProperty,
			}

			err := tx.db().Model(q).Where(q).Select("id", "value", "metadata").Take(q).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}

			if q.Metadata.Update(v.Metadata) {
				v.Metadata = q.Metadata
			}

			if v.RawValue == q.RawValue {
				v.ID = q.ID
				continue
			}

			if err := tx.db().Model(q).Where("`id` = ?", q.ID).Delete(q).Error; err != nil {
				return err
			}
		}

		return tx.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"metadata", "updated_at"}),
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
			if s.Overlay != nil && s.OverlayID == 0 {
				s.OverlayID = s.Overlay.ID
			}

			q := &model.Schema{
				Key:             s.Key,
				OverlayID:       s.OverlayID,
				OverlayProperty: s.OverlayProperty,
			}

			err := tx.db().Model(q).Where(q).Select("id", "type", "encoding", "metadata").Take(q).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}

			if q.Metadata.Update(s.Metadata) {
				s.Metadata = q.Metadata
			}

			if s.Codec() == q.Codec() {
				s.ID = q.ID
				continue
			}

			if err := tx.db().Model(q).Where("`id` = ?", q.ID).Delete(q).Error; err != nil {
				return err
			}
		}

		return tx.db().Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"metadata", "updated_at"}),
		}).Create(s).Error
	}
}

func (g Gorm) ListOverlays() (model.Overlays, error) {
	var n model.Overlays

	return n, g.db().
		Where("`priority` > -1").
		Order("`priority` ASC").
		Find(&n).Error
}

func (g Gorm) ListIndexes() (model.Indexes, error) {
	var i model.Indexes

	return i, g.db().
		Order("`name`, `property` ASC").
		Find(&i).
		Error
}

func (g Gorm) ListSchemas(o *model.Overlay, key string) (model.Schemas, error) {
	var (
		s  model.Schemas
		db = g.db()
	)

	if o != nil {
		db = db.Where("`overlay_id` = ? AND `overlay_property` = ?", o.ID, o.Property)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

	if err := db.Order("`key` ASC").Find(&s).Error; err != nil {
		return nil, err
	}

	s.SetOverlay(o)

	return s, nil
}

func (g Gorm) ListValues(o *model.Overlay, key string) (model.Values, error) {
	var (
		v  model.Values
		db = g.db()
	)

	if o != nil {
		db = db.Where("`overlay_id` = ? AND `overlay_property` = ?", o.ID, o.Property)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

	if err := db.Order("`key` ASC").Find(&v).Error; err != nil {
		return nil, err
	}

	v.SetOverlay(o)

	return v, nil
}

func (g Gorm) DeleteOverlay(n *model.Overlay) error {
	return g.db().
		Where(n).
		Delete((*model.Overlay)(nil)).
		Error
}

func (g Gorm) DeleteIndex(i *model.Index) error {
	return g.db().
		Where(i).
		Delete((*model.Index)(nil)).
		Error
}
