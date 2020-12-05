module github.com/spotlightpa/sourcesdb

go 1.15

// +heroku goVersion go1.15
// +heroku install ./cmd/...

require (
	github.com/algolia/algoliasearch-client-go/v3 v3.12.1
	github.com/carlmjohnson/exitcode v0.20.2
	github.com/carlmjohnson/flagext v0.20.2
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.6.0
	github.com/henvic/ctxsignal v1.0.0
	gocloud.dev v0.20.0
)
