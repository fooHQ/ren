package builtins

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

// Pack serializes a map of values into a little-endian byte buffer according to
// a schema. It takes two arguments: the schema (a list of field descriptors)
// and the data map. See parseSchema for the schema format.
func Pack(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("pack", 2, len(args))
	}
	schema, ok := args[0].(*object.List)
	if !ok {
		return nil, fmt.Errorf("pack: expected list for schema, got %s", args[0].Type())
	}
	m, ok := args[1].(*object.Map)
	if !ok {
		return nil, fmt.Errorf("pack: expected map for data, got %s", args[1].Type())
	}
	fields, err := parseSchema(schema)
	if err != nil {
		return nil, err
	}
	totalSize := 0
	for _, f := range fields {
		s, err := fieldSize(f)
		if err != nil {
			return nil, err
		}
		totalSize += s
	}
	buf := make([]byte, totalSize)
	offset := 0
	for _, f := range fields {
		size, err := packField(f, m, buf, offset)
		if err != nil {
			return nil, err
		}
		offset += size
	}
	return object.NewBytes(buf), nil
}

// Packsize returns the total byte size of a schema without packing any data. It
// takes a single schema argument.
func Packsize(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, object.NewArgsError("packsize", 1, len(args))
	}
	schema, ok := args[0].(*object.List)
	if !ok {
		return nil, fmt.Errorf("packsize: expected list for schema, got %s", args[0].Type())
	}
	fields, err := parseSchema(schema)
	if err != nil {
		return nil, err
	}
	total := 0
	for _, f := range fields {
		s, err := fieldSize(f)
		if err != nil {
			return nil, err
		}
		total += s
	}
	return object.NewInt(int64(total)), nil
}

// Unpack deserializes a little-endian byte buffer into a map according to a
// schema, the inverse of Pack. It takes two arguments: the schema and the byte
// buffer. Fields named "_" are decoded but omitted from the result.
func Unpack(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, object.NewArgsError("unpack", 2, len(args))
	}
	schema, ok := args[0].(*object.List)
	if !ok {
		return nil, fmt.Errorf("unpack: expected list for schema, got %s", args[0].Type())
	}
	buf, ok := args[1].(*object.Bytes)
	if !ok {
		return nil, fmt.Errorf("unpack: expected bytes for buffer, got %s", args[1].Type())
	}
	m, _, err := unpackStruct(schema, buf.Value(), 0)
	if err != nil {
		return nil, err
	}
	return object.NewMap(m), nil
}

// field is a single parsed schema entry: a named scalar or nested struct,
// optionally repeated count times.
type field struct {
	name   string
	typ    string       // scalar type name, e.g. "int32", "uint16", "float64"
	nested *object.List // non-nil when field is a nested schema
	count  int          // number of elements; 1 if omitted
}

// parseSchema converts a schema list into fields. Each entry is a list of
// [name, type, count?] where type is a scalar type name or a nested schema
// list, and the optional count gives the number of repeated elements.
func parseSchema(schema *object.List) ([]field, error) {
	items := schema.Value()
	fields := make([]field, 0, len(items))
	for i, item := range items {
		entry, ok := item.(*object.List)
		if !ok {
			return nil, fmt.Errorf("pack: schema entry %d must be a list", i)
		}
		vals := entry.Value()
		if len(vals) < 2 {
			return nil, fmt.Errorf("pack: schema entry %d must have at least 2 elements", i)
		}
		name, err := object.AsString(vals[0])
		if err != nil {
			return nil, fmt.Errorf("pack: schema entry %d name: %w", i, err)
		}
		f := field{name: name, count: 1}
		switch v := vals[1].(type) {
		case *object.String:
			f.typ = v.Value()
			switch f.typ {
			case "int8", "int16", "int32", "int64",
				"uint8", "uint16", "uint32",
				"float32", "float64":
			default:
				return nil, fmt.Errorf("pack: schema entry %d: unknown type %q", i, f.typ)
			}
		case *object.List:
			f.nested = v
		default:
			return nil, fmt.Errorf("pack: schema entry %d: type must be a string or list", i)
		}
		if len(vals) >= 3 {
			count, err := object.AsInt(vals[2])
			if err != nil {
				return nil, fmt.Errorf("pack: schema entry %d count: %w", i, err)
			}
			f.count = int(count)
		}
		fields = append(fields, f)
	}
	return fields, nil
}

// scalarSize returns the byte size of a scalar type, or 0 if unknown.
func scalarSize(typ string) int {
	switch typ {
	case "int8", "uint8":
		return 1
	case "int16", "uint16":
		return 2
	case "int32", "uint32", "float32":
		return 4
	case "int64", "float64":
		return 8
	}
	return 0
}

