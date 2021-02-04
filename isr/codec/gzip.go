package codec

import (
	"bytes"
	stdgzip "compress/gzip"
	"io"
)

var gzip gzipCodec

type gzipCodec struct{}

func (gzipCodec) Marshal(p []byte) ([]byte, error) {
	var (
		buf bytes.Buffer
		w   = stdgzip.NewWriter(&buf)
	)

	if _, err := w.Write(p); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (gzipCodec) Unmarshal(p []byte) ([]byte, error) {
	r, err := stdgzip.NewReader(bytes.NewReader(p))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
