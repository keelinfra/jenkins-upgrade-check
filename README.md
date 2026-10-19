# jenkins-upgrade-check

Plan a Jenkins upgrade, then rehearse it before you touch the real controller.

Upgrading an old controller is rarely one step. Jenkins 2.346.3 cannot go straight
to today's LTS: no Java version runs both, so the controller and every agent have
to change Java on the way, and the plugin versions you can install depend on the
core you are on at each stop. Plugins that break on an API change (the move to
Jakarta EE 9 in 2.475, for example) do not show it in any metadata: they fail when
Jenkins starts.

`jenkins-upgrade-check` does two things:

1. **plan** works out the path (the Java version and the core for each stop), the
   plugin versions for each step including new dependencies, and what the update
   center already knows is wrong: security warnings today and after the upgrade,
   deprecated plugins, plugins no longer distributed, and agents that cannot follow
   (for example Solaris, HP-UX or old AIX agents that have no new enough Java).
2. **dry-run** starts the target Jenkins in a throwaway Docker container with your
   plugin set and reports which plugins fail to load and why.

## Quick start

Get your plugin list. In **Manage Jenkins → Script Console** run:

```groovy
Jenkins.instance.pluginManager.plugins.sort{it.shortName}.each{println "${it.shortName}:${it.version}"}
```

Save the output as `plugins.txt`, then:

```sh
go install github.com/keelinfra/jenkins-upgrade-check/cmd/jenkins-upgrade-check@latest

jenkins-upgrade-check plan --plugins-file plugins.txt --core 2.426.3 --to lts \
  -o plan.md --json plan.json
jenkins-upgrade-check dry-run --plan plan.json -o rehearsal.md   # needs Docker
```

Or let it read the controller (GET requests only, see below):

```sh
export JENKINS_TOKEN=...   # an API token, ideally of a read-only user
jenkins-upgrade-check collect --url https://jenkins.example.com --user me -o snapshot.json
jenkins-upgrade-check plan --snapshot snapshot.json --to lts -o plan.md --json plan.json
```

`collect` also records each agent's OS, and its Java version when the Versions Node
Monitors plugin (`versioncolumn`) is installed, so the plan can say which agents
will not reconnect after the upgrade.

## What it reads, what it sends

- `collect` sends three GET requests to the controller (core version, plugins,
  agents) and refuses any other method. It writes nothing to the controller and
  needs no plugin on it. The token is read from `JENKINS_TOKEN` so it stays out of
  shell history.
- `plan` downloads public data from the Jenkins update center
  (`updates.jenkins.io`): plugin releases, security warnings, deprecations.
- `dry-run` downloads the plugin files from the update center, checks their
  SHA-256, and starts `jenkins/jenkins:<version>-jdk<N>` with only those plugins.
  No job, credential or configuration from your controller goes into the container.

## Limits

- A plugin that loads can still break at run time: a pipeline step, a UI page, an
  integration. The rehearsal tells you what fails at startup, which is where most
  upgrade breakage shows, not everything.
- Private plugins (not on the update center) are listed but not rehearsed.
- Configuration as Code, Groovy init scripts and job configuration are not checked yet.

## Commands

```
jenkins-upgrade-check collect  --url URL --user USER [-o snapshot.json]
jenkins-upgrade-check plan     (--snapshot snapshot.json | --plugins-file FILE --core VERSION) [--to lts]
jenkins-upgrade-check dry-run  (--plan plan.json | --plugins-file FILE --core VERSION [--resolve])
jenkins-upgrade-check resolve  --core VERSION (--plugins-file FILE | --names a,b | --top-plugins CSV --top N)
```

`hack/experiment.sh` rehearses a typical controller (the most installed plugins
from the Jenkins usage statistics) going from an older version to the current LTS,
once with only the core swapped and once with the plan. The "Upgrade experiment"
workflow runs it for several starting versions.

## Help with an upgrade

If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks
and upgrades: https://keelinfra.io/jenkins/pricing/

## License

Apache-2.0
