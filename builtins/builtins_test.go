package builtins_test

import (
	"context"
	"testing"

	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/stretchr/testify/require"

	"github.com/foohq/ren"
	"github.com/foohq/ren/builtins"
	"github.com/foohq/ren/testutils"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name       string
		args       []object.Object
		wantOutput string
		wantErr    error
	}{
		{
			name: "string",
			args: []object.Object{
				object.NewString("hello world"),
			},
			wantOutput: "hello world\n",
		},
		{
			name:    "without args",
			wantErr: object.NewArgsError("print", 1, 0),
		},
		{
			name: "bytes",
			args: []object.Object{
				object.NewString("string"),
				object.NewBytes([]byte("bytes")),
				object.NewInt(32),
				object.False,
				object.True,
				object.Nil,
			},
			wantOutput: "string bytes 32 false true null\n",
		},
	}

	m := &testutils.MockOS{}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = ren.WithOS(ctx, m)

	stdout := ren.NewPipe()
	m.On("Stdout").Return(stdout)

	outputCh := make(chan string, 1)
	go func() {
		for ctx.Err() == nil {
			output := make([]byte, 100)
			n, _ := stdout.Read(output)
			outputCh <- string(output[:n])
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := builtins.Print(ctx, tt.args...)
			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantOutput, <-outputCh)
		})
	}

	cancel()
}

func TestPrintf(t *testing.T) {
	tests := []struct {
		name       string
		args       []object.Object
		wantOutput string
		wantErr    bool
	}{
		{
			name: "single string",
			args: []object.Object{
				object.NewString("hello world"),
			},
			wantOutput: "hello world",
		},
		{
			name: "string with args",
			args: []object.Object{
				object.NewString("hello %s, age %d"),
				object.NewString("world"),
				object.NewInt(32),
			},
			wantOutput: "hello world, age 32",
		},
		{
			name:    "without args",
			wantErr: true,
		},
	}

	m := &testutils.MockOS{}
	ctx := ren.WithOS(t.Context(), m)

	stdout := ren.NewPipe()
	m.On("Stdout").Return(stdout)

	outputCh := make(chan string, 1)
	go func() {
		for ctx.Err() == nil {
			output := make([]byte, 100)
			n, _ := stdout.Read(output)
			outputCh <- string(output[:n])
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := builtins.Printf(ctx, tt.args...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantOutput, <-outputCh)
		})
	}
}
