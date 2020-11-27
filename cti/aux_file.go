package cti

import (
	"os"
	"os/user"
	"path"
	"strconv"
	"syscall"
)

type File struct {
	Mode  os.FileMode `json:"m,omitempty" yaml:"m,omitempty"`
	User  string      `json:"u,omitempty" yaml:"u,omitempty"`
	Group string      `json:"g,omitempty" yaml:"g,omitempty"`
}

func (auxexpand) File(key string, o Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	switch x := n.Aux.(type) {
	case File:
		// ok
	default:
		var f File

		if err := reencode(x, &f); err != nil {
			panic("unexpected error: " + err.Error())
		}

		n.Aux = f
		o[k] = n
	}

	return true
}

func newFile(fi os.FileInfo) File {
	f := File{
		Mode: fi.Mode(),
	}

	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		if u, err := user.LookupId(strconv.Itoa(int(stat.Uid))); err == nil {
			f.User = u.Username
		}

		if g, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid))); err == nil {
			f.Group = g.Name
		}
	}

	return f
}
