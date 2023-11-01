package ndjson

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestEncodeError(t *testing.T) {
	type badType chan int

	records := []badType{make(chan int)}
	_, err := Marshal[badType](records)
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.HasPrefix(err.Error(), "json: unsupported type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEncodeLineSep(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder[record](&buf)
	enc.LineSep = []byte("-x-")

	err := enc.EncodeAll(makeRecords(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "{\"foo\":\"foo0\",\"x\":0}-x-{\"foo\":\"foo1\",\"x\":1}-x-"
	if expected != buf.String() {
		t.Errorf("unexpected encoding: %s", buf.String())
	}
}

type errWriter struct{}

func (e errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("ioerror")
}

func TestEncodeIOError(t *testing.T) {
	enc := NewEncoder[record](errWriter{})
	err := enc.EncodeAll(makeRecords((1)))
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if err.Error() != "ioerror" {
		t.Errorf("unexpected error: %v", err)
	}
}
