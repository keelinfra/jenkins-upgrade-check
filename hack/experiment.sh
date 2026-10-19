#!/usr/bin/env bash
# Rehearse upgrading a typical controller and record what breaks.
#
# The "typical controller" runs FROM with the TOP most installed plugins
# (Jenkins usage statistics), each at the newest release FROM can load: the
# best case, an administrator who kept plugins current. Three rehearsals:
#
#   baseline  FROM with those plugins (checks the harness itself)
#   as-is     TO with the same plugin files: upgrading the core and nothing else
#   planned   TO with the plugin set `plan` computes
#
# Usage: hack/experiment.sh FROM [TO] [TOP] [MONTH]
# Results go to results/FROM-topTOP.
set -euo pipefail

from=$1
to=${2:-lts}
top=${3:-100}
month=${4:-202603}
bin=${BIN:-./jenkins-upgrade-check}
out=results/$from-top$top
mkdir -p "$out"

if [ "$to" = lts ]; then
  to=$(curl -fsSL https://updates.jenkins.io/stable/latestCore.txt)
fi

curl -fsSL "https://raw.githubusercontent.com/jenkins-infra/infra-statistics/gh-pages/jenkins-stats/svg/${month}-top-plugins1000.csv" -o "$out/top.csv"

"$bin" resolve --core "$from" --top-plugins "$out/top.csv" --top "$top" -o "$out/plugins-$from.txt"

"$bin" dry-run --core "$from" --plugins-file "$out/plugins-$from.txt" \
  --out-dir "$out" -o "$out/baseline.md" --json "$out/baseline.json" || true

"$bin" dry-run --core "$to" --plugins-file "$out/plugins-$from.txt" \
  --out-dir "$out/as-is" -o "$out/as-is.md" --json "$out/as-is.json" || true

"$bin" plan --plugins-file "$out/plugins-$from.txt" --core "$from" --to "$to" \
  -o "$out/plan.md" --json "$out/plan.json" --plugins-out "$out/plugins-$to.txt"

"$bin" dry-run --plan "$out/plan.json" \
  --out-dir "$out/planned" -o "$out/planned.md" --json "$out/planned.json" || true

summary() {
  jq -r --arg name "$1" '"| \($name) | \(.image) | \(if .started then "started" else "did not start" end) | \(.plugins) | \(.failed // [] | length) | \(.roots // [] | length) | \(.broken // [] | length) | \(.skipped // [] | length) |"' "$2"
}
causes() {
  jq -r --arg name "$1" '(.roots // [])[] | "- \($name): **\(.name) \(.version // "")** \(.diagnosis.kind): \(.diagnosis.why)\(if (.cascade // [] | length) > 0 then " (takes down \(.cascade | length) more: \(.cascade | join(", ")))" else "" end)"' "$2"
}
{
  echo "### $from → $to, $top most installed plugins"
  echo
  echo "| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |"
  echo "|---|---|---|---|---|---|---|---|"
  summary "baseline: $from" "$out/baseline.json"
  summary "as-is: $to, plugins unchanged" "$out/as-is.json"
  summary "planned: $to, plugins from plan" "$out/planned.json"
  echo
  causes baseline "$out/baseline.json"
  causes as-is "$out/as-is.json"
  causes planned "$out/planned.json"
  echo
} | tee "$out/summary.md"
