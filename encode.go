package ndjson

import (
	"bytes"
	"encoding/json"
	"io"
)

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

// Encoder writes NDJSON to an output stream.
type Encoder[T any] struct {
	LineSep []byte

	w io.Writer
}

// NewEncoder creates a new Encoder that writes to w.
func NewEncoder[T any](w io.Writer) *Encoder[T] {
	return &Encoder[T]{
		LineSep: []byte{'\n'},
		w:       w,
	}
}

// EncodeNext encodes and writes item to the output stream.
func (e *Encoder[T]) EncodeNext(item T) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	if err := writeAll(e.w, data); err != nil {
		return err
	}
	if err := writeAll(e.w, e.LineSep); err != nil {
		return err
	}
	return nil
}

// EncodeAll encodes and writes all items to the output stream.
func (e *Encoder[T]) EncodeAll(items []T) error {
	for _, item := range items {
		if err := e.EncodeNext(item); err != nil {
			return err
		}
	}
	return nil
}

// Marshal returns the NDJSON encoding of items.
func Marshal[T any](items []T) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder[T](&buf).EncodeAll(items); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
