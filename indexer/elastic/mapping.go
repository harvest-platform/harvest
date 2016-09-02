package elastic

import "github.com/harvest-platform/harvest/indexer"

func mapFields(fields map[string]*indexer.Field) map[string]interface{} {
	m := make(map[string]interface{})

	var v map[string]interface{}

	for _, f := range fields {
		if f.Multiple && f.Type == indexer.ObjectType {
			v = nestedMapping(mapFields(f.Fields))
		} else {
			switch f.Type {
			case indexer.ObjectType:
				v = objectMapping(mapFields(f.Fields))

			case indexer.StringType:
				v = stringMapping()

			case indexer.TextType:
				v = textMapping()

			case indexer.IntegerType:
				v = integerMapping(f.Bits)

			case indexer.FloatType:
				v = floatMapping(f.Bits)

			case indexer.BooleanType:
				v = booleanMapping()

			case indexer.DateType:
				v = dateMapping()
			}
		}

		m[f.Name] = v
	}

	return m
}

func objectMapping(props map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"dynamic":    "strict",
		"properties": props,
	}
}

func nestedMapping(props map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":       "strict",
		"dynamic":    false,
		"properties": props,
	}
}

func stringMapping() map[string]interface{} {
	return map[string]interface{}{
		"type":       "string",
		"index":      "not_analyzed",
		"doc_values": true,
		"fields": map[string]interface{}{
			"analyzed": map[string]interface{}{
				"type":          "string",
				"index":         "analyzed",
				"index_options": "positions",
				"fielddata": map[string]interface{}{
					"format": "disabled",
				},
			},
		},
	}
}

func textMapping() map[string]interface{} {
	return map[string]interface{}{
		"type":          "string",
		"index":         "analyzed",
		"index_options": "positions",
		"fielddata": map[string]interface{}{
			"format": "disabled",
		},
	}
}

func integerMapping(size uint8) map[string]interface{} {
	var typ string

	if size <= 8 {
		typ = "byte"
	} else if size <= 16 {
		typ = "short"
	} else if size <= 32 {
		typ = "integer"
	} else if size <= 64 {
		typ = "long"
	} else {
		panic("integer size must be 8, 16, 32 or 64")
	}

	return map[string]interface{}{
		"type":             typ,
		"coerce":           false,
		"ignore_malformed": false,
		"index":            "not_analyzed",
		"doc_values":       true,
	}
}

func floatMapping(size uint8) map[string]interface{} {
	var typ string

	if size <= 32 {
		typ = "float"
	} else if size <= 64 {
		typ = "double"
	} else {
		panic("float size must be 32 or 64")
	}

	return map[string]interface{}{
		"type":             typ,
		"coerce":           false,
		"ignore_malformed": false,
		"index":            "not_analyzed",
		"doc_values":       true,
	}
}

func booleanMapping() map[string]interface{} {
	return map[string]interface{}{
		"type":       "boolean",
		"index":      "not_analyzed",
		"doc_values": true,
	}
}

func dateMapping() map[string]interface{} {
	return map[string]interface{}{
		"type":             "date",
		"index":            "not_analyzed",
		"format":           "strict_date_optional_time",
		"ignore_malformed": false,
	}
}

// Generate generatees an Elasticsearch mapping for the fields.
func Generate(fields map[string]*indexer.Field) map[string]interface{} {
	return mapFields(fields)
}
