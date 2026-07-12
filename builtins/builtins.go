// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

// Package builtins defines the global functions available to every Ren script.
// It re-exports a curated subset of Risor's built-ins and adds Ren-specific
// ones such as import, print, printf, and the pack/unpack family. It also
// registers the "utf16" encode/decode codec on import.
package builtins

import (
	"bytes"
	"context"
	"fmt"
	"maps"

	modbuiltins "github.com/deepnoodle-ai/risor/v2/pkg/builtins"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"

	"github.com/foohq/ren"
)

var builtins = map[string]*object.Builtin{
	"all":      object.NewBuiltin("all", modbuiltins.All),
	"any":      object.NewBuiltin("any", modbuiltins.Any),
	"assert":   object.NewBuiltin("assert", modbuiltins.Assert),
	"bool":     object.NewBuiltin("bool", modbuiltins.Bool),
	"byte":     object.NewBuiltin("byte", modbuiltins.Byte),
	"bytes":    object.NewBuiltin("bytes", modbuiltins.Bytes),
	"call":     object.NewBuiltin("call", modbuiltins.Call),
	"chunk":    object.NewBuiltin("chunk", modbuiltins.Chunk),
	"coalesce": object.NewBuiltin("coalesce", modbuiltins.Coalesce),
	"decode":   object.NewBuiltin("decode", modbuiltins.Decode),
	"encode":   object.NewBuiltin("encode", modbuiltins.Encode),
	"error":    object.NewBuiltin("error", modbuiltins.Error),
	"filter":   object.NewBuiltin("filter", modbuiltins.Filter),
	"float":    object.NewBuiltin("float", modbuiltins.Float),
	"getattr":  object.NewBuiltin("getattr", modbuiltins.GetAttr),
	"import":   object.NewBuiltin("import", Import),
	"int":      object.NewBuiltin("int", modbuiltins.Int),
	"keys":     object.NewBuiltin("keys", modbuiltins.Keys),
	"len":      object.NewBuiltin("len", modbuiltins.Len),
	"list":     object.NewBuiltin("list", modbuiltins.List),
	"range":    object.NewBuiltin("range", modbuiltins.Range),
	"reversed": object.NewBuiltin("reversed", modbuiltins.Reversed),
	"sorted":   object.NewBuiltin("sorted", modbuiltins.Sorted),
	"sprintf":  object.NewBuiltin("sprintf", modbuiltins.Sprintf),
	"string":   object.NewBuiltin("string", modbuiltins.String),
	"type":     object.NewBuiltin("type", modbuiltins.Type),
	"pack":     object.NewBuiltin("pack", Pack),
	"packsize": object.NewBuiltin("packsize", Packsize),
	"unpack":   object.NewBuiltin("unpack", Unpack),
	"print":    object.NewBuiltin("print", Print),
	"printf":   object.NewBuiltin("printf", Printf),
}

// Print writes its arguments to standard output separated by spaces and
// followed by a newline. Strings and bytes are written verbatim; other values
// are written via their Inspect representation. It accepts 1 to 64 arguments.
func Print(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) < 1 || len(args) > 64 {
		return nil, object.NewArgsRangeError("print", 1, 64, len(args))
	}
	var b bytes.Buffer
	for i := range args {
		if i > 0 {
			b.WriteByte(' ')
		}
		var err error
		switch obj := args[i].(type) {
		case *object.String:
			_, err = b.WriteString(obj.Value())
		case *object.Bytes:
			_, err = b.Write(obj.Value())
		default:
			_, err = b.WriteString(obj.Inspect())
		}
		if err != nil {
			return nil, err
		}
	}
	b.WriteByte('\n')
	_, err := ren.GetOS(ctx).Stdout().Write(b.Bytes())
	if err != nil {
		return nil, err
	}
	return object.Nil, nil
}

// Printf formats its trailing arguments according to the first (format-string)
// argument and writes the result to standard output. It accepts 1 to 64
// arguments.
func Printf(ctx context.Context, args ...object.Object) (object.Object, error) {
	if len(args) < 1 || len(args) > 64 {
		return nil, object.NewArgsRangeError("printf", 1, 64, len(args))
	}
	fs, err := object.AsString(args[0])
	if err != nil {
		return nil, err
	}
	fmtArgs := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		fmtArgs[i] = v.Interface()
	}
	b := []byte(fmt.Sprintf(fs, fmtArgs...))
	_, err = ren.GetOS(ctx).Stdout().Write(b)
	if err != nil {
		return nil, err
	}
	return object.Nil, nil
}

// Builtins returns a copy of the global built-in functions, keyed by name.
func Builtins() map[string]*object.Builtin {
	result := make(map[string]*object.Builtin, len(builtins))
	maps.Copy(result, builtins)
	return result
}
