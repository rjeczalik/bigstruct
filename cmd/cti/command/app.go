package command

import (
	"github.com/spf13/pflag"
)

type App struct {
	Encoding string
}

func (app *App) Register(f *pflag.FlagSet) {
	f.StringVarP(&app.Encoding, "encoding", "e", "", "")
}
