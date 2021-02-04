package command

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/query"
	"github.com/rjeczalik/bigstruct/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App struct {
	context.Context

	Home   string
	Format string

	Config  Config
	Storage *storage.Gorm
	Query   *query.Query
}

func (app *App) Register(f *pflag.FlagSet) {
	f.StringVar(&app.Home, "home", DefaultHome(), "Default home directory")
	f.StringVarP(&app.Format, "format", "f", "text", "Output formatting type")
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

	app.Query = &query.Query{Storage: app.Storage}

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
	type textRenderer interface {
		RenderText(io.Writer)
	}

	switch app.Format {
	case "json":
		fmt.Println(types.MakePrettyJSON(v))
	case "yaml":
		fmt.Println(types.MakeYAML(v))
	case "text":
		switch v := v.(type) {
		case textRenderer:
			tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
			v.RenderText(tw)
			err = tw.Flush()
		default:
			fmt.Fprintln(os.Stdout, v)
		}
	default:
		err = fmt.Errorf("unsupported format type: %q", app.Format)
	}

	return err
}

func DefaultHome() string {
	if dir := os.Getenv("CONFETTI_HOME"); dir != "" {
		return dir
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return filepath.Join(dir, "bigstruct")
}
