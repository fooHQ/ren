package builtins_test

import (
	"context"
	"encoding/binary"
	"math"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren/builtins"
)

// Schema-building helpers keep the table entries readable.

func str(s string) object.Object {
	return object.NewString(s)
}

func i(n int64) object.Object {
	return object.NewInt(n)
}

func entry(parts ...object.Object) object.Object {
	return object.NewList(parts)
}

func schemaOf(entries ...object.Object) *object.List {
	return object.NewList(entries)
}

func f32Bytes(v float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
	return b
}

func f64Bytes(v float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	return b
}

func TestPacksize(t *testing.T) {
	filetime := schemaOf(
		entry(str("low"), str("int32")),
		entry(str("high"), str("int32")),
	)

	tests := []struct {
		name   string
		schema *object.List
		want   int64
	}{
		{"scalars", schemaOf(entry(str("a"), str("int8")), entry(str("b"), str("int64"))), 9},
		{"unsigned", schemaOf(entry(str("a"), str("uint8")), entry(str("b"), str("uint32"))), 5},
		{"floats", schemaOf(entry(str("a"), str("float32")), entry(str("b"), str("float64"))), 12},
		{"count", schemaOf(entry(str("a"), str("int16"), i(4))), 8},
		{"padding", schemaOf(entry(str("_"), str("int32")), entry(str("a"), str("int8"))), 5},
		{"nested", schemaOf(entry(str("t"), filetime)), 8},
		{"nested array", schemaOf(entry(str("t"), filetime, i(3))), 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtins.Packsize(context.Background(), tt.schema)
			require.NoError(t, err)
			require.Equal(t, object.NewInt(tt.want), got)
		})
	}
}

func TestUnpackScalarValues(t *testing.T) {
	tests := []struct {
		name string
		typ  string
		buf  []byte
		want object.Object
	}{
		{"int8 negative", "int8", []byte{0xFF}, object.NewInt(-1)},
		{"uint8 max", "uint8", []byte{0xFF}, object.NewInt(255)},
		{"int16 little-endian", "int16", []byte{0x00, 0x01}, object.NewInt(256)},
		{"int16 negative", "int16", []byte{0xFF, 0xFF}, object.NewInt(-1)},
		{"uint16 high bit", "uint16", []byte{0x40, 0x9C}, object.NewInt(40000)},
		{"int32 negative", "int32", []byte{0xFF, 0xFF, 0xFF, 0xFF}, object.NewInt(-1)},
		{"uint32 above int32 max", "uint32", []byte{0x00, 0x5E, 0xD0, 0xB2}, object.NewInt(3000000000)},
		{"int64", "int64", []byte{0x01, 0, 0, 0, 0, 0, 0, 0}, object.NewInt(1)},
		{"float32", "float32", f32Bytes(1.5), object.NewFloat(1.5)},
		{"float64", "float64", f64Bytes(2.5), object.NewFloat(2.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := schemaOf(entry(str("x"), str(tt.typ)))
			out, err := builtins.Unpack(context.Background(), schema, object.NewBytes(tt.buf))
			require.NoError(t, err)
			require.Equal(t, tt.want, out.(*object.Map).Get("x"))
		})
	}
}

func TestUnpackAggregates(t *testing.T) {
	ctx := context.Background()

	t.Run("uint8 count yields bytes", func(t *testing.T) {
		schema := schemaOf(entry(str("x"), str("uint8"), i(3)))
		out, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{1, 2, 3}))
		require.NoError(t, err)
		require.Equal(t, object.NewBytes([]byte{1, 2, 3}), out.(*object.Map).Get("x"))
	})

	t.Run("int16 count yields list", func(t *testing.T) {
		schema := schemaOf(entry(str("x"), str("int16"), i(2)))
		out, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{0x01, 0x00, 0xFF, 0xFF}))
		require.NoError(t, err)
		want := object.NewList([]object.Object{object.NewInt(1), object.NewInt(-1)})
		require.Equal(t, want, out.(*object.Map).Get("x"))
	})

	t.Run("float64 count yields list", func(t *testing.T) {
		schema := schemaOf(entry(str("x"), str("float64"), i(2)))
		buf := append(f64Bytes(1.5), f64Bytes(-2.5)...)
		out, err := builtins.Unpack(ctx, schema, object.NewBytes(buf))
		require.NoError(t, err)
		want := object.NewList([]object.Object{object.NewFloat(1.5), object.NewFloat(-2.5)})
		require.Equal(t, want, out.(*object.Map).Get("x"))
	})

	t.Run("nested yields map", func(t *testing.T) {
		filetime := schemaOf(entry(str("low"), str("int32")), entry(str("high"), str("int32")))
		schema := schemaOf(entry(str("t"), filetime))
		out, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{1, 0, 0, 0, 2, 0, 0, 0}))
		require.NoError(t, err)
		want := object.NewMap(map[string]object.Object{"low": object.NewInt(1), "high": object.NewInt(2)})
		require.Equal(t, want, out.(*object.Map).Get("t"))
	})

	t.Run("nested array yields list of maps", func(t *testing.T) {
		pair := schemaOf(entry(str("a"), str("int8")), entry(str("b"), str("int8")))
		schema := schemaOf(entry(str("p"), pair, i(2)))
		out, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{1, 2, 3, 4}))
		require.NoError(t, err)
		want := object.NewList([]object.Object{
			object.NewMap(map[string]object.Object{"a": object.NewInt(1), "b": object.NewInt(2)}),
			object.NewMap(map[string]object.Object{"a": object.NewInt(3), "b": object.NewInt(4)}),
		})
		require.Equal(t, want, out.(*object.Map).Get("p"))
	})

	t.Run("padding is skipped", func(t *testing.T) {
		schema := schemaOf(entry(str("_"), str("int32")), entry(str("a"), str("int8")))
		out, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{9, 9, 9, 9, 7}))
		require.NoError(t, err)
		m := out.(*object.Map)
		require.Equal(t, object.NewInt(7), m.Get("a"))
		require.Equal(t, object.Nil, m.Get("_"))
	})
}

