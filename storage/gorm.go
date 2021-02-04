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

	return cfg
}

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

func (g Gorm) With(tx *gorm.DB) Gorm {
	return Gorm{DB: tx}
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
	var (
		cur model.Index
		tx  = g.DB.Begin()
	)

	switch err := tx.Where("`name` = ? AND `property` = ?", i.Name, i.Property).First(&cur).Error; {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := tx.Create(i).Error; err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit().Error
	case err != nil:
		_ = tx.Rollback()
		return err
	}

	var (
		db     = tx.Table("index").Where("`id` = ?", cur.ID)
		update = new(model.Index)
	)

	if old := cur.ValueIndex; !cur.ValueIndex.Merge(i.ValueIndex.Get()).Equal(old) {
		update.ValueIndex = cur.ValueIndex
	}

	if old := cur.SchemaIndex; !cur.SchemaIndex.Merge(i.SchemaIndex.Get()).Equal(old) {
		update.SchemaIndex = cur.SchemaIndex
	}

	if err := db.Updates(update).Error; err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (g Gorm) UpsertValues(v model.Values) error {
	tx := g.DB.Session(&gorm.Session{
		CreateBatchSize: len(v),
		PrepareStmt:     true,
	}).Begin()

	for _, v := range v {
		q := &model.Value{
			Key:         v.Key,
			NamespaceID: v.NamespaceID,
		}

		err := tx.Model(q).Where(q).Select("id").Take(&v.ID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(v).Error

	if err != nil {
		_ = tx.Rollback().Error

		return err
	}

	return tx.Commit().Error
}

func (g Gorm) UpsertSchemas(s model.Schemas) error {
	tx := g.DB.Session(&gorm.Session{
		CreateBatchSize: len(s),
		PrepareStmt:     true,
	}).Begin()

	for _, s := range s {
		q := &model.Schema{
			Key:         s.Key,
			NamespaceID: s.NamespaceID,
		}

		err := tx.Model(q).Where(q).Select("id").Take(&s.ID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "encoding", "schema"}),
	}).Create(s).Error

	if err != nil {
		_ = tx.Rollback().Error

		return err
	}

	return tx.Commit().Error
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

func (g Gorm) ListSchemas(namespaceID uint64, key string) (model.Schemas, error) {
	var (
		s  model.Schemas
		db = g.DB
	)

	if namespaceID != 0 {
		db = db.Where("`namespace_id` = ?", namespaceID)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

	return s, db.
		Preload("Namespace").
		Order("`key` ASC").
		Find(&s).
		Error
}

func (g Gorm) ListValues(namespaceID uint64, key string) (model.Values, error) {
	var (
		v  model.Values
		db = g.DB
	)

	if namespaceID != 0 {
		db = db.Where("`namespace_id` = ?", namespaceID)
	}

	if key != "" {
		db = db.Where("`key` LIKE ?", key+"%")
	}

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
