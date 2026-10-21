# jenkins-upgrade-check

[![CI](https://github.com/keelinfra/jenkins-upgrade-check/actions/workflows/ci.yml/badge.svg)](https://github.com/keelinfra/jenkins-upgrade-check/actions/workflows/ci.yml)
[![Upgrade matrix](https://github.com/keelinfra/jenkins-upgrade-check/actions/workflows/experiment.yml/badge.svg)](https://github.com/keelinfra/jenkins-upgrade-check/tree/matrix)
[![Release](https://img.shields.io/github/v/release/keelinfra/jenkins-upgrade-check)](https://github.com/keelinfra/jenkins-upgrade-check/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/keelinfra/jenkins-upgrade-check)](go.mod)
[![License](https://img.shields.io/github/license/keelinfra/jenkins-upgrade-check)](LICENSE)

Plan a Jenkins upgrade, rehearse it on a copy of your controller, and know within minutes
of the real upgrade whether to keep it or roll back, and why.

It is for the people who run a Jenkins controller and have put off upgrading it. One
static binary; the rehearsal needs Docker. Nothing is installed on the controller.

**Why rehearse.** We upgraded a controller with the 300 most installed plugins from the
last release of every LTS line since 2.60.3 to 2.580.1
([upgrade matrix](https://keelinfra.io/jenkins/upgrade-matrix/), 2026-10-21). Replacing only
the war, package or image, Jenkins did not start from 12 of the 38 lines and started with
plugins failing from 12 more; 13.6% of the controllers that report usage are on those lines.
Following the path `plan` gives, every stop loaded every plugin, from all 38. From 2.356 or
older the core-only upgrade fails with network access too: 2.580.1 no longer ships the plugins
that were split out of core.

Upgrading an old controller is rarely one step. Jenkins 2.346.3 cannot go straight
to today's LTS: no Java version runs both, so the controller and every agent have
to change Java on the way, and the plugin versions you can install depend on the
core you are on at each stop. Plugins that break on an API change (the move to
Jakarta EE 9 in 2.475, for example) do not show it in any metadata: they fail when
Jenkins starts.

`jenkins-upgrade-check` covers the whole upgrade:

1. **plan** works out the path (the Java version and the core for each stop), the
   plugin versions for each step including new dependencies, and what the update
   center already knows is wrong: security warnings today and after the upgrade,
   deprecated plugins, plugins no longer distributed, and agents that cannot follow
   (for example Solaris, HP-UX or old AIX agents that have no new enough Java).
   With `--runbook` it also writes the commands for the upgrade and the rollback,
   for a Debian or RPM package, a war or the Docker image.
2. **dry-run** starts the target Jenkins in a throwaway Docker container, with no
   network and no executors, with your plugin set and, with `--home`, a copy of your
   configuration: jobs, folders, nodes, credentials and settings, without build
   records. It reports what failed, grouped by cause, why, and what to do: a plugin
   that needs a newer core or Java, a missing or too old dependency, a plugin built
   for the pre-Jakarta servlet API, the jobs that would disappear and the plugin
   they need, the settings Jenkins could no longer read.
3. **verify** checks the real controller right after the upgrade against a
   snapshot from before (plugins, jobs and folders, agents, settings) and says GO or
   ROLL BACK, and why.

## How the upgrade stays safe for users

The runbook follows the same method for every install type:

- Rehearse on a copy until it is clean. Nothing reaches users or other systems.
- In the window, Jenkins stops taking new builds and waits for running ones to
  finish (quiet down), so nothing is killed. Then it stops and is backed up.
- The new version starts in quiet mode: an init script from `jenkins-upgrade-check
  probe --quiet` keeps builds from starting and records what loaded and what did not.
- `verify` says GO or ROLL BACK. On GO you cancel quiet mode and users are back.
  On ROLL BACK you restore the backup and the old version; since no build ran on the
  new version, nothing is lost.

Downtime is the time to stop, back up, start and verify. Open-source Jenkins cannot
be upgraded with no downtime at all: two controllers cannot share one JENKINS_HOME.

## Install

On Linux or macOS (amd64 or arm64):

```sh
curl -fsSL https://github.com/keelinfra/jenkins-upgrade-check/releases/latest/download/install.sh | sh
```

It downloads the binary for your system from the latest release, checks it against
the release's `sha256sums.txt`, and installs it in `/usr/local/bin`, or `~/.local/bin`
when that is not writable. `sh -s -- --dir DIR --version 0.1.0` chooses the place and
the release. On Windows, download `jenkins-upgrade-check_windows_amd64.exe` from the
[releases](https://github.com/keelinfra/jenkins-upgrade-check/releases). It is one
static binary with no dependencies; `dry-run` also needs Docker. With Go 1.24 or newer,
`go install github.com/keelinfra/jenkins-upgrade-check/cmd/jenkins-upgrade-check@latest`
works too.

## Quick start

Get your plugin list. In **Manage Jenkins → Script Console** run:

```groovy
Jenkins.instance.pluginManager.plugins.sort{it.shortName}.each{println "${it.shortName}:${it.version}"}
```

Save the output as `plugins.txt`, then:

```sh
jenkins-upgrade-check plan --plugins-file plugins.txt --core 2.426.3 --to lts \
  -o plan.md --json plan.json --runbook runbook.md --install deb
jenkins-upgrade-check dry-run --plan plan.json --home /path/to/jenkins_home -o rehearsal.md   # needs Docker
```

Or let it read the controller (GET requests only, see below):

```sh
export JENKINS_TOKEN=...   # an API token, ideally of a read-only user
jenkins-upgrade-check collect --url https://jenkins.example.com --user me -o snapshot.json
jenkins-upgrade-check plan --snapshot snapshot.json --to lts -o plan.md --json plan.json
```

`collect` also records each agent's OS and Java version, so the plan can say which
agents will not reconnect after the upgrade. The Java version comes from the Versions
Node Monitors plugin (`versioncolumn`) when it is installed, and otherwise from each
online agent's System Information page, which needs Overall/Administer.

## What it reads, what it sends

- `collect` and `verify` send GET requests to the controller (core version and job
  names, plugins, agents; without `versioncolumn`, one more per online agent for its
  Java; after an upgrade, one per agent that did not come back, for its connection
  log) and refuse any other method. They write nothing to the
  controller and need no plugin on it. The token is read from `JENKINS_TOKEN` so it
  stays out of shell history.
- `plan` downloads public data from the Jenkins update center
  (`updates.jenkins.io`): plugin releases, security warnings, deprecations.
- `dry-run` and `download` fetch the plugin files from the update center and check
  their SHA-256 (plugins the update center no longer lists come from its archive,
  which has no checksum). `dry-run` starts `jenkins/jenkins:<version>-jdk<N>` with
  `--network none`. With `--home`, the configuration files of JENKINS_HOME, including
  `secret.key` and `secrets/` so credentials still decrypt, are copied into the
  container (no image is built); the container and its work directory are removed
  afterwards unless you pass `--keep`.
- The probe script writes one file, `jenkins-upgrade-check-probe.json`, into
  JENKINS_HOME. With `--quiet` it also puts Jenkins into quiet mode; nothing else.

## Limits

- A plugin that loads can still break at run time: a pipeline step, a UI page, an
  integration. The rehearsal tells you what fails at startup, which is where most
  upgrade breakage shows, not everything.
- Private plugins (not on the update center) are listed but not rehearsed.
- The rehearsal has no network, so whatever a plugin does with an outside system
  (an SCM, a cloud, an LDAP server) is not exercised.
- `verify` compares job names four folder levels deep.

## Commands

```
jenkins-upgrade-check collect  --url URL --user USER [-o snapshot.json]
jenkins-upgrade-check plan     (--snapshot snapshot.json | --plugins-file FILE --core VERSION) [--to lts]
                               [--runbook runbook.md --install deb|rpm|war|docker]
jenkins-upgrade-check dry-run  (--plan plan.json [--step N] | --plugins-file FILE --core VERSION [--resolve] [--from VERSION]) [--home DIR]
jenkins-upgrade-check download --plan plan.json [--step N] --dir DIR
jenkins-upgrade-check probe    [--quiet]
jenkins-upgrade-check verify   --before before.json --url URL --user USER [--plan plan.json] [--core VERSION]
                               [--home JENKINS_HOME] [--log FILE]
jenkins-upgrade-check resolve  --core VERSION (--plugins-file FILE | --names a,b | --top-plugins CSV --top N)
```

`hack/experiment.sh` rehearses a typical controller (the most installed plugins
from the Jenkins usage statistics) going from an older version to the current LTS:
once with only the core swapped, and at every stop of the plan. The "Upgrade matrix"
workflow runs it from the last release of every LTS line since 2.60.3 whenever a new
LTS is out, and publishes the table to the
[matrix branch](https://github.com/keelinfra/jenkins-upgrade-check/tree/matrix).

## Status and help

jenkins-upgrade-check is new (0.x): commands and report formats may still change. Bugs and
questions go to [GitHub issues](https://github.com/keelinfra/jenkins-upgrade-check/issues);
security problems to [SECURITY.md](SECURITY.md). It gets free security updates until at least
31 December 2031.

To have your own plugin list rehearsed, or help with the upgrade itself, see
[keelinfra.io/jenkins/upgrade-check](https://keelinfra.io/jenkins/upgrade-check/) or write to
hello@keelinfra.io.

## License

Apache-2.0
