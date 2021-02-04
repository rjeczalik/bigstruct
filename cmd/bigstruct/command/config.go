package command

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/rjeczalik/bigstruct/internal/types"
)

type Endpoint struct {
	URI string `yaml:"uri,omitempty"`
}

func (e Endpoint) Validate() error {
	if e.URI == "" {
		return errors.New(`"uri" is missing or empty`)
	}
	if _, err := url.Parse(e.URI); err != nil {
		return fmt.Errorf(`"uri" is invalid: %w`, err)
	}
	return nil
}

type Config struct {
	Backend Endpoint `yaml:"backend,omitempty"`
}

func (c Config) Validate() error {
	if err := c.Backend.Validate(); err != nil {
		return fmt.Errorf(`"backend" is invalid: %w`, err)
	}
	return nil
}

func (c Config) YAML() types.YAML {
	return types.MakeYAML(&c)
}

func (c *Config) FromYAML(s types.YAML) error {
	if err := s.Unmarshal(c); err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}
