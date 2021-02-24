package codec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/objects"
)

type Func func(typ string, m Map) error

type Map map[string]struct {
	Codec    big.Codec
	Priority int
	Children Map
}

var _ big.Codec = Map(nil)

func (m Map) Register(name string, c big.Codec) Map {
	n := m[name]
	n.Codec = c
	m[name] = n

	return m
}

func (m Map) RegisterMap(name string, priority int, cm Map) Map {
	n := m[name]
	n.Codec = cm
	n.Priority = priority
	n.Children = cm
	m[name] = n

	return cm
}

func (m Map) Codec(typ string) big.Codec {
	var (
		parent = m
		dir    = path.Dir(typ)
		key    = path.Base(typ)
	)

	for _, k := range objects.Split(dir) {
		n, ok := parent[k]
		if !ok {
			return nil
		}
		parent = n.Children
	}

	n, ok := parent[key]
	if !ok {
		return nil
	}

	return n.Codec
}

func (m Map) Encode(ctx context.Context, key string, o big.Struct) error {
	var (
		k   = path.Base(key)
		n   = o[k]
		err error
	)

	for i := len(n.Type); i != -1; i = strings.LastIndexByte(n.Type[:i], '/') {
		var (
			typ = n.Type[:i]
		)

		if c := m.Codec(typ); c != nil {
			if e := c.Encode(ctx, key, o); e != nil {
				err = (&big.Error{
					Type: typ,
					Op:   "encode",
					Key:  key,
					Err:  e,
				}).Chain(err)

				continue
			}

			return nil
		}
	}

	if n.Type != "" {
		err = (&big.Error{
			Type: n.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("no suitable codec found"),
		}).Chain(err)
	}

	return err
}

func (m Map) Decode(ctx context.Context, key string, o big.Struct) error {
	var (
		err error
		k   = path.Base(key)
		n   = o[k]
	)

	if n.Value == nil {
		return nil // early skip
	}

	for i := len(n.Type); i > 0; i = strings.LastIndexByte(n.Type[:i], '/') {
		var (
			typ = n.Type[:i]
		)

		if c := m.Codec(typ); c != nil {
			if e := c.Decode(ctx, key, o); e != nil {
				err = (&big.Error{
					Type: typ,
					Op:   "decode",
					Key:  key,
					Err:  e,
				}).Chain(err)

				continue
			}

			return nil
		}
	}

	if n.Type != "" {
		err = (&big.Error{
			Type: n.Type,
			Op:   "decode",
			Key:  key,
			Err:  errors.New("no suitable codec found"),
		}).Chain(err)

		return err
	}

	// Try to guess content type (n.Type) during Decode. It is useful for
	// generating schema from raw data. While this part is pretty accurate,
	// proper design would be to decouple schema building into a separate routine.
	// This part may require refactoring when more complex schema building
	// would be required - e.g. guessing content type using http

	for _, typ := range m.Keys() {
		if e := m[typ].Codec.Decode(ctx, key, o); e != nil {
			err = (&big.Error{
				Type: typ,
				Op:   "decode",
				Key:  key,
				Err:  e,
			}).Chain(err)

			continue
		}

		// fixme: more complex codecs may require more accurate
		// feedback about dencoding result than quessing here
		// (it currently assumes the type has been set on success)
		if n = o[k]; n.Type != "" && n.Type != typ {
			n.Type = path.Join(typ, n.Type)
			o[k] = n
		}

		return nil
	}

	return err
}

func (m Map) Walk(fn Func) error {
	type elm struct {
		parent Map
		key    string
		left   []string
	}

	if len(m) == 0 {
		return nil
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: m, key: "", left: m.Keys()}}
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		key := path.Join(it.key, k)

		if err := fn(key, it.parent); err != nil {
			return err
		}

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if parent := it.parent[k].Children; len(parent) != 0 {
			queue = append(queue, elm{parent: parent, key: key, left: parent.Keys()})
		}
	}

	return nil
}

func (m Map) WriteTab(w io.Writer) (n int64, err error) {
	if m, err := fmt.Fprintln(w, "TYPE\tCODEC"); err != nil {
		return int64(m), err
	}

	err = m.Walk(func(typ string, m Map) error {
		var (
			k = path.Base(typ)
			u = m[k]
		)

		i, err := fmt.Fprintf(w, "%s\t%#v\n",
			typ,
			nonil(u.Codec, "-"),
		)

		n += int64(i)

		return err
	})

	return n, err
}

func (m Map) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := m.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

func (m Map) String() string {
	var buf strings.Builder

	if _, err := m.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

func (m Map) GoString() string {
	var buf strings.Builder

	buf.WriteString("codec.Map{")

	if keys := m.Keys(); len(keys) > 0 {
		fmt.Fprintf(&buf, "%q", keys[0])

		for _, k := range keys[1:] {
			fmt.Fprintf(&buf, ", %q", k)
		}
	}

	buf.WriteRune('}')

	return buf.String()
}

type priokey []struct {
	prio int
	key  string
}

func (set priokey) Less(i, j int) bool {
	if set[i].prio == set[j].prio {
		return set[i].key < set[j].key
	}
	return set[i].key >= set[j].key
}

func (set priokey) Keys() []string {
	keys := make([]string, 0, len(set))

	for _, it := range set {
		keys = append(keys, it.key)
	}

	return keys
}

func (m Map) Keys() []string {
	var (
		set = make(priokey, 0, len(m))
		it  = make(priokey, 1)[0]
	)

	for k := range m {
		it.key = k
		it.prio = m[k].Priority

		set = append(set, it)
	}

	sort.Slice(set, set.Less)

	return set.Keys()
}
