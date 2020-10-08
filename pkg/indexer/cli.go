package indexer

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/carlmjohnson/flagext"
)

const AppName = "indexer"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}
	if err = app.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)

	src := flagext.File(flagext.StdIO)
	app.src = src
	fl.Var(src, "src", "source file or URL")

	app.l = log.New(nil, AppName+" ", log.LstdFlags)
	flagext.LoggerVar(
		fl, app.l, "verbose", flagext.LogVerbose, "log debug output")

	appID := fl.String("algolia-app-id", "", "`ID` for Algolia app")
	apiKey := fl.String("algolia-api-key", "", "`key` for Algolia API")
	indexName := fl.String("algolia-index-name", "", "`name` for Algolia index")

	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `indexer - sends JSON object array to Algolia

Options:
`)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output(), "")
	}
	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagext.ParseEnv(fl, AppName); err != nil {
		return err
	}

	if err := flagext.MustHave(fl, "algolia-app-id", "algolia-api-key", "algolia-index-name"); err != nil {
		return err
	}

	client := search.NewClient(*appID, *apiKey)
	app.i = client.InitIndex(*indexName)

	return nil
}

type appEnv struct {
	src io.ReadCloser
	i   *search.Index
	l   *log.Logger
}

func (app *appEnv) logf(format string, args ...interface{}) {
	app.l.Printf(format, args...)
}

func (app *appEnv) Exec() (err error) {
	app.logf("starting")

	raw, err := ioutil.ReadAll(app.src)
	if err != nil {
		return err
	}

	app.logf("decoding source")
	var objs []map[string]interface{}
	if err = json.Unmarshal(raw, &objs); err != nil {
		return err
	}
	app.logf("saving to Algolia")
	res, err := app.i.SaveObjects(objs)
	if err != nil {
		return err
	}
	app.logf("saved to %d objects to Algolia", len(res.ObjectIDs()))

	return err
}
