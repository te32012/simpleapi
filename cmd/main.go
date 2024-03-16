package main

import (
	"fmt"
	"vktestgo2024/internal/app"
)

//go:generate swagger generate spec -o ./swagger.yml -m
func main() {
	fmt.Println("start")
	app.Run()
}
