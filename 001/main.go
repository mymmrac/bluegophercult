package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/cjoudrey/gluahttp"
	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/x"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/go", yaegi)
	app.Post("/go-template", goTemplate)
	app.Post("/lua", gopherLua)

	const port = "8080"
	err := app.Listen(":" + port)
	x.Assert(err == nil, err)
}

func yaegi(c *fiber.Ctx) error {
	i := interp.New(interp.Options{
		Stdout: c.Response().BodyWriter(),
	})
	err := i.Use(stdlib.Symbols)
	if err != nil {
		return fmt.Errorf("use std: %w", err)
	}

	_, err = i.Eval(string(c.Body()))
	if err != nil {
		return fmt.Errorf("eval: %T: %w", err, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func gopherLua(c *fiber.Ctx) error {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	write := func(L *lua.LState) int {
		lv := L.Get(1)
		if str, ok := lv.(lua.LString); ok {
			_, err := c.WriteString(string(str))
			if err != nil {
				fmt.Println("ERROR: ", err)
				return 0
			}
		} else {
			fmt.Println("ERROR: unknown type:", lv)
		}
		return 0
	}
	L.SetGlobal("write", L.NewFunction(write))

	if err := L.DoString(string(c.Body())); err != nil {
		return fmt.Errorf("eval: %T: %w", err, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func goTemplate(c *fiber.Ctx) error {
	t := template.New("")

	t = t.Funcs(template.FuncMap{
		"getPage": http.Get,
		"readAll": func(r *http.Response) (string, error) {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	})

	t, err := t.Parse(string(c.Body()))
	if err != nil {
		return fmt.Errorf("parse template: %T: %w", err, err)
	}

	err = t.Execute(c, nil)
	if err != nil {
		return fmt.Errorf("execute template: %T: %w", err, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
