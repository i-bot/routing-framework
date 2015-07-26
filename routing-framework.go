package main

import (
	"options"
	"os"
)

func main() {
	options.Interpret(os.Args[1:])
}
