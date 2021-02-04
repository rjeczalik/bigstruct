package storage

import (
	"fmt"
	"net/url"
)

func Open(uri string) (*Gorm, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "sqlite":
		gdb, err := sqliteOpenOrCreate(u)
		if err != nil {
			return nil, err
		}

		return &Gorm{DB: gdb}, nil
	case "mysql":
		gdb, err := mysqlConnect(u)
		if err != nil {
			return nil, err
		}
		return &Gorm{DB: gdb}, nil
	default:
		return nil, fmt.Errorf("unsupported storage: %q", u.Scheme)
	}
}

func uristring(uri *url.URL) string {
	u := *uri
	u.Scheme = ""
	q := u.Query()
	delete(q, "debug")
	u.RawQuery = q.Encode()
	return u.String()
}

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}
