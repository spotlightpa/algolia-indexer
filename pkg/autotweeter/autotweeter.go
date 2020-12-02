package autotweeter

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/carlmjohnson/flagext"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spotlightpa/sourcesdb/internal/redis"
)

const AppName = "autotweeter"

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

	fl.BoolVar(&app.mock, "mock", false, "mock calls rather than use real thing")

	accessToken := fl.String("twitter-access-token", "", "")
	accessTokenSecret := fl.String("twitter-access-token-secret", "", "")
	consumerKey := fl.String("twitter-consumer-key", "", "")
	consumerSecret := fl.String("twitter-consumer-secret", "", "")

	getRedis := redis.Var(fl, "redis-url", "`URL` for Redis (mock if not set)")

	app.l = log.New(nil, AppName+" ", log.LstdFlags)
	flagext.LoggerVar(
		fl, app.l, "silent", flagext.LogSilent, "don't log debug output")

	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `autotweeter - sends Tweets about people
Options:
`)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output(), "")
	}

	flagext.Callback(fl, "template", "", "Go-style `template` for Tweet text", func(s string) error {
		var err error
		app.tmpl, err = template.New("tweet").Parse(s)
		return err
	})

	src := flagext.FileOrURL(flagext.StdIO, nil)
	fl.Var(src, "src", "`file or URL` source for Tweets")
	app.src = src

	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagext.ParseEnv(fl, AppName); err != nil {
		return err
	}

	musthave := []string{"template"}

	if !app.mock {
		musthave = append(musthave,
			"twitter-access-token",
			"twitter-access-token-secret",
			"twitter-consumer-key",
			"twitter-consumer-secret",
			"redis-url")
	}

	if err := flagext.MustHave(fl, musthave...); err != nil {
		return err
	}

	if !app.mock {
		config := oauth1.NewConfig(*consumerKey, *consumerSecret)
		token := oauth1.NewToken(*accessToken, *accessTokenSecret)
		app.cl = config.Client(context.Background(), token)
	}
	var err error
	app.store, err = getRedis(app.l)

	return err
}

type appEnv struct {
	l     *log.Logger
	mock  bool
	cl    *http.Client
	store redis.Storable
	src   io.ReadCloser
	tmpl  *template.Template
}

func (app *appEnv) logf(format string, args ...interface{}) {
	app.l.Printf(format, args...)
}

func (app *appEnv) Exec() (err error) {
	app.logf("starting")

	// lock Redis
	unlock, err := app.store.GetLock(AppName)
	defer unlock()
	if err != nil {
		return err
	}

	app.logf("got lock")

	// get old tweets
	priorTweets := map[string]time.Time{}
	if err = app.store.Get("prior-tweets", &priorTweets); err != nil &&
		err != redis.ErrNil {
		return err
	}

	app.logf("found %d prior Tweets", len(priorTweets))

	// get list of contexts
	ctxs, err := app.getContext()
	if err != nil {
		return err
	}

	// remove items that were tweeted
	filteredCtxs := ctxs[:0]
	for _, ctx := range ctxs {
		if _, ok := priorTweets[ctx["id"]]; !ok {
			filteredCtxs = append(filteredCtxs, ctx)
		}
	}

	app.logf("found %d contexts, %d unused",
		len(ctxs), len(filteredCtxs))

	// random choice
	if !app.mock {
		rand.Seed(time.Now().UnixNano())
	}
	ctx := filteredCtxs[rand.Intn(len(filteredCtxs))]

	app.logf("chose %q", ctx["id"])

	// build text
	var buf strings.Builder
	if err = app.tmpl.Execute(&buf, ctx); err != nil {
		return err
	}

	// tweet
	if err = app.tweet(buf.String()); err != nil {
		return err
	}

	// update list of old Tweets
	priorTweets[ctx["id"]] = time.Now()
	if err = app.store.Set("prior-tweets", &priorTweets); err != nil {
		return err
	}

	return err
}

func (app *appEnv) getContext() ([]map[string]string, error) {
	dec := json.NewDecoder(app.src)
	var data []map[string]string
	err := dec.Decode(&data)
	return data, err
}

func (app *appEnv) tweet(text string) (err error) {
	app.logf("tweeting: %q\n", text)
	if app.mock {
		return nil
	}

	client := twitter.NewClient(app.cl)
	if _, _, err := client.Statuses.Update(text, nil); err != nil {
		return err
	}
	return nil
}
