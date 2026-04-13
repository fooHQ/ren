// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package builtins

import (
	"bytes"
	"context"
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
	"print":    object.NewBuiltin("print", Print),
}

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

func Builtins() map[string]*object.Builtin {
	result := make(map[string]*object.Builtin, len(builtins))
	maps.Copy(result, builtins)
	return result
}
