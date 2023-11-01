package ndjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type record struct {
	Foo string `json:"foo"`
	X   int    `json:"x"`
}

func makeRecords(n int) []record {
	var ret []record
	for i := 0; i < n; i++ {
		ret = append(ret, record{Foo: fmt.Sprintf("foo%d", i), X: i})
	}
	return ret
}

func makeRecordPtrs(n int) []*record {
	var ret []*record
	for i := 0; i < n; i++ {
		ret = append(ret, &record{Foo: fmt.Sprintf("foo%d", i), X: i})
	}
	return ret
}

func TestRoundTrip(t *testing.T) {
	for i := 0; i <= 5; i++ {
		t.Run(fmt.Sprintf("size%d", i), func(t *testing.T) {
			records := makeRecords(i)
			data, err := Marshal[record](records)
			if err != nil {
				t.Errorf("unexpected encoding error: %v", err)
			}

			recordsOut, err := Unmarshal[record](data)
			if err != nil {
				t.Errorf("unexpected decoding error: %v", err)
			}

			if !reflect.DeepEqual(records, recordsOut) {
				t.Errorf("roundtrip did not match: %v != %v", records, recordsOut)
			}
		})
	}
}

func TestRoundTripPointers(t *testing.T) {
	for i := 0; i <= 5; i++ {
		t.Run(fmt.Sprintf("size%d", i), func(t *testing.T) {
			records := makeRecordPtrs(i)
			data, err := Marshal[*record](records)
			if err != nil {
				t.Errorf("unexpected encoding error: %v", err)
			}

			recordsOut, err := Unmarshal[*record](data)
			if err != nil {
				t.Errorf("unexpected decoding error: %v", err)
			}

			if !reflect.DeepEqual(records, recordsOut) {
				recordStr, err := json.Marshal(records)
				if err != nil {
					t.Fatal(err)
				}
				t.Errorf("roundtrip did not match: %s vs %s", recordStr, data)
			}
		})
	}
}
