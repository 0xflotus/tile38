#!/bin/bash
set -e

VERSION="1.0.4"
PROTECTED_MODE="no"

# Hardcode some values to the core package
LDFLAGS="$LDFLAGS -X github.com/tidwall/tile38/core.Version=${VERSION}"
if [ -d ".git" ]; then
	LDFLAGS="$LDFLAGS -X github.com/tidwall/tile38/core.GitSHA=$(git rev-parse --short HEAD)"
fi
LDFLAGS="$LDFLAGS -X github.com/tidwall/tile38/core.BuildTime=$(date +%FT%T%z)"
if [ "$PROTECTED_MODE" == "no" ]; then
	LDFLAGS="$LDFLAGS -X github.com/tidwall/tile38/core.ProtectedMode=no"
fi

export GO15VENDOREXPERIMENT=1

cd $(dirname "${BASH_SOURCE[0]}")
OD="$(pwd)"

# temp directory for storing isolated environment.
TMP="$(mktemp -d -t tile38.XXXX)"
function rmtemp {
  	rm -rf "$TMP"
}
trap rmtemp EXIT

if [ "$NOCOPY" != "1" ]; then
	# copy all files to an isloated directory.
	WD="$TMP/src/github.com/tidwall/tile38"
	GOPATH="$TMP"
	for file in `find . -type f`; do
		# TODO: use .gitignore to ignore, or possibly just use git to determine the file list.
		if [[ "$file" != "." && "$file" != ./.git* && "$file" != ./data* && "$file" != ./tile38-* ]]; then
			mkdir -p "$WD/$(dirname "${file}")"
			cp -P "$file" "$WD/$(dirname "${file}")"
		fi
	done
	cd $WD
fi

#core/gen.sh

# build and store objects into original directory.
go build -ldflags "$LDFLAGS" -o "$OD/tile38-server" cmd/tile38-server/*.go
go build -ldflags "$LDFLAGS" -o "$OD/tile38-cli" cmd/tile38-cli/*.go

# test if requested
if [ "$1" == "test" ]; then
	$OD/tile38-server -p 9876 -d "$TMP" -q &
	PID=$!
	function testend {
	  	kill $PID &
	}
	trap testend EXIT
	go test $(go list ./... | grep -v /vendor/)
fi

