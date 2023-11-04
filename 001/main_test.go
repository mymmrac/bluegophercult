package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestEvalValue(t *testing.T) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	v, err := i.Eval(`1`)
	require.NoError(t, err)

	require.Equal(t, reflect.Int, v.Kind())
	require.Equal(t, int64(1), v.Int())
}

func TestEvalValue2(t *testing.T) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	v, err := i.Eval(`f, err := os.Open("")`)
	require.NoError(t, err)

	require.Equal(t, reflect.Int, v.Kind())
	require.Equal(t, int64(1), v.Int())
}

func TestEvalMultiple(t *testing.T) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	_, err := i.Eval(`
package x

var X = "Hello\n"
`)
	require.NoError(t, err)

	_, err = i.Eval(`
package main

import "x"

func main() {

	print(x.X)
}
`)
	require.NoError(t, err)
}
