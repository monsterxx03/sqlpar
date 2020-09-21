package engine

import (
	"fmt"
	"github.com/xitongsys/parquet-go/parquet"
	"strings"
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

func (f *Field) String() string {
	if len(f.Fields) == 0 {
		if f.SE.ConvertedType != nil {
			return f.SE.ConvertedType.String()
		}
		if f.SE.Type != nil {
			return f.SE.Type.String()
		}
	}
	if f.SE.ConvertedType == nil {
		fields := make([]string, 0, len(f.Fields))
		for _, _f := range f.Fields {
			fields = append(fields, fmt.Sprintf("%s:%s", _f.Name, _f))
		}
		return fmt.Sprintf("struct<%s>", strings.Join(fields, ","))
	}
	switch *f.SE.ConvertedType {
	case parquet.ConvertedType_LIST:
		// LIST.ELEMENT
		return fmt.Sprintf("array<%s>", f.Fields[0].Fields[0])
	case parquet.ConvertedType_MAP:
		// map<KEY_VALUE.KEY, KEY_VALUE.VALUE>
		return fmt.Sprintf("map<%s,%s>", f.Fields[0].Fields[0], f.Fields[0].Fields[1])
	case parquet.ConvertedType_JSON:
		return "json"
	default:
		return f.SE.ConvertedType.String()
	}
}

type ParquetSchema struct {
	Name    string
	MFields map[string]*Field
	Fields  []*Field
}

func (s *ParquetSchema) GetName() string {
	return s.Name
}

func (s *ParquetSchema) GetAllFieldNames() []string {
	cols := make([]string, 0, len(s.Fields))
	for _, f := range s.Fields {
		cols = append(cols, f.Name)
	}
	return cols
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
