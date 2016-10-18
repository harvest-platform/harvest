package json

import (
	"bytes"
	"strings"
	"testing"
)

type record struct {
	ID int `json:"id"`
}

func TestEncoder(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)

	items := []*record{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	for _, item := range items {
		if err := e.Encode(item); err != nil {
			t.Fatal(err)
		}
	}

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}

	if err := e.Encode(&record{}); err == nil {
		t.Error("expected closed error")
	}

	if err := e.Close(); err == nil {
		t.Error("expected closed error")
	}

	// Normalize for comparison.
	act := strings.Replace(b.String(), "\n", "", -1)
	exp := `[{"id":1},{"id":2},{"id":3}]`

	if act != exp {
		t.Errorf("expected %s, got %s", exp, act)
	}
}

func TestDecoder(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(`[
		{"id": 1},
		{"id": 2},
		{"id": 3}
	]`)

	d := NewDecoder(&b)

	exps := []record{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	var i int

	for d.More() {
		exp := exps[i]

		var act record
		if err := d.Decode(&act); err != nil {
			t.Fatal(err)
		}

		if act != exp {
			t.Errorf("expected %v, got %v", exp, act)
		}

		i++
	}

	if i != 3 {
		t.Errorf("expected 3 elements, got %d", i)
	}

	if err := d.Err(); err != nil {
		t.Error(err)
	}
}
