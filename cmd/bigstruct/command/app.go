package command

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rjeczalik/bigstruct"
	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App struct {
	context.Context

	Home   string
	Format string
	Raw    bool

	Config  Config
	Storage *storage.Gorm
	Client  *bigstruct.Client
}

func (app *App) Register(f *pflag.FlagSet) {
	f.StringVar(&app.Home, "home", DefaultHome(), "Default home directory")
	f.StringVarP(&app.Format, "format", "f", "text", "Output formatting type")
	f.BoolVarP(&app.Raw, "raw", "r", false, "Print big.Struct representation")
}

func (app *App) DefaultHome(cmd *cobra.Command, home string) {
	if !cmd.Flags().Changed("home") {
		app.Home = home
	}
}

func (app *App) DefaultFormat(cmd *cobra.Command, format string) {
	if !cmd.Flags().Changed("format") {
		app.Format = format
	}
}

func (app *App) Init(*cobra.Command, []string) error {
	if err := os.MkdirAll(app.Home, 0755); err != nil {
		return err
	}

	cfg := filepath.Join(app.Home, "config.yaml")

	p, err := ioutil.ReadFile(cfg)
	switch {
	case os.IsNotExist(err):
		p = app.DefaultConfig().YAML().Bytes()

		if err := ioutil.WriteFile(cfg, p, 0644); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	if err := app.Config.FromYAML(types.YAML(p)); err != nil {
		return err
	}

	if app.Storage, err = storage.Open(app.Config.Backend.URI); err != nil {
		return err
	}

	app.Client = &bigstruct.Client{
		Transport: &bigstruct.Server{
			Storage: app.Storage,
		},
	}

	return nil
}

func (app *App) Close() (err error) {
	if app.Storage != nil {
		err = app.Storage.Close()
	}
	return err
}

func (app *App) DefaultConfig() Config {
	return Config{
		Backend: Endpoint{
			URI: (&url.URL{
				Scheme: "sqlite",
				Opaque: "file:" + filepath.Join(app.Home, "storage.db"),
				RawQuery: url.Values{
					"cache":    {"shared"},
					"_locking": {"EXCLUSIVE"},
					"_journal": {"OFF"},
					"_sync":    {"OFF"},
					"debug":    {"0"},
				}.Encode(),
			}).String(),
		},
	}
}

func (app *App) Render(v interface{}) (err error) {
	type fielder interface {
		Fields() big.Fields
	}

	type structer interface {
		Struct() big.Struct
	}

	if app.Raw {
		switch w := v.(type) {
		case fielder:
			v = w.Fields().Struct()
		case structer:
			v = w.Struct()
		}
	}

	switch app.Format {
	case "json":
		fmt.Println(types.MakePrettyJSON(v))
	case "yaml":
		fmt.Println(types.MakeYAML(v))
	case "text":
		switch v := v.(type) {
		case []byte:
			fmt.Fprintf(os.Stdout, "%s\n", bytes.TrimSpace(v))
		case nil:
			// skip
		default:
			fmt.Fprint(os.Stdout, v)
		}
	default:
		err = fmt.Errorf("unsupported format type: %q", app.Format)
	}

	return err
}

func DefaultHome() string {
	if dir := os.Getenv("BIGSTRUCT_HOME"); dir != "" {
		return dir
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return filepath.Join(dir, "bigstruct")
}
