#!/bin/sh
# Installs jenkins-upgrade-check: downloads the binary for this OS and CPU from a GitHub
# release, checks it against the release's sha256sums.txt, and installs it.
#
#   curl -fsSL https://github.com/keelinfra/jenkins-upgrade-check/releases/latest/download/install.sh | sh
#   curl -fsSL .../install.sh | sh -s -- --dir ~/bin --version 0.1.0
#
# Options:
#   --version V   release to install, e.g. 0.1.0 (default: the latest)
#   --dir DIR     where to install (default: /usr/local/bin when writable, else ~/.local/bin)
#   --base URL    take the release files from URL/<file> instead of GitHub, e.g. a mirror
#
# Linux and macOS, amd64 and arm64. On Windows, download jenkins-upgrade-check_windows_amd64.exe
# from the release page. Nothing is installed unless the checksum matches.
set -eu

REPO=https://github.com/keelinfra/jenkins-upgrade-check
NAME=jenkins-upgrade-check
VERSION='' DIR='' BASE=''
while [ $# -gt 0 ]; do
    case $1 in
        --version) VERSION=${2#v}; shift 2 ;;
        --dir) DIR=$2; shift 2 ;;
        --base) BASE=${2%/}; shift 2 ;;
        -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
        *) echo "$NAME: unknown option $1 (see --help)" >&2; exit 2 ;;
    esac
done

die() { echo "$NAME: $*" >&2; exit 1; }
say() { echo "$NAME: $*"; }
have() { command -v "$1" >/dev/null 2>&1; }

case $(uname -s) in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) die "no build for $(uname -s); download one from $REPO/releases" ;;
esac
case $(uname -m) in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) die "no build for $(uname -m); download one from $REPO/releases" ;;
esac
asset=${NAME}_${os}_$arch

if [ -z "$BASE" ]; then
    if [ -n "$VERSION" ]; then BASE=$REPO/releases/download/v$VERSION
    else BASE=$REPO/releases/latest/download
    fi
fi
if [ -z "$DIR" ]; then
    if [ -w /usr/local/bin ]; then DIR=/usr/local/bin; else DIR=$HOME/.local/bin; fi
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
fetch() { # url file
    if have curl; then curl -fsSL -o "$2" "$1"
    elif have wget; then wget -q -O "$2" "$1"
    else die "needs curl or wget"
    fi
}
sha256() { # file -> its SHA-256
    if have sha256sum; then sha256sum < "$1" | cut -d' ' -f1
    elif have shasum; then shasum -a 256 < "$1" | cut -d' ' -f1
    else die "needs sha256sum or shasum to check the download"
    fi
}

say "downloading $asset from $BASE"
fetch "$BASE/$asset" "$work/$NAME" || die "cannot download $BASE/$asset"
fetch "$BASE/sha256sums.txt" "$work/sha256sums.txt" || die "cannot download $BASE/sha256sums.txt"
want=$(awk -v f="$asset" '$2 == f || $2 == "*" f { print $1 }' "$work/sha256sums.txt")
[ -n "$want" ] || die "$asset is not listed in sha256sums.txt"
[ "$(sha256 "$work/$NAME")" = "$want" ] || die "$asset does not match its SHA-256 in sha256sums.txt (download damaged or altered); nothing installed"

mkdir -p "$DIR"
# Copy next to the target, try it there (/tmp may be noexec), then rename it over the old
# one, so a half-written or broken file is never left in its place.
new=$DIR/.$NAME.new
cp "$work/$NAME" "$new"
chmod 755 "$new"
version=$("$new" version) || { rm -f "$new"; die "the downloaded $asset does not run here"; }
mv "$new" "$DIR/$NAME"
say "installed $version as $DIR/$NAME"
case :$PATH: in
    *:"$DIR":*) ;;
    *) say "$DIR is not on your PATH; add it, or run $DIR/$NAME" ;;
esac
