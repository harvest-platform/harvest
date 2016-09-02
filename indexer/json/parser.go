package json

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"time"

	"github.com/harvest-platform/harvest/bitsize"
	"github.com/harvest-platform/harvest/indexer"
)

const (
	floatingPointError = 1e-9
	textThreshold      = 50
)

var dateFormats = []string{
	"2006-01-02",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05Z07:00",
}

func isDate(s string) bool {
	for _, layout := range dateFormats {
		if _, err := time.Parse(layout, s); err == nil {
			return true
		}
	}

	return false
}

func isInteger(f float64) bool {
	_, frac := math.Modf(math.Abs(f))
	return frac < floatingPointError || frac > 1.0-floatingPointError
}

// generalizeType takes two types and attempts to generalize the type between them.
func generalizeType(t1, t2 string) (string, bool) {
	if t1 == t2 {
		return t1, true
	}

	if t1 == "" {
		return t2, true
	}

	if t2 == "" {
		return t1, true
	}

	switch t1 {
	case indexer.IntegerType:
		if t2 == indexer.FloatType {
			return indexer.FloatType, true
		}

	case indexer.FloatType:
		if t2 == indexer.IntegerType {
			return indexer.FloatType, true
		}

	case indexer.TextType:
		switch t2 {
		case indexer.StringType, indexer.DateType:
			return indexer.TextType, true
		}

	case indexer.StringType:
		switch t2 {
		case indexer.TextType, indexer.DateType:
			return indexer.TextType, true
		}

	case indexer.DateType:
		switch t2 {
		case indexer.StringType:
			return indexer.StringType, true
		case indexer.TextType:
			return indexer.TextType, true
		}
	}

	return "", false
}

func mergeFields(f1, f2 *indexer.Field) *indexer.Field {
	if f1 == nil {
		return f2
	}

	if f2 == nil {
		return f1
	}

	// Different cardinality.
	if f1.Multiple != f2.Multiple {
		log.Print("Cardinality conflict")
	}

	if f2.Nullable == true {
		f1.Nullable = true
	}

	if f1.Fields == nil {
		f1.Fields = f2.Fields
	} else {
		for k, c1 := range f1.Fields {
			if c2, ok := f2.Fields[k]; ok {
				f1.Fields[k] = mergeFields(c1, c2)
			}
		}

		for k, c2 := range f2.Fields {
			if _, ok := f1.Fields[k]; !ok {
				f1.Fields[k] = c2
			}
		}
	}

	// Different types, attempt to generalize.
	typ, ok := generalizeType(f1.Type, f2.Type)
	if !ok {
		log.Printf("Type conflict: %s, %s", f1.Type, f2.Type)
	} else {
		f1.Type = typ
	}

	return f1
}

func parseArray(k string, a []interface{}) *indexer.Field {
	f := &indexer.Field{
		Name: k,
	}

	for _, v := range a {
		if _, ok := v.([]interface{}); ok {
			log.Print("array cannot contain sub-arrays")
			return nil
		}

		f = mergeFields(f, parseFieldValue(k, v))
	}

	return f
}

func parseFieldValue(k string, v interface{}) *indexer.Field {
	f := &indexer.Field{
		Name: k,
	}

	switch x := v.(type) {
	case nil:
		f.Nullable = true

	// Nested object.
	case map[string]interface{}:
		f.Type = indexer.ObjectType
		f.Fields = parseMap(x)

	// Possibly nested objects or primitives.
	case []interface{}:
		if len(x) > 0 {
			f = parseArray(k, x)
		}
		f.Multiple = true

	case bool:
		f.Type = indexer.BooleanType

	case string:
		// TODO: more sophisticated way to detect text.
		if isDate(x) {
			f.Type = indexer.DateType
		} else {
			f.Length = uint16(len(x))

			if f.Length >= textThreshold {
				f.Type = indexer.TextType
			} else {
				f.Type = indexer.StringType
			}
		}

	case float64:
		if x < 0 {
			f.Signed = true
		}

		if isInteger(x) {
			f.Type = indexer.IntegerType

			if x < 0 {
				f.Bits = bitsize.Int(int64(x))
			} else {
				f.Bits = bitsize.Uint(uint64(x))
			}
		} else {
			f.Type = indexer.FloatType
		}
	}

	return f
}

func parseMap(m map[string]interface{}) map[string]*indexer.Field {
	fields := make(map[string]*indexer.Field, len(m))

	for k, v := range m {
		fields[k] = parseFieldValue(k, v)
	}

	return fields
}

// Infer takes JSON decoded record and infers the schema from the document.
func Infer(v map[string]interface{}) map[string]*indexer.Field {
	return parseMap(v)
}

func UnmarshalInfer(b []byte) (map[string]*indexer.Field, error) {
	var m map[string]interface{}

	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return Infer(m), nil
}

func DecodeInfer(r io.Reader) (map[string]*indexer.Field, error) {
	var m map[string]interface{}

	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}

	return Infer(m), nil

}
