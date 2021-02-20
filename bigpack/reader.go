package bigpack

import (
	"fmt"
	"io"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/codec"
	"github.com/rjeczalik/bigstruct/bigpack/pak"
)

var globalReader = &Reader{}

type Reader struct {
	Codec big.Codec // or codec.Default if nil
}

func (r *Reader) Read(fs pak.FS) (*pak.Pak, error) {
	var p pak.Pak

	err := fs(func(path string, rc io.ReadCloser) error {
		fmt.Println(path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// post-process

	return &p, nil
}

func (r *Reader) codec() big.Codec {
	if r.Codec != nil {
		return r.Codec
	}
	return codec.Default
}

func Read(fs pak.FS) (*pak.Pak, error) {
	return globalReader.Read(fs)
}
