package blob

import (
	"context"
	"flag"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Storable interface {
	Get(ctx context.Context, key string, v interface{}) error
	Set(ctx context.Context, key string, v interface{}) error
}

func Var(fl *flag.FlagSet, name, usage string) func(context.Context, Logger) (Storable, error) {
	if fl == nil {
		fl = flag.CommandLine
	}
	urlstr := fl.String(name, "", usage)
	return func(ctx context.Context, l Logger) (Storable, error) {
		if *urlstr != "" {
			return New(ctx, *urlstr, l)
		}
		return NewMock("", "sourcesdb-blob", l), nil
	}
}
