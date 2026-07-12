package builtins_test

import (
	"context"
	"testing"

	modbuiltins "github.com/deepnoodle-ai/risor/v2/pkg/builtins"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	_ "github.com/foohq/ren/builtins" // registers the "utf16" codec via init
)

var utf16Codec = object.NewString("utf16")

func encodeUTF16(t *testing.T, s string) []byte {
	t.Helper()
	enc, err := modbuiltins.Encode(context.Background(), object.NewString(s), utf16Codec)
	require.NoError(t, err)
	return enc.(*object.Bytes).Value()
}

func decodeUTF16(t *testing.T, b []byte) string {
	t.Helper()
	dec, err := modbuiltins.Decode(context.Background(), object.NewBytes(b), utf16Codec)
	require.NoError(t, err)
	return dec.(*object.String).Value()
}

func TestUTF16RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"ascii", "kernel32.dll"},
		{"empty", ""},
		{"unicode bmp", "café → ☂"},
		{"astral", "emoji 😀"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.text, decodeUTF16(t, encodeUTF16(t, tt.text)))
		})
	}
}

func TestUTF16LittleEndian(t *testing.T) {
	// "AB" -> U+0041, U+0042 encoded little-endian.
	require.Equal(t, []byte{0x41, 0x00, 0x42, 0x00}, encodeUTF16(t, "AB"))
}

func TestUTF16DecodeIsFaithful(t *testing.T) {
	// A trailing NUL code unit must be preserved, not trimmed.
	require.Equal(t, "A\x00", decodeUTF16(t, []byte{0x41, 0x00, 0x00, 0x00}))
}

func TestUTF16DecodeOddLength(t *testing.T) {
	// A trailing half code unit is dropped rather than erroring.
	require.Equal(t, "A", decodeUTF16(t, []byte{0x41, 0x00, 0x99}))
}

func TestUTF16EncodeRejectsNonString(t *testing.T) {
	_, err := modbuiltins.Encode(context.Background(), object.NewInt(1), utf16Codec)
	require.Error(t, err)
}
