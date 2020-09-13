package engine

import (
	"github.com/xitongsys/parquet-go/parquet"
)

type Schema interface {
}

type Field struct {
	Name    string
	SE      *parquet.SchemaElement
	Type    string
	MFields map[string]*Field
	Fields  []*Field
}

type ParquetSchema struct {
	Name    string
	MFields map[string]*Field
	Fields  []*Field
}

func (s *ParquetSchema) GetName() string {
	return s.Name
}

func NewParquetSchema(schemas []*parquet.SchemaElement) *ParquetSchema {
	stack := make([]*Field, 0)
	root := NewField(schemas[0])
	stack = append(stack, root)
	pos := 1
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		numChildren := int(node.SE.GetNumChildren())
		if len(node.Fields) < numChildren {
			field := NewField(schemas[pos])
			node.Fields = append(node.Fields, field)
			node.MFields[field.Name] = field
			stack = append(stack, field)
			pos++
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return &ParquetSchema{Name: root.Name, MFields: root.MFields, Fields: root.Fields}
}

func NewField(schema *parquet.SchemaElement) *Field {
	return &Field{
		Name:    schema.Name,
		SE:      schema,
		MFields: make(map[string]*Field),
		Fields:  make([]*Field, 0, schema.GetNumChildren()),
	}
}
