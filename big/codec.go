package big

import "context"

type Codec interface {
	Encode(ctx context.Context, key string, s Struct) error
	Decode(ctx context.Context, key string, s Struct) error
}
