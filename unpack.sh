#!/usr/bin/env bash

# Copy dot files too
shopt -s dotglob

ROOT=`pwd`
CURRENT_DIR_NAME=`basename "$ROOT"`
SRC="$ROOT/src/$CURRENT_DIR_NAME"
TMPDIR=`mktemp -d`

# Replace the import statements
if [ "$CURRENT_DIR_NAME" != "api-barebone" ]; then
  find . -name "*.go"|xargs sed -i "s/api-barebone/${CURRENT_DIR_NAME}/g"
fi

mv $ROOT/* $TMPDIR/
mkdir -p $SRC
mv $TMPDIR/* $SRC/
rm -rf $TMPDIR

echo "GOPATH=$ROOT" > $ROOT/.env

if [ "`which go 2>/dev/null`" != "" ]; then
  export GOPATH=$ROOT && go get -t ./...
fi
