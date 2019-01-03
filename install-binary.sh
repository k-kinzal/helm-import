#!/bin/sh

set -eu

DIR=`dirname $0`

VERSION="v0.1.0"
USER="k-kinzal"
REPO="helm-import"
OS=`uname | tr '[:upper:]' '[:lower:]'`
ARCH=`
case $(uname -m) in
  x86_64 ) echo "amd64";;
esac
`
URL="https://github.com/${USER}/${REPO}/releases/download/${VERSION}/${OS}-${ARCH}.tar.gz"

if which curl >/dev/null; then
  curl -sSL "$URL" | tar xvz -C $DIR >/dev/null
else
  wget -O - "$URL" | tar xvz -C $DIR >/dev/null
fi

chmod +x $DIR/import
