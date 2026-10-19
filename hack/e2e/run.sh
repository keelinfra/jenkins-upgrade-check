#!/usr/bin/env bash
# Upgrade a real controller by running the runbook the tool writes, block by
# block, and check that its users get back what they had.
#
#   hack/e2e/run.sh SCENARIO [FROM] [TO]
#
#   go              the runbook as written: every verify must say GO, and
#                   afterwards the jobs, their build history and the agent
#                   must be there and builds must pass on the new version.
#   forgot-plugins  at the first stop the operator forgets to copy the new
#                   plugins: verify must say ROLL BACK, and the runbook's
#                   rollback must bring back the controller as it was.
#   missing-plugin  the plugin pipeline jobs need disappears during the
#                   upgrade: verify must say ROLL BACK and name it, and the
#                   rollback must bring back the controller as it was.
#   old-agent       the agent runs Java 17, which the target does not
#                   support: the plan must name it, verify must say ROLL
#                   BACK and why, and the rollback must bring it back.
#
# The controller starts on FROM (default 2.346.3) with the TOP (default 300)
# most installed plugins, configured by seed.groovy, with one inbound agent,
# and runs every job once before the upgrade. Runs as root, as the runbook
# says (sudo -E env "PATH=$PATH" hack/e2e/run.sh ...); needs Docker with
# compose, curl, jq and jenkins-upgrade-check on PATH; takes over the
# container names jenkins and agent1, the volume jenkins_home and port 8080.
set -euo pipefail

scenario=${1:?scenario: go, forgot-plugins, missing-plugin or old-agent}
from=${2:-2.346.3}
to=${3:-lts}
top=${TOP:-300}
here=$(cd "$(dirname "$0")" && pwd)
work=${E2E_WORK:-$PWD/e2e-$scenario}
from_image=${FROM_IMAGE:-jenkins/jenkins:$from}
agent_image=jenkins/inbound-agent:3384.v60d89463d9e0-1-jdk21@sha256:6274f49f86429c5382b3943e96df81b65fd81712677347ae2a50b235b6c4e286
if [ "$scenario" = old-agent ]; then
  agent_image=jenkins/inbound-agent:3384.v60d89463d9e0-1-jdk17@sha256:b0c41a6b69468ac7c72ad25c2a364e6b5f116a1ea44510642bd629a6323f0fab
fi
url=http://localhost:8080
home=/var/lib/docker/volumes/jenkins_home/_data
password=e2e-admin
seed_plugins=cloudbees-folder,credentials,plain-credentials,credentials-binding,matrix-auth,workflow-aggregator,pipeline-model-definition,timestamper,git,email-ext,ws-cleanup,build-timeout

case $scenario in go | forgot-plugins | missing-plugin | old-agent) ;; *)
  echo "unknown scenario $scenario" >&2
  exit 2
  ;;
esac
[ "$(id -u)" = 0 ] || {
  echo "run as root, as the runbook says: sudo -E env \"PATH=\$PATH\" $0 $*" >&2
  exit 2
}

rm -rf "$work"
mkdir -p "$work/blocks"
cd "$work"
summary=$work/summary.md
failures=()
t0=$(date +%s)

log() { printf '%s  %s\n' "$(date -u +%H:%M:%S)" "$*" | tee -a timeline.txt; }
fail() {
  failures+=("$*")
  log "FAIL: $*"
}
wait_up() {
  for _ in $(seq 120); do
    curl -fs -o /dev/null "$url/login" && return 0
    sleep 5
  done
  return 1
}
api() { curl -fsS -u "admin:$TOKEN" "$@"; }
disk_free() { df -h --output=avail / | tail -1 | tr -d ' '; }
core_now() { curl -fs -o /dev/null -D - "$url/login" | tr -d '\r' | awk -F': ' 'tolower($1)=="x-jenkins"{print $2}'; }

# Every job with its builds, one line per build: "full/name #n RESULT".
tree='jobs[fullName,builds[number,result],jobs[fullName,builds[number,result],jobs[fullName,builds[number,result]]]]'
builds() {
  api -g "$url/api/json?tree=$tree" |
    jq -r '.. | objects | select(has("fullName")) | .fullName as $j | (.builds // [])[] | "\($j) #\(.number) \(.result)"' | sort
}
list_jobs() {
  api -g "$url/api/json?tree=$tree" | jq -r '.. | objects | select(has("builds")) | .fullName' | sort
}
job_url() { printf '%s/job/%s' "$url" "${1//\//\/job\/}"; }

