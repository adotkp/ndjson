// Package ndjson provides encoding and decoding utilities for working with
// newline-delimeted JSON (ndjson).
package ndjson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

// Decoder reads and decodes NDJSON from the input stream.
type Decoder[T any] struct {
	// If true, skip records that do not successfully decode.
	SkipErrors bool

	scanner *bufio.Scanner
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder[T any](r io.Reader) *Decoder[T] {
	return &Decoder[T]{
		SkipErrors: false,
		scanner:    bufio.NewScanner(r),
	}
}

func (d *Decoder[T]) maybeDecode(item *T) (bool, error) {
	if !d.scanner.Scan() {
		err := d.scanner.Err()
		if err != nil {
			return false, err
		}
		return false, io.EOF
	}
	if err := json.Unmarshal(d.scanner.Bytes(), item); err != nil {
		if d.SkipErrors {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DecodeNext returns the next decoded value from the input stream.
// If there are no more records, returns io.EOF.
func (d *Decoder[T]) DecodeNext() (T, error) {
	for {
		var item T
		if ok, err := d.maybeDecode(&item); err != nil {
			return item, err
		} else if ok {
			return item, nil
		}
	}
}

// DecodeAll decodes and returns all the values from the input stream
// until EOF or error.
func (d *Decoder[T]) DecodeAll() ([]T, error) {
	var ret []T
	for {
		item, err := d.DecodeNext()
		if err == io.EOF {
			return ret, nil
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
}

// Unmarshal parses the NDJSON-encoded data and returns the result.
func Unmarshal[T any](data []byte) ([]T, error) {
	return NewDecoder[T](bytes.NewBuffer(data)).DecodeAll()
}
