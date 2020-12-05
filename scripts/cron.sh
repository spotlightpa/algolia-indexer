#!/bin/bash

set -euxo pipefail

curl -fsS -m 10 --retry 5 "$HC_URL/start"

OUTPUT=$(autotweeter 2>&1)
curl -fsS -m 10 --retry 5 --data-raw "$OUTPUT" "$HC_URL/$?"