# Runs every job once and waits for it; prints "name result".
run_jobs() {
  local j u n params r
  for j in $(list_jobs); do
    u=$(job_url "$j")
    n=$(api "$u/api/json?tree=nextBuildNumber" | jq -r .nextBuildNumber)
    params=$(api -g "$u/api/json?tree=property[parameterDefinitions[name]]" | jq '[.property[]?.parameterDefinitions[]?] | length')
    if [ "$params" -gt 0 ]; then
      api -X POST "$u/buildWithParameters" -o /dev/null
    else
      api -X POST "$u/build" -o /dev/null
    fi
    r=""
    for _ in $(seq 100); do
      r=$(api "$u/$n/api/json?tree=building,result" 2>/dev/null | jq -r 'if .building then "" else .result end' || true)
      [ -n "$r" ] && break
      sleep 3
    done
    if [ -z "$r" ]; then
      r="did-not-finish($(api -g "$url/queue/api/json?tree=items[task[url],why]" | jq -r --arg u "$u/" '[.items[] | select(.task.url | endswith($u)) | .why] | first // "not queued"' | tr ' ' '_'))"
    fi
    echo "$j $r"
  done
}

agent_online() {
  for _ in $(seq "${1:-24}"); do
    [ "$(api "$url/computer/agent1/api/json?tree=offline" | jq -r .offline)" = false ] && return 0
    sleep 5
  done
  return 1
}

