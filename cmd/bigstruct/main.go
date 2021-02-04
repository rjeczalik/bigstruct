package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/xq"

	"github.com/spf13/cobra"
)

func die(err error) {
	if _, ok := err.(fmt.Formatter); ok {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}

func main() {
	var (
		app = &command.App{
			Context: context.Background(),
		}
		cmd *cobra.Command
	)

	defer func() {
		if err := app.Close(); err != nil {
			die(err)
		}
	}()

	if filepath.Base(os.Args[0]) == "xq" {
		cmd = xq.NewCommand(app)
	} else {
		cmd = NewCommand(app)
	}

	if err := cmd.Execute(); err != nil {
		die(err)
	}
}
