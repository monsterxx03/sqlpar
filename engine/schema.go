package engine

import (
	"fmt"
	"github.com/xitongsys/parquet-go/parquet"
	"strconv"
	"strings"
)

const (
	TYPE_STRUCT = "STRUCT"
	TYPE_LIST   = "LIST"
	TYPE_MAP    = "MAP"
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
	switch f.Type {
	case TYPE_LIST:
		// LIST.ELEMENT
		return fmt.Sprintf("array<%s>", f.Fields[0].Fields[0])
	case TYPE_MAP:
		// map<KEY_VALUE.KEY, KEY_VALUE.VALUE>
		return fmt.Sprintf("map<%s,%s>", f.Fields[0].Fields[0], f.Fields[0].Fields[1])
	case TYPE_STRUCT:
		fields := make([]string, 0, len(f.Fields))
		for _, _f := range f.Fields {
			fields = append(fields, fmt.Sprintf("%s:%s", _f.Name, _f))
		}
		return fmt.Sprintf("struct<%s>", strings.Join(fields, ","))
	default:
		return f.Type
	}
}

type ParquetSchema struct {
	Name    string
	MFields map[string]*Field
	Fields  []*Field
}

func (s *ParquetSchema) GetFieldPath(subs []string) (string, error) {
	paths := make([]string, 0)
	f := s.MFields[subs[0]]
	paths = append(paths, subs[0])
	subs = subs[1:]
	for len(subs) > 0 {
		sub := subs[0]
		idx, err := strconv.Atoi(sub)
		if err != nil {
			switch f.Type{
			case TYPE_MAP:
				paths = append(paths, []string{"Key_value", "Value"}...)
				f = f.Fields[0].Fields[1]
			case TYPE_STRUCT:
				found := false
				for _, _f := range f.Fields{
					if _f.Name == sub {
						found = true
						paths = append(paths, _f.Name)
						f = _f
						break
					}
				}
				if !found {
					return "", fmt.Errorf("can't find %s on %s:%s", sub, f.Name, f)
				}
			default:
				return "", fmt.Errorf("unsuport to get `%s` on %s:%s", sub, f.Name, f)
			}
		} else {
			if f.Type != TYPE_LIST {
				return "", fmt.Errorf("unsupport to retrive index `%d` on %s:%s", idx, f.Name, f)
			}
			paths = append(paths, "List", "Element")
			f = f.Fields[0].Fields[0]
		}
		subs = subs[1:]
	}
	return strings.Join(paths, "."), nil
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

func NewField(se *parquet.SchemaElement) *Field {
	TYPE := ""
	if se.GetNumChildren() == 0 {
		if se.ConvertedType != nil {
			TYPE = se.ConvertedType.String()
		} else if se.Type != nil {
			TYPE = se.Type.String()
		}
	} else {
		if se.ConvertedType == nil {
			TYPE = TYPE_STRUCT
		} else {
			switch *se.ConvertedType {
			case parquet.ConvertedType_MAP:
				TYPE = TYPE_MAP
			case parquet.ConvertedType_LIST:
				TYPE = TYPE_LIST
			default:
				TYPE = se.ConvertedType.String()
			}
		}
	}
	return &Field{
		Name:    se.Name,
		SE:      se,
		Type:    TYPE,
		MFields: make(map[string]*Field),
		Fields:  make([]*Field, 0, se.GetNumChildren()),
	}
}
