package main

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/cli/root"
	_ "github.com/vj396/milton/src/cli/run"
)

func main() {
	err := root.GetRoot().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not execute cli, er: %+v", err)
		os.Exit(1)
	}
}
