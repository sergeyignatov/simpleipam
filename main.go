package main

import (
	"fmt"
	"github.com/sergeyignatov/simpleipam/cmd"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("error: %v", err)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s", msg))
		os.Exit(1)
	}
}
