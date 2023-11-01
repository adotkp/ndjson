package ndjson

import (
	"bytes"
	"strings"
	"testing"
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
