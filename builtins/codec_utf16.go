package builtins

import (
	"context"
	"encoding/binary"
	"unicode/utf16"

	modbuiltins "github.com/deepnoodle-ai/risor/v2/pkg/builtins"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
)

// Registers a "utf16" codec so scripts can convert between strings and
// little-endian UTF-16 bytes via the standard encode/decode builtins:
//
//	encode(s, "utf16")   // string -> UTF-16LE bytes
//	decode(b, "utf16")   // UTF-16LE bytes -> string
//
// Decoding is faithful: it does not trim at a NUL terminator. Callers
// reading fixed-size C string fields should trim it themselves.
func init() {
	err := modbuiltins.RegisterCodec("utf16", &modbuiltins.Codec{
		Encode: encodeUTF16,
		Decode: decodeUTF16,
	})
	if err != nil {
		panic(err)
	}
}

// encodeUTF16 encodes a string as little-endian UTF-16 bytes.
func encodeUTF16(_ context.Context, obj object.Object) (object.Object, error) {
	s, err := object.AsString(obj)
	if err != nil {
		return nil, err
	}
	units := utf16.Encode([]rune(s))
	buf := make([]byte, len(units)*2)
	for i, u := range units {
		binary.LittleEndian.PutUint16(buf[i*2:], u)
	}
	return object.NewBytes(buf), nil
}

// decodeUTF16 decodes little-endian UTF-16 bytes into a string. It does not
// stop at a NUL terminator.
func decodeUTF16(_ context.Context, obj object.Object) (object.Object, error) {
	b, err := object.AsBytes(obj)
	if err != nil {
		return nil, err
	}
	units := make([]uint16, len(b)/2)
	for i := range units {
		units[i] = binary.LittleEndian.Uint16(b[i*2:])
	}
	return object.NewString(string(utf16.Decode(units))), nil
}