finish() {
  local rc=$?
  {
    echo "## Scenario: $scenario"
    echo
    echo "Jenkins $from → ${target:-$to}, the $top most installed plugins, ${hops:-?} on the way. $(($(date +%s) - t0)) s in all."
    echo
    if [ ${#failures[@]} -eq 0 ] && [ $rc -eq 0 ]; then
      echo "**Result: as expected.**"
    else
      echo "**Result: NOT as expected.**"
      for f in "${failures[@]}"; do echo "- $f"; done
      if [ $rc -ne 0 ]; then echo "- the script stopped with exit code $rc (see the job log)"; fi
    fi
    echo
    if [ -n "${windows:-}" ]; then printf 'Window for users (quiet down → open again):%s\n\n' "$windows"; fi
    echo '```'
    cat timeline.txt 2>/dev/null
    echo '```'
    echo
    for f in /var/tmp/jenkins-upgrade/rehearsal.md /var/tmp/jenkins-upgrade/verify-*.md; do
      [ -f "$f" ] || continue
      echo "<details><summary>$(basename "$f")</summary>"
      echo
      sed 's/^#/###/' "$f"
      echo
      echo "</details>"
      echo
    done
  } >"$summary"
  docker logs jenkins >jenkins.log 2>&1 || true
  docker logs agent1 >agent.log 2>&1 || true
  if [ $rc -eq 0 ] && [ ${#failures[@]} -gt 0 ]; then exit 1; fi
  exit $rc
}
trap finish EXIT

# --- The controller before the upgrade ------------------------------------
log "setting up Jenkins $from ($from_image) with the $top most installed plugins; $(disk_free) free on /"
docker rm -f jenkins agent1 >/dev/null 2>&1 || true
docker volume rm jenkins_home >/dev/null 2>&1 || true
rm -rf /var/tmp/jenkins-upgrade /var/backups/jenkins-before-upgrade.tgz
cp "$here/compose.yaml" compose.yaml
echo "JENKINS_IMAGE=$from_image" >.env

curl -fsSL "https://raw.githubusercontent.com/jenkins-infra/infra-statistics/gh-pages/jenkins-stats/svg/202603-top-plugins1000.csv" -o top.csv
jenkins-upgrade-check resolve --core "$from" --top-plugins top.csv --top "$top" -o top.txt
names=$( (cut -d: -f1 top.txt && tr , '\n' <<<"$seed_plugins") | sort -u | paste -sd,)
jenkins-upgrade-check resolve --core "$from" --names "$names" -o plugins.txt
mkdir -p plugins
tr : ' ' <plugins.txt | xargs -P 8 -n 2 sh -c 'curl -fsSL --retry 3 -o "plugins/$0.jpi" "https://updates.jenkins.io/download/plugins/$0/$1/$0.hpi" || { echo "could not download $0 $1" >&2; exit 1; }'

docker compose up --no-start --quiet-pull 2>&1 | tail -2
mkdir -p "$home/plugins" "$home/init.groovy.d"
cp plugins/*.jpi "$home/plugins/"
cp "$here/seed.groovy" "$home/init.groovy.d/seed.groovy"
chown -R 1000:1000 "$home"
docker compose up -d
wait_up || {
  fail "Jenkins $from did not start"
  exit 1
}
log "Jenkins $(core_now) is up with $(wc -l <plugins.txt) plugins"

curl -fsS -c cookies -u "admin:$password" "$url/crumbIssuer/api/json" >crumb.json
TOKEN=$(curl -fsS -b cookies -u "admin:$password" -H "$(jq -r '.crumbRequestField + ": " + .crumb' crumb.json)" \
  --data newTokenName=e2e "$url/user/admin/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken" | jq -r .data.tokenValue)
export TOKEN
secret=$(api --data-urlencode 'script=print(jenkins.model.Jenkins.get().getComputer("agent1").getJnlpMac())' "$url/scriptText")
docker run -d --name agent1 --network ci --restart unless-stopped \
  -e JENKINS_URL=http://jenkins:8080 -e JENKINS_SECRET="$secret" -e JENKINS_AGENT_NAME=agent1 \
  -e JENKINS_AGENT_WORKDIR=/home/jenkins/agent "$agent_image" >/dev/null
if agent_online; then log "agent1 is online"; else log "agent1 did not connect to $from (see agent.log)"; fi

run_jobs | tee first-builds.txt
run_jobs >/dev/null
if grep -v ' SUCCESS$' first-builds.txt >/dev/null; then fail "builds failed on $from before the upgrade: $(grep -v ' SUCCESS$' first-builds.txt | paste -sd' ')"; fi
builds >builds-before.txt
log "$(list_jobs | wc -l) jobs, $(wc -l <builds-before.txt) builds before the upgrade"

# --- What the administrator does -----------------------------------------
export JENKINS_URL=$url JENKINS_USER=admin JENKINS_TOKEN=$TOKEN
jenkins-upgrade-check collect -o snapshot.json
jenkins-upgrade-check plan --snapshot snapshot.json --to "$to" -o plan.md --json plan.json --runbook runbook.md --install docker
target=$(jq -r .to plan.json)
hops=$(jq -r '[.steps[] | select(.upgrade) | .core] | join(" → ")' plan.json)
first_hop=$(jq -r '[.steps[] | select(.upgrade) | .core][0]' plan.json)
log "plan: $from → $hops"
java_seen=$(jq -r '.agents[]? | select(.name == "agent1") | .java // empty' snapshot.json)
[ -n "$java_seen" ] || fail "collect did not read the Java of agent1"
at_risk=$(jq -r '[.agents[]?.name] | join(" ")' plan.json)
case $scenario in
old-agent) [ "$at_risk" = agent1 ] || fail "the plan does not name agent1 (Java $java_seen) as unable to follow" ;;
*) [ -z "$at_risk" ] || fail "the plan says these agents cannot follow: $at_risk" ;;
esac

# Split the runbook into its shell blocks, each with its section and lead.
awk -v dir=blocks '
  /^## / { sec = substr($0, 4); lead = "" }
  /^### / { lead = substr($0, 5) }
  /^(\*\*|[0-9]+\. \*\*)/ { lead = $0; sub(/^[0-9]+\. /, "", lead); sub(/^\*\*/, "", lead); sub(/\*\*.*$/, "", lead); sub(/[.:]$/, "", lead) }
  /^```sh$/ { n++; f = sprintf("%s/%02d.sh", dir, n); printf "%02d\t%s\t%s\n", n, sec, lead > dir "/index.tsv"; inb = 1; next }
  /^```$/ { if (inb) { inb = 0; close(f) } next }
  inb { print > f }
' runbook.md

# verify and dry-run report through their exit code; the runbook's reader
# acts on it, so the driver records it instead of stopping.
verdict=""
jenkins-upgrade-check() {
  local rc=0
  case $1 in
  verify)
    command jenkins-upgrade-check "$@" || rc=$?
    case $rc in 0) verdict=GO ;; 1) verdict="ROLL BACK" ;; *) return $rc ;; esac
    log "verify: $verdict"
    ;;
  dry-run)
    command jenkins-upgrade-check "$@" || rc=$?
    [ $rc -le 1 ] || return $rc
    log "rehearsal: $([ $rc -eq 0 ] && echo clean || echo 'found problems')"
    [ $rc -eq 0 ] || [ "$scenario" != go ] || fail "the rehearsal of the plan was not clean"
    ;;
  *) command jenkins-upgrade-check "$@" ;;
  esac
}

