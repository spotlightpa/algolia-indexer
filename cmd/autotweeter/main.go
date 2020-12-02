package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/spotlightpa/sourcesdb/pkg/autotweeter"
)

func main() {
	exitcode.Exit(autotweeter.CLI(os.Args[1:]))
}
