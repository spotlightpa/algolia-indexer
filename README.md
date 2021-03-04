# Algolia-Indexer

A tool to power your static site's search with [Algolia](https://www.algolia.com). Build a JSON site index and send it to Algolia with this script. To work, the Algolia script API key must have permission to create/drop an index because it creates a temporary table, sends all the data to the temporary table, then swaps it in.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN=$(pwd) go install github.com/spotlightpa/algolia-indexer@latest
```

## Usage

```
$ algolia-indexer -h
indexer - sends JSON object array to Algolia

Options may also be set as environment variables prefixed with INDEXER_.

Options:
  -algolia-api-key key
        key for Algolia API
  -algolia-app-id ID
        ID for Algolia app
  -algolia-index-name name
        name for Algolia index
  -src value
        source file or URL (default stdin)
  -verbose
        log debug output (default silent)
```