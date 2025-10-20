package main

import (
	"github.com/kyanite/focus/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
