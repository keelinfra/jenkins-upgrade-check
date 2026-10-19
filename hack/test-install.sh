#!/bin/sh
# Runs install.sh against release files in DIST, as the release workflow builds them: it
# must install the binary for this host, and must refuse one that does not match
# sha256sums.txt.
#
#   hack/test-install.sh DIST VERSION
set -eu

dist=$(cd "$1" && pwd)
version=$2
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
fail() { echo "test-install: FAIL: $*" >&2; exit 1; }

sh install.sh --base "file://$dist" --dir "$tmp/bin"
got=$("$tmp/bin/jenkins-upgrade-check" version)
[ "$got" = "jenkins-upgrade-check $version" ] || fail "installed binary says '$got'"

cp -R "$dist" "$tmp/bad"
for f in "$tmp"/bad/jenkins-upgrade-check_*; do printf x >>"$f"; done
if sh install.sh --base "file://$tmp/bad" --dir "$tmp/bad-bin" 2>"$tmp/err"; then
    fail "installed a binary that does not match its checksum"
fi
grep -q 'does not match its SHA-256' "$tmp/err" || fail "unexpected error: $(cat "$tmp/err")"
[ ! -e "$tmp/bad-bin/jenkins-upgrade-check" ] || fail "left a binary behind"
echo "test-install: ok ($got; a damaged download was refused)"
