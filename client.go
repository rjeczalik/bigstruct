package bigstruct

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/storage/model"
)

type Op struct {
	Type     string // LIST, GET, SET, DEBUG
	Encode   bool
	Encoding string
	Index    *model.Index
	Overlay  *model.Overlay
	Struct   big.Struct

	Debug struct {
		Values  model.Values
		Schemas model.Schemas
	}
}

type Ref struct {
	Name string
	Prop string
}

func (r Ref) String() string {
	if r.Prop != "" {
		return r.Name + "=" + r.Prop
	}
	return r.Name
}

func (r Ref) IsZero() bool {
	return r.Name == ""
}

type Client struct {
	Transport
}

func (c *Client) List(ctx context.Context, index Ref) (big.Struct, error) {
	op := &Op{
		Type:  "LIST",
		Index: new(model.Index),
	}

	if err := op.Index.SetRef(index.String()); err != nil {
		return nil, err
	}

	if err := c.Do(ctx, op); err != nil {
		return nil, err
	}

	return op.Struct, nil
}

func (c *Client) Get(ctx context.Context, index Ref, key string) (big.Struct, error) {
	op := &Op{
		Type:   "GET",
		Encode: true,
		Index:  new(model.Index),
		Struct: big.Fields{{Key: key}}.Struct(),
	}

	if err := op.Index.SetRef(index.String()); err != nil {
		return nil, err
	}

	if err := c.Do(ctx, op); err != nil {
		return nil, err
	}

	return op.Struct, nil
}

func (c *Client) Diff(ctx context.Context, index Ref, ovStart, ovEnd Ref, key string) (big.Struct, error) {
	var (
		s   = big.Fields{{Key: key}}.Struct()
		ops = []*Op{{
			Type:    "GET",
			Encode:  false,
			Index:   new(model.Index),
			Overlay: new(model.Overlay),
			Struct:  s,
		}, {
			Type:    "GET",
			Encode:  false,
			Index:   new(model.Index),
			Overlay: new(model.Overlay),
			Struct:  s,
		}}
	)

	for _, op := range ops {
		if err := op.Index.SetRef(index.String()); err != nil {
			return nil, err
		}
	}

	for i, ref := range []Ref{ovStart, ovEnd} {
		if err := ops[i].Overlay.SetRef(ref.String()); err != nil {
			return nil, err
		}
	}

	if err := c.Do(ctx, ops...); err != nil {
		return nil, err
	}

	return ops[1].Struct.Sub(ops[0].Struct), nil
}

func (c *Client) Debug(ctx context.Context, index, overlay Ref, key string) (*model.Index, model.Schemas, model.Values, error) {
	op := &Op{
		Type:   "DEBUG",
		Index:  new(model.Index),
		Struct: big.Fields{{Key: key}}.Struct(),
	}

	if err := op.Index.SetRef(index.String()); err != nil {
		return nil, nil, nil, err
	}

	if !overlay.IsZero() {
		op.Overlay = new(model.Overlay)

		if err := op.Overlay.SetRef(overlay.String()); err != nil {
			return nil, nil, nil, err
		}
	}

	if err := c.Do(ctx, op); err != nil {
		return nil, nil, nil, err
	}

	return op.Index, op.Debug.Schemas, op.Debug.Values, nil
}

func (c *Client) Struct(ctx context.Context, index Ref, key string, v interface{}) error {
	op := &Op{
		Type:     "GET",
		Encode:   true,
		Encoding: "struct/json",
		Index:    new(model.Index),
		Struct:   big.Fields{{Key: key}}.Struct(),
	}

	if err := op.Index.SetRef(index.String()); err != nil {
		return err
	}

	if err := c.Do(ctx, op); err != nil {
		return err
	}

	p := op.Struct.Fields().At(0).Bytes()

	if err := json.Unmarshal(p, v); err != nil {
		return fmt.Errorf("failed to unmarshal %q: %w", key, err)
	}

	return nil
}

func (c *Client) Set(ctx context.Context, index, overlay Ref, s big.Struct) error {
	op := &Op{
		Type:    "SET",
		Index:   new(model.Index),
		Overlay: new(model.Overlay),
		Struct:  s,
	}

	if err := op.Index.SetRef(index.String()); err != nil {
		return err
	}

	if err := op.Overlay.SetRef(overlay.String()); err != nil {
		return err
	}

	return c.Do(ctx, op)
}
