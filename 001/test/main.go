package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/dop251/goja"
	"github.com/mymmrac/x"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	lua "github.com/yuin/gopher-lua"
)

const fileName = "hello"

func main() {
	mainGo()
	mainLua()
	mainJS()
}

// ==== // ==== //
// https://github.com/traefik/yaegi

func mainGo() {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	_, err := i.Eval(`import "os"`)
	x.Assert(err == nil, err)

	err = i.Use(map[string]map[string]reflect.Value{
		"my/": {
			"FileName": reflect.ValueOf(filenameGo),
		},
	})
	x.Assert(err == nil, err)
	i.ImportUsed()

	_, err = i.Eval(`file, err := os.Open(my.FileName("` + fileName + `"))`)
	x.Assert(err == nil, err)

	g := i.Globals()
	x.Assert(len(g) != 0)

	file := g["file"].Interface().(*os.File)
	if !g["err"].IsNil() {
		err = g["err"].Interface().(error)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := io.ReadAll(file)
	x.Assert(err == nil, err)
	fmt.Println(string(data))
}

func filenameGo(name string) string {
	return name + ".go.txt"
}

// ==== // ==== //
// https://github.com/yuin/gopher-lua

func mainLua() {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("file_name", L.NewFunction(filenameLua))

	err := L.DoString(`
local file = assert(io.open(file_name("` + fileName + `"), "r"))
content = file:read("*all")
file:close()
`)
	x.Assert(err == nil, err)

	fmt.Println(L.GetGlobal("content").String())
}

func filenameLua(L *lua.LState) int {
	name := L.ToString(1)
	L.Push(lua.LString(name + ".lua.txt"))
	return 1
}

// ==== // ==== //
// https://github.com/dop251/goja

func mainJS() {
	vm := goja.New()

	err := vm.Set("fileName", filenameJS)
	x.Assert(err == nil, err)

	v, err := vm.RunString(`fileName("` + fileName + `")`)
	x.Assert(err == nil, err)

	fileName := v.Export().(string)
	file, err := os.Open(fileName)
	x.Assert(err == nil, err)

	data, err := io.ReadAll(file)
	x.Assert(err == nil, err)
	fmt.Println(string(data))
}

func filenameJS(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	name := call.Argument(0).Export().(string)
	return vm.ToValue(name + ".js.txt")
}
