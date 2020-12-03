package main

import (
	"fmt"
	"os"

	"github.com/glaucusio/confetti/cmd/cti/command"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	app := new(command.App)
	cmd := command.New(app)

	if err := cmd.Execute(); err != nil {
		die(err)
	}
}
