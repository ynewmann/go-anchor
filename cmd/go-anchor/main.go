package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-anchor <command> [args...]")
		fmt.Fprintln(os.Stderr, "Commands: idl (fetch, validate, convert)")
		os.Exit(1)
	}
	// TODO: implement idl fetch, validate, convert
	fmt.Printf("go-anchor %s (not yet implemented)\n", os.Args[1])
}
