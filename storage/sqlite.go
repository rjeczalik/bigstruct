package storage

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/rjeczalik/bigstruct/storage/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tables = []interface{}{
	new(model.Namespace),
	new(model.Index),
	new(model.Value),
	new(model.Schema),
	new(model.Meta),
}

func sqliteFile(uri *url.URL) string {
	return strings.TrimPrefix(uri.Opaque, "file:")
}

func sqliteURI(uri *url.URL) string {
	u := *uri
	u.Scheme = ""
	return u.String()
}

func sqliteOpenOrCreate(uri *url.URL) (*gorm.DB, error) {
	if _, err := os.Stat(sqliteFile(uri)); os.IsNotExist(err) {
		return sqliteCreate(uri)
	}

	return sqliteOpen(uri)
}

func sqliteOpen(uri *url.URL) (*gorm.DB, error) {
	var (
		dialect = sqlite.Open(sqliteURI(uri))
		config  = newConfig(uri.Query())
	)

	return gorm.Open(dialect, config)
}

func sqliteCreate(uri *url.URL) (*gorm.DB, error) {
	var (
		path    = sqliteFile(uri)
		dialect = sqlite.Open(sqliteURI(uri))
		config  = newConfig(uri.Query())
	)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(dialect, config)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		if err := gdb.AutoMigrate(table); err != nil {
			return nil, fmt.Errorf("auto-migrate failed for table %T: %w", table, err)
		}
	}

	return gdb, nil
}
