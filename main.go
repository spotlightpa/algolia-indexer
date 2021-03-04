package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/spotlightpa/algolia-indexer/indexer"
)

func main() {
	exitcode.Exit(indexer.CLI(os.Args[1:]))
}