windows=""
quiet_at=""
run_block() {
  local n=$1 sec=$2 lead=$3 f=blocks/$1.sh img
  # The parts a person does by hand: the token, and the image in compose.
  sed -i 's/^read -rs JENKINS_TOKEN.*$/: # JENKINS_TOKEN is set/' "$f"
  img=$(grep -o 'change the image to [^ ]*' "$f" | awk '{print $NF}' || true)
  [ -z "$img" ] || echo "JENKINS_IMAGE=$img" >.env
  if grep -q 'put the image back' "$f"; then
    echo "JENKINS_IMAGE=$from_image" >.env
  fi
  # The mistakes.
  if [ "$lead" = "Upgrade to $first_hop" ]; then
    case $scenario in
    forgot-plugins) sed -i '/step[0-9]*-plugins/d' "$f" ;;
    missing-plugin) echo 'rm -rf "$JENKINS_HOME/plugins/workflow-job.jpi" "$JENKINS_HOME/plugins/workflow-job"' >>"$f" ;;
    esac
  fi
  log "runbook: $sec: $lead"
  if grep -q '/quietDown"' "$f"; then quiet_at=$(date +%s); fi
  verdict=""
  # shellcheck disable=SC1090
  source "$f" </dev/null
  if grep -q '/cancelQuietDown"' "$f" && [ -n "$quiet_at" ]; then
    windows="$windows $(( ($(date +%s) - quiet_at + 59) / 60 )) min ($sec)"
    quiet_at=""
  fi
}

rolled_back=""
while IFS=$'\t' read -r -u 9 n sec lead; do
  if [ "$sec" = "Roll back" ]; then continue; fi
  run_block "$n" "$sec" "$lead"
  if [ "$verdict" = "ROLL BACK" ]; then
    rolled_back="$(core_now || echo 'did not start')"
    break
  fi
done 9<blocks/index.tsv

if [ -n "$rolled_back" ]; then
  log "rolling back from $rolled_back; $(disk_free) free on /"
  while IFS=$'\t' read -r -u 9 n sec lead; do
    if [ "$sec" = "Roll back" ]; then run_block "$n" "$sec" "$lead"; fi
  done 9<blocks/index.tsv
  if [ -n "$quiet_at" ]; then windows="$windows $(( ($(date +%s) - quiet_at + 59) / 60 )) min (rolled back)"; fi
  [ "$verdict" = GO ] || fail "verify after the rollback says $verdict, not GO"
fi

# --- Did users get back what they had? ------------------------------------
case $scenario in
go)
  [ -z "$rolled_back" ] || fail "verify said ROLL BACK at $rolled_back"
  want=$target
  ;;
*)
  [ -n "$rolled_back" ] || fail "verify never said ROLL BACK"
  want=$from
  ;;
esac
case $scenario in
missing-plugin)
  grep -l workflow-job /var/tmp/jenkins-upgrade/verify-"$first_hop".md >/dev/null 2>&1 ||
    fail "the ROLL BACK report does not name workflow-job"
  ;;
old-agent)
  grep -l 'agent1: runs Java 17' /var/tmp/jenkins-upgrade/verify-2*.md >/dev/null 2>&1 ||
    fail "the ROLL BACK report does not say agent1's Java is too old"
  ;;
esac
wait_up || fail "Jenkins is not up at the end"
now=$(core_now)
[ "$now" = "$want" ] || fail "Jenkins is $now at the end, not $want"
[ "$(api "$url/api/json?tree=quietingDown" | jq -r .quietingDown)" = false ] || fail "Jenkins is still in quiet mode at the end"
builds >builds-after.txt
lost=$(comm -23 builds-before.txt builds-after.txt)
[ -z "$lost" ] || fail "builds from before the upgrade are gone: $(echo "$lost" | paste -sd' ')"
if agent_online 36; then log "agent1 is online on $now"; else fail "agent1 is not online on $now"; fi
run_jobs | tee last-builds.txt
bad=$(grep -v ' SUCCESS$' last-builds.txt || true)
[ -z "$bad" ] || fail "builds on $now did not pass: $(echo "$bad" | paste -sd' ')"
log "Jenkins $now: $(list_jobs | wc -l) jobs, $(wc -l <builds-before.txt) old builds kept, $(grep -c ' SUCCESS$' last-builds.txt || true) of $(wc -l <last-builds.txt) new builds passed; $(disk_free) free on /"
