package ndjson

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

const contents = "{\"foo\": \"bar0\"}\n{\"foo\": bar1\"}\n{\"foo\": \"bar2\"}"

func TestDecodeError(t *testing.T) {
	_, err := Unmarshal[record]([]byte(contents))
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.HasPrefix(err.Error(), "invalid character 'b'") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDecodeSkipErrors(t *testing.T) {
	dec := NewDecoder[record](bytes.NewBuffer([]byte(contents)))
	dec.SkipErrors = true

	records, err := dec.DecodeAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(records) != 2 || records[0].Foo != "bar0" || records[1].Foo != "bar2" {
		t.Errorf("unexpected result: %v", records)
	}
}

func TestDecodeIOError(t *testing.T) {
	dec := NewDecoder[record](iotest.ErrReader(errors.New("ioerror")))
	dec.SkipErrors = true
	_, err := dec.DecodeAll()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if err.Error() != "ioerror" {
		t.Errorf("unexpected error: %v", err)
	}
}
