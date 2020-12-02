package redis

import (
	"flag"

	"github.com/carlmjohnson/flagext"
	"github.com/gomodule/redigo/redis"
)

type Dialer = func() (redis.Conn, error)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Storable interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}) error
	GetSet(key string, getv, setv interface{}) error
	GetLock(key string) (unlock func(), err error)
}

// Var adds an option to the specified FlagSet (or flag.CommandLine if nil)
// that creates a Redis dialer for the specified URL.
// URLs should have the format redis://user:password@host:port/db
// where db is an integer and username is ignored.
// Use the callback after parsing options to retrieve the dialer.
func Var(fl *flag.FlagSet, name, usage string) func(Logger) (Storable, error) {
	if fl == nil {
		fl = flag.CommandLine
	}
	var d Dialer
	set := false
	flagext.Callback(fl, name, "", usage, func(s string) error {
		set = true
		var err error
		d, err = Parse(s)
		return err
	})
	return func(l Logger) (Storable, error) {
		if set {
			return New(d, l)
		}
		return NewMock("", "sourcesdb-redis", l), nil
	}
}
