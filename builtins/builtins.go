package builtins

import (
	modbuiltins "github.com/risor-io/risor/builtins"
	modfmt "github.com/risor-io/risor/modules/fmt"
	"github.com/risor-io/risor/object"
)

var builtins = map[string]object.Object{
	"all":         object.NewBuiltin("all", modbuiltins.All),
	"any":         object.NewBuiltin("any", modbuiltins.Any),
	"assert":      object.NewBuiltin("assert", modbuiltins.Assert),
	"bool":        object.NewBuiltin("bool", modbuiltins.Bool),
	"buffer":      object.NewBuiltin("buffer", modbuiltins.Buffer),
	"byte_slice":  object.NewBuiltin("byte_slice", modbuiltins.ByteSlice),
	"byte":        object.NewBuiltin("byte", modbuiltins.Byte),
	"call":        object.NewBuiltin("call", modbuiltins.Call),
	"chan":        object.NewBuiltin("chan", modbuiltins.Chan),
	"chr":         object.NewBuiltin("chr", modbuiltins.Chr),
	"chunk":       object.NewBuiltin("chunk", modbuiltins.Chunk),
	"close":       object.NewBuiltin("close", modbuiltins.Close),
	"coalesce":    object.NewBuiltin("coalesce", modbuiltins.Coalesce),
	"decode":      object.NewBuiltin("decode", modbuiltins.Decode),
	"delete":      object.NewBuiltin("delete", modbuiltins.Delete),
	"encode":      object.NewBuiltin("encode", modbuiltins.Encode),
	"error":       object.NewBuiltin("error", modbuiltins.Error),
	"float_slice": object.NewBuiltin("float_slice", modbuiltins.FloatSlice),
	"float":       object.NewBuiltin("float", modbuiltins.Float),
	"getattr":     object.NewBuiltin("getattr", modbuiltins.GetAttr),
	"hash":        object.NewBuiltin("hash", modbuiltins.Hash),
	"int":         object.NewBuiltin("int", modbuiltins.Int),
	"is_hashable": object.NewBuiltin("is_hashable", modbuiltins.IsHashable),
	"iter":        object.NewBuiltin("iter", modbuiltins.Iter),
	"keys":        object.NewBuiltin("keys", modbuiltins.Keys),
	"len":         object.NewBuiltin("len", modbuiltins.Len),
	"list":        object.NewBuiltin("list", modbuiltins.List),
	"make":        object.NewBuiltin("make", modbuiltins.Make),
	"map":         object.NewBuiltin("map", modbuiltins.Map),
	"ord":         object.NewBuiltin("ord", modbuiltins.Ord),
	"reversed":    object.NewBuiltin("reversed", modbuiltins.Reversed),
	"set":         object.NewBuiltin("set", modbuiltins.Set),
	"sorted":      object.NewBuiltin("sorted", modbuiltins.Sorted),
	"sprintf":     object.NewBuiltin("sprintf", modbuiltins.Sprintf),
	"string":      object.NewBuiltin("string", modbuiltins.String),
	"try":         object.NewBuiltin("try", modbuiltins.Try),
	"type":        object.NewBuiltin("type", modbuiltins.Type),
	"print":       object.NewBuiltin("print", modfmt.Println),
	"printf":      object.NewBuiltin("printf", modfmt.Printf),
	"errorf":      object.NewBuiltin("errorf", modfmt.Errorf),
}

func Builtins() []string {
	result := make([]string, 0, len(builtins))
	for name := range builtins {
		result = append(result, name)
	}
	return result
}

func Globals() map[string]any {
	result := make(map[string]any)
	for name, fn := range builtins {
		result[name] = fn
	}
	return result
}