// fieldSize returns the total byte size of a field, recursing into nested
// schemas and accounting for the element count.
func fieldSize(f field) (int, error) {
	if f.nested != nil {
		fields, err := parseSchema(f.nested)
		if err != nil {
			return 0, err
		}
		structSize := 0
		for _, sf := range fields {
			s, err := fieldSize(sf)
			if err != nil {
				return 0, err
			}
			structSize += s
		}
		return structSize * f.count, nil
	}
	return scalarSize(f.typ) * f.count, nil
}

// unpackStruct decodes the fields of a schema starting at offset, returning the
// resulting map and the offset just past the decoded bytes.
func unpackStruct(schema *object.List, buf []byte, offset int) (map[string]object.Object, int, error) {
	fields, err := parseSchema(schema)
	if err != nil {
		return nil, 0, err
	}
	result := make(map[string]object.Object, len(fields))
	for _, f := range fields {
		val, size, err := unpackField(f, buf, offset)
		if err != nil {
			return nil, 0, err
		}
		if f.name != "_" {
			result[f.name] = val
		}
		offset += size
	}
	return result, offset, nil
}

// unpackField decodes a single field, returning its value and byte size.
func unpackField(f field, buf []byte, offset int) (object.Object, int, error) {
	if f.nested != nil {
		return unpackNested(f, buf, offset)
	}
	return unpackScalar(f, buf, offset)
}

// unpackNested decodes a nested-struct field, yielding a map for a single
// element or a list of maps when repeated.
func unpackNested(f field, buf []byte, offset int) (object.Object, int, error) {
	if f.count == 1 {
		m, end, err := unpackStruct(f.nested, buf, offset)
		if err != nil {
			return nil, 0, err
		}
		return object.NewMap(m), end - offset, nil
	}
	items := make([]object.Object, f.count)
	totalSize := 0
	for i := range f.count {
		m, end, err := unpackStruct(f.nested, buf, offset)
		if err != nil {
			return nil, 0, err
		}
		size := end - offset
		items[i] = object.NewMap(m)
		offset += size
		totalSize += size
	}
	return object.NewList(items), totalSize, nil
}

// unpackScalar decodes a scalar field, yielding a single value, a list of
// values, or (for 8-bit integer arrays) a bytes value.
func unpackScalar(f field, buf []byte, offset int) (object.Object, int, error) {
	elemSize := scalarSize(f.typ)
	totalSize := elemSize * f.count

	if offset+totalSize > len(buf) {
		return nil, 0, fmt.Errorf("unpack: buffer too small at offset %d: need %d bytes, have %d",
			offset, totalSize, len(buf)-offset)
	}
	data := buf[offset : offset+totalSize]

	switch f.typ {
	case "float32", "float64":
		if f.count == 1 {
			v, err := readFloat(f.typ, data)
			if err != nil {
				return nil, 0, err
			}
			return object.NewFloat(v), totalSize, nil
		}
		items := make([]object.Object, f.count)
		for i := range f.count {
			v, err := readFloat(f.typ, data[i*elemSize:])
			if err != nil {
				return nil, 0, err
			}
			items[i] = object.NewFloat(v)
		}
		return object.NewList(items), totalSize, nil
	}

	if f.count == 1 {
		v, err := readInt(f.typ, data)
		if err != nil {
			return nil, 0, err
		}
		return object.NewInt(v), totalSize, nil
	}

	if f.typ == "int8" || f.typ == "uint8" {
		result := make([]byte, f.count)
		copy(result, data)
		return object.NewBytes(result), totalSize, nil
	}

	items := make([]object.Object, f.count)
	for i := range f.count {
		v, err := readInt(f.typ, data[i*elemSize:])
		if err != nil {
			return nil, 0, err
		}
		items[i] = object.NewInt(v)
	}
	return object.NewList(items), totalSize, nil
}

// readInt reads a little-endian integer of the given type from data.
func readInt(typ string, data []byte) (int64, error) {
	switch typ {
	case "int8":
		return int64(int8(data[0])), nil
	case "int16":
		return int64(int16(binary.LittleEndian.Uint16(data))), nil
	case "int32":
		return int64(int32(binary.LittleEndian.Uint32(data))), nil
	case "int64":
		return int64(binary.LittleEndian.Uint64(data)), nil
	case "uint8":
		return int64(data[0]), nil
	case "uint16":
		return int64(binary.LittleEndian.Uint16(data)), nil
	case "uint32":
		return int64(binary.LittleEndian.Uint32(data)), nil
	}
	return 0, fmt.Errorf("unpack: unknown type %q", typ)
}

// readFloat reads a little-endian float of the given type from data.
func readFloat(typ string, data []byte) (float64, error) {
	switch typ {
	case "float32":
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(data))), nil
	case "float64":
		return math.Float64frombits(binary.LittleEndian.Uint64(data)), nil
	}
	return 0, fmt.Errorf("unpack: unknown type %q", typ)
}

