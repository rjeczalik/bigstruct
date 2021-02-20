package storage

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sqliteFile(uri *url.URL) string {
	if uri.Opaque != "" {
		return strings.TrimPrefix(uri.Opaque, "file:")
	}
	return strings.TrimPrefix(uri.Host, "file:")
}

func sqliteURI(uri *url.URL) string {
	u := *uri
	u.Scheme = ""
	if u.Opaque == "" {
		u.Opaque = u.Host
		u.Host = ""
	}
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
		config  = newConfig(uri.Query())
		dialect = sqlite.Open(uristring(uri))
	)

	return gorm.Open(dialect, config)
}

func sqliteCreate(uri *url.URL) (*gorm.DB, error) {
	var (
		path    = sqliteFile(uri)
		dialect = sqlite.Open(sqliteURI(uri))
		config  = newConfig(uri.Query())
	)

	if path != ":memory:" {
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
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
