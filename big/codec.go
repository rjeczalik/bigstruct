package big

type Codec interface {
	Encode(key string, s Struct) error
	Decode(key string, s Struct) error
}
