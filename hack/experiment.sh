#!/usr/bin/env bash
# Rehearse upgrading a typical controller and record what breaks.
#
# The "typical controller" runs FROM with the TOP most installed plugins
# (Jenkins usage statistics), each at the newest release FROM can load: the
# best case, an administrator who kept plugins current. The rehearsals:
#
#   baseline  FROM with those plugins (checks the starting point itself)
#   as-is     TO with the same plugin files: upgrading the core and nothing else
#   stop-N    each core upgrade of the plan, with the plugin set `plan`
#             computes for it; the last stop is TO
#
# Usage: hack/experiment.sh FROM [TO] [TOP] [MONTH]
# Results go to results/FROM-topTOP: a Markdown summary and row.json, which
# hack/matrix.sh turns into the table for every starting version.
set -euo pipefail

from=$1
to=${2:-lts}
top=${3:-100}
month=${4:-$(go run ./hack/matrix month)}
bin=${BIN:-./jenkins-upgrade-check}
out=results/$from-top$top
mkdir -p "$out"

if [ "$to" = lts ]; then
  to=$(curl -fsSL https://updates.jenkins.io/stable/latestCore.txt)
fi

curl -fsSL "https://raw.githubusercontent.com/jenkins-infra/infra-statistics/gh-pages/jenkins-stats/svg/${month}-top-plugins1000.csv" -o "$out/top.csv"

"$bin" resolve --core "$from" --top-plugins "$out/top.csv" --top "$top" -o "$out/plugins-$from.txt"

"$bin" dry-run --core "$from" --plugins-file "$out/plugins-$from.txt" \
  --out-dir "$out/baseline" -o "$out/baseline.md" --json "$out/baseline.json" \
  2> >(tee "$out/baseline.err" >&2) || true

"$bin" dry-run --core "$to" --from "$from" --plugins-file "$out/plugins-$from.txt" \
  --out-dir "$out/as-is" -o "$out/as-is.md" --json "$out/as-is.json" \
  2> >(tee "$out/as-is.err" >&2) || true

"$bin" plan --plugins-file "$out/plugins-$from.txt" --core "$from" --to "$to" \
  -o "$out/plan.md" --json "$out/plan.json" --plugins-out "$out/plugins-$to.txt"

# Step 0 updates plugins on FROM, which the baseline already did; every later
# step is a core upgrade.
steps=$(jq '.steps | length' "$out/plan.json")
for ((i = 1; i < steps; i++)); do
  "$bin" dry-run --plan "$out/plan.json" --step "$i" \
    --out-dir "$out/stop-$i" -o "$out/stop-$i.md" --json "$out/stop-$i.json" \
    2> >(tee "$out/stop-$i.err" >&2) || true
done

# A full disk makes every plugin after it fail; that says nothing about the
# upgrade, so the row is not published.
if grep -qs "No space left on device" "$out"/*/dry-run-*.log "$out"/*.err; then
  echo "the runner ran out of disk space; these results are not usable" >&2
  df -h / >&2
  exit 1
fi

# brief turns a rehearsal into the fields the table needs; a rehearsal that
# could not run at all (no image, say) has no JSON, only its last error.
brief() {
  local json=$1.json err=$1.err
  if [ -s "$json" ]; then
    jq -c '{core, java, javaReason, image, started, reason, plugins, upgrading, detached,
      failed: (.failed // [] | length), broken: (.broken // [] | length),
      brokenNames: [(.broken // [])[] | .name],
      skipped: (.skipped // [] | length), skippedWhy: (.skipped // [])[:5],
      roots: [(.roots // [])[] | {name, version, kind: .diagnosis.kind, why: .diagnosis.why, cascade: (.cascade // [] | length)}]}' "$json"
  else
    jq -cn --arg core "$2" --arg why "$(tail -n 1 "$err" 2>/dev/null)" \
      '{core: $core, started: false, reason: (if $why == "" then "no result" else $why end), roots: []}'
  fi
}

stops="[]"
for ((i = 1; i < steps; i++)); do
  core=$(jq -r ".steps[$i].core" "$out/plan.json")
  java=$(jq ".steps[$i].java // 0" "$out/plan.json")
  stops=$(jq -c --argjson s "$(brief "$out/stop-$i" "$core")" --argjson java "$java" '. + [$s + {planJava: $java}]' <<<"$stops")
done
jq -n --arg from "$from" --arg to "$to" --argjson top "$top" --arg month "$month" \
  --argjson baseline "$(brief "$out/baseline" "$from")" \
  --argjson asis "$(brief "$out/as-is" "$to")" \
  --argjson stops "$stops" \
  --argjson problems "$(jq '[.steps[].problems // [] | length] | add' "$out/plan.json")" \
  '{from: $from, to: $to, top: $top, month: $month, baseline: $baseline, asIs: $asis, stops: $stops, planProblems: $problems}' \
  >"$out/row.json"

summary() {
  jq -r --arg name "$1" '"| \($name) | \(.image // "") | \(if .started then "started" else "did not start" end) | \(.plugins // 0) | \(.failed // 0) | \(.roots | length) | \(.broken // 0) | \(.skipped // 0) |"'
}
causes() {
  jq -r --arg name "$1" '(.roots[] | "- \($name): **\(.name) \(.version // "")** \(.kind): \(.why)\(if .cascade > 0 then " (takes down \(.cascade) more)" else "" end)"),
    (select((.brokenNames // []) != []) | "- \($name): loaded in part (extensions failed): \(.brokenNames | join(", "))"),
    (select(.javaReason) | "- \($name): Java \(.java) because \(.javaReason)"),
    (select((.detached // []) != []) | "- \($name): upgrading from \(.upgrading), Jenkins installed from its war: \(.detached | join(", "))"),
    ((.skippedWhy // [])[] | "- \($name): not downloaded: \(.)"),
    (select(.started | not) | "- \($name): \(.reason // "did not start")")'
}
{
  echo "### $from → $to, $top most installed plugins"
  echo
  echo "| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |"
  echo "|---|---|---|---|---|---|---|---|"
  jq '.baseline' "$out/row.json" | summary "baseline: $from"
  jq '.asIs' "$out/row.json" | summary "as-is: $to, plugins unchanged"
  jq -c '.stops[]' "$out/row.json" | while read -r s; do
    summary "planned: $(jq -r .core <<<"$s"), plugins from plan" <<<"$s"
  done
  echo
  jq '.baseline' "$out/row.json" | causes baseline
  jq '.asIs' "$out/row.json" | causes as-is
  jq -c '.stops[]' "$out/row.json" | while read -r s; do
    causes "planned $(jq -r .core <<<"$s")" <<<"$s"
  done
  echo
} | tee "$out/summary.md"