// packField encodes a single field from m into buf, returning its byte size.
func packField(f field, m *object.Map, buf []byte, offset int) (int, error) {
	if f.nested != nil {
		return packNested(f, m, buf, offset)
	}
	return packScalar(f, m, buf, offset)
}

// packNested encodes a nested-struct field, reading a map for a single element
// or a list of maps when repeated.
func packNested(f field, m *object.Map, buf []byte, offset int) (int, error) {
	val := m.GetWithDefault(f.name, object.Nil)

	if f.count == 1 {
		nested, ok := val.(*object.Map)
		if !ok {
			return 0, fmt.Errorf("pack: field %q: expected map, got %s", f.name, val.Type())
		}
		return packStruct(f.nested, nested, buf, offset)
	}

	list, ok := val.(*object.List)
	if !ok {
		return 0, fmt.Errorf("pack: field %q: expected list, got %s", f.name, val.Type())
	}
	totalSize := 0
	for i, item := range list.Value() {
		nested, ok := item.(*object.Map)
		if !ok {
			return 0, fmt.Errorf("pack: field %q[%d]: expected map, got %s", f.name, i, item.Type())
		}
		size, err := packStruct(f.nested, nested, buf, offset)
		if err != nil {
			return 0, err
		}
		offset += size
		totalSize += size
	}
	return totalSize, nil
}

// packStruct encodes all fields of a schema from m into buf, returning the
// number of bytes written.
func packStruct(schema *object.List, m *object.Map, buf []byte, offset int) (int, error) {
	fields, err := parseSchema(schema)
	if err != nil {
		return 0, err
	}
	start := offset
	for _, f := range fields {
		size, err := packField(f, m, buf, offset)
		if err != nil {
			return 0, err
		}
		offset += size
	}
	return offset - start, nil
}

// packScalar encodes a scalar field from m into buf. Fields named "_" are
// skipped, leaving zero padding. Repeated 8-bit integers are read from a bytes
// value; other repeated scalars are read from a list.
func packScalar(f field, m *object.Map, buf []byte, offset int) (int, error) {
	elemSize := scalarSize(f.typ)
	totalSize := elemSize * f.count

	if f.name == "_" {
		return totalSize, nil
	}

	val := m.GetWithDefault(f.name, object.Nil)

	switch f.typ {
	case "float32", "float64":
		if f.count == 1 {
			x, err := object.AsFloat(val)
			if err != nil {
				return 0, fmt.Errorf("pack: field %q: %w", f.name, err)
			}
			writeFloat(f.typ, buf[offset:], x)
			return totalSize, nil
		}
		list, ok := val.(*object.List)
		if !ok {
			return 0, fmt.Errorf("pack: field %q: expected list, got %s", f.name, val.Type())
		}
		for i, item := range list.Value() {
			x, err := object.AsFloat(item)
			if err != nil {
				return 0, fmt.Errorf("pack: field %q[%d]: %w", f.name, i, err)
			}
			writeFloat(f.typ, buf[offset+i*elemSize:], x)
		}
		return totalSize, nil
	}

	if f.count == 1 {
		n, err := object.AsInt(val)
		if err != nil {
			return 0, fmt.Errorf("pack: field %q: %w", f.name, err)
		}
		writeInt(f.typ, buf[offset:], n)
		return totalSize, nil
	}

	if f.typ == "int8" || f.typ == "uint8" {
		b, ok := val.(*object.Bytes)
		if !ok {
			return 0, fmt.Errorf("pack: field %q: expected bytes, got %s", f.name, val.Type())
		}
		copy(buf[offset:offset+totalSize], b.Value())
		return totalSize, nil
	}

	list, ok := val.(*object.List)
	if !ok {
		return 0, fmt.Errorf("pack: field %q: expected list, got %s", f.name, val.Type())
	}
	for i, item := range list.Value() {
		n, err := object.AsInt(item)
		if err != nil {
			return 0, fmt.Errorf("pack: field %q[%d]: %w", f.name, i, err)
		}
		writeInt(f.typ, buf[offset+i*elemSize:], n)
	}
	return totalSize, nil
}

// writeInt writes v as a little-endian integer of the given type into buf.
func writeInt(typ string, buf []byte, v int64) {
	switch typ {
	case "int8", "uint8":
		buf[0] = byte(v)
	case "int16", "uint16":
		binary.LittleEndian.PutUint16(buf, uint16(v))
	case "int32", "uint32":
		binary.LittleEndian.PutUint32(buf, uint32(v))
	case "int64":
		binary.LittleEndian.PutUint64(buf, uint64(v))
	}
}

// writeFloat writes v as a little-endian float of the given type into buf.
func writeFloat(typ string, buf []byte, v float64) {
	switch typ {
	case "float32":
		binary.LittleEndian.PutUint32(buf, math.Float32bits(float32(v)))
	case "float64":
		binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	}
}
