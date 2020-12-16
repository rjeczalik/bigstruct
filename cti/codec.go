package cti

var DefaultCodec Codec

type Codec interface {
	Encode(key string, o Object) error
	Decode(key string, o Object) error
}
