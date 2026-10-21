#!/usr/bin/env bash
# Commit the matrix for one target to the matrix branch, which keeps one
# directory per target version and a README that lists them.
#
# Usage: hack/publish-matrix.sh DIR TO
# DIR holds README.md (the table), matrix.json and one directory per row.
set -euo pipefail

src=$(cd "$1" && pwd)
to=$2
work=$(mktemp -d)

if git fetch -q --depth 1 origin matrix 2>/dev/null; then
  git worktree add -q --detach "$work" FETCH_HEAD
else
  git worktree add -q --detach "$work"
  git -C "$work" checkout -q --orphan matrix
  git -C "$work" rm -rqf --ignore-unmatch .
fi

cd "$work"
rm -rf "$to"
mkdir -p "$to"
cp -R "$src"/. "$to"/
{
  echo "# Jenkins upgrade matrix"
  echo
  echo "Every Jenkins LTS line, upgraded to a later release in a rehearsal by"
  echo "[jenkins-upgrade-check](https://github.com/keelinfra/jenkins-upgrade-check)."
  echo "One table per target version, newest first:"
  echo
  for d in $(find . -maxdepth 1 -type d -name '2.*' | sed 's|^\./||' | sort -t. -k2,2nr -k3,3nr); do
    echo "- [$d]($d/README.md)"
  done
} >README.md
git add -A
git -c user.name='github-actions[bot]' \
  -c user.email='41898282+github-actions[bot]@users.noreply.github.com' \
  -c commit.gpgsign=false \
  commit -q -m "Upgrade matrix for Jenkins $to"
git push -q origin HEAD:refs/heads/matrix
echo "published $to to the matrix branch"
