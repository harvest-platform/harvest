package indexer

const (
	StringType  = "string"
	TextType    = "text"
	IntegerType = "integer"
	FloatType   = "float"
	BooleanType = "boolean"
	DateType    = "date"
	ObjectType  = "object"
)

type Field struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Signed   bool              `json:"signed"`
	Bits     uint8             `json:"bits"`
	Length   uint16            `json:"length"`
	Multiple bool              `json:"multiple"`
	Nullable bool              `json:"nullable"`
	Fields   map[string]*Field `json:"fields,omitempty"`
}
