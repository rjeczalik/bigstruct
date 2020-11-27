package cti

import (
	"fmt"
)

var DefaultAuxs = []Aux{{
	ID:     AttrFile,
	Name:   "file",
	Expand: expand.File,
}}

type Aux struct {
	ID     Attr
	Name   string
	Expand Func
}

func (x *Aux) String() string {
	return fmt.Sprintf("%s (%b)", x.Name, x.ID)
}

var expand auxexpand

type auxexpand struct{}
