//go:build go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func Curl() {
	href := "https://google.com"

	resp, err := http.Get(href)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}

func main() {
	Curl()
}