func TestPackUnpackRoundTrip(t *testing.T) {
	ctx := context.Background()

	filetime := schemaOf(entry(str("low"), str("int32")), entry(str("high"), str("int32")))
	schema := schemaOf(
		entry(str("i8"), str("int8")),
		entry(str("u8"), str("uint8")),
		entry(str("u32"), str("uint32")),
		entry(str("f32"), str("float32")),
		entry(str("f64"), str("float64")),
		entry(str("raw"), str("uint8"), i(3)),
		entry(str("nums"), str("int16"), i(2)),
		entry(str("t"), filetime),
		entry(str("_"), str("int32")),
	)

	in := object.NewMap(map[string]object.Object{
		"i8":   object.NewInt(-5),
		"u8":   object.NewInt(200),        // would be -56 if sign-extended
		"u32":  object.NewInt(3000000000), // above int32 max
		"f32":  object.NewFloat(1.5),
		"f64":  object.NewFloat(3.141592653589793),
		"raw":  object.NewBytes([]byte{10, 20, 30}),
		"nums": object.NewList([]object.Object{object.NewInt(-1), object.NewInt(1000)}),
		"t":    object.NewMap(map[string]object.Object{"low": object.NewInt(11), "high": object.NewInt(22)}),
	})

	packed, err := builtins.Pack(ctx, schema, in)
	require.NoError(t, err)

	size, err := builtins.Packsize(ctx, schema)
	require.NoError(t, err)
	require.Len(t, packed.(*object.Bytes).Value(), int(size.(*object.Int).Value()))

	out, err := builtins.Unpack(ctx, schema, packed)
	require.NoError(t, err)
	m := out.(*object.Map)

	require.Equal(t, object.NewInt(-5), m.Get("i8"))
	require.Equal(t, object.NewInt(200), m.Get("u8"))
	require.Equal(t, object.NewInt(3000000000), m.Get("u32"))
	require.Equal(t, object.NewFloat(1.5), m.Get("f32"))
	require.Equal(t, object.NewFloat(3.141592653589793), m.Get("f64"))
	require.Equal(t, object.NewBytes([]byte{10, 20, 30}), m.Get("raw"))
	require.Equal(t, object.NewList([]object.Object{object.NewInt(-1), object.NewInt(1000)}), m.Get("nums"))
	require.Equal(t, object.NewMap(map[string]object.Object{"low": object.NewInt(11), "high": object.NewInt(22)}), m.Get("t"))
	require.Equal(t, object.Nil, m.Get("_"))
}

func TestPackPaddingIsZeroFilled(t *testing.T) {
	schema := schemaOf(entry(str("_"), str("int16")), entry(str("a"), str("int8")))
	in := object.NewMap(map[string]object.Object{"a": object.NewInt(7)})
	packed, err := builtins.Pack(context.Background(), schema, in)
	require.NoError(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x07}, packed.(*object.Bytes).Value())
}

func TestSchemaErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		schema *object.List
		errMsg string
	}{
		{"unknown type", schemaOf(entry(str("a"), str("int24"))), `unknown type "int24"`},
		{"entry not a list", object.NewList([]object.Object{str("a")}), "must be a list"},
		{"entry too short", schemaOf(entry(str("a"))), "at least 2 elements"},
		{"bad type value", schemaOf(entry(str("a"), i(1))), "must be a string or list"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := builtins.Packsize(ctx, tt.schema)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestUnpackErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("buffer too small", func(t *testing.T) {
		schema := schemaOf(entry(str("a"), str("int32")))
		_, err := builtins.Unpack(ctx, schema, object.NewBytes([]byte{1, 2}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "buffer too small")
	})

	t.Run("wrong argument count", func(t *testing.T) {
		_, err := builtins.Unpack(ctx, schemaOf())
		require.Error(t, err)
	})

	t.Run("schema not a list", func(t *testing.T) {
		_, err := builtins.Unpack(ctx, str("nope"), object.NewBytes(nil))
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected list")
	})

	t.Run("buffer not bytes", func(t *testing.T) {
		_, err := builtins.Unpack(ctx, schemaOf(entry(str("a"), str("int8"))), str("nope"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected bytes")
	})
}

func TestPackErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("data not a map", func(t *testing.T) {
		_, err := builtins.Pack(ctx, schemaOf(entry(str("a"), str("int8"))), str("nope"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected map")
	})

	t.Run("field wrong type", func(t *testing.T) {
		schema := schemaOf(entry(str("a"), str("int32")))
		in := object.NewMap(map[string]object.Object{"a": str("not an int")})
		_, err := builtins.Pack(ctx, schema, in)
		require.Error(t, err)
	})

	t.Run("byte field wrong type", func(t *testing.T) {
		schema := schemaOf(entry(str("a"), str("uint8"), i(2)))
		in := object.NewMap(map[string]object.Object{"a": object.NewInt(5)})
		_, err := builtins.Pack(ctx, schema, in)
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected bytes")
	})
}
