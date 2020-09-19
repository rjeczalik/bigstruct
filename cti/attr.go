package cti

import (
	"path"
	"strings"
)

const (
	AttrDelete Attr = 1 << iota
	_
	_
	_
	_
	_
	_
	_
	AttrJSON
	AttrINI
	AttrFlag
	AttrTOML
	AttrHCL
	AttrYAML
	_
	_
)

const (
	AttrOp  = AttrDelete
	AttrEnc = AttrJSON | AttrINI | AttrFlag | AttrTOML | AttrHCL | AttrYAML
)

var attrs = map[Attr]string{
	AttrDelete: "delete",
	AttrJSON:   "json",
	AttrINI:    "ini",
	AttrFlag:   "flag",
	AttrTOML:   "toml",
	AttrHCL:    "hcl",
	AttrYAML:   "yaml",
}

type Attr uint32

func (a Attr) Has(attr Attr) bool {
	return a&attr == attr
}

func (a Attr) String() string {
	var attr []string

	for i := 0; i < 32; i++ {
		if x := Attr(1 << i); a&x != 0 {
			if s, ok := attrs[x]; ok {
				attr = append(attr, s)
			}
		}
	}

	return strings.Join(attr, "|")
}

var (
	_ Func = Attr(0).Set
	_ Func = Attr(0).Add
	_ Func = Attr(0).Del
)

func (a Attr) Set(key string, o Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	n.Attr = a
	o[k] = n

	return true
}

func (a Attr) Add(key string, o Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	n.Attr |= a
	o[k] = n

	return true
}

func (a Attr) Del(key string, o Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	n.Attr &= ^a
	o[k] = n

	return true
}