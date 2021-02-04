package storage

import (
	"fmt"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlConnect(uri *url.URL) (*gorm.DB, error) {
	var (
		config  = newConfig(uri.Query())
		dialect = mysql.Open(uristring(uri))
	)

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
