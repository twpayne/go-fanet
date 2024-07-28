package fanet

import (
	"bytes"
	"encoding/json"
)

type Optional[T any] struct {
	Value T
	Valid bool
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		Value: value,
		Valid: true,
	}
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*o = Optional[T]{}
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*o = Optional[T]{
		Value: value,
		Valid: true,
	}
	return nil
}
