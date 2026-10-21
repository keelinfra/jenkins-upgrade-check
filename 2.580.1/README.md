# Upgrading every Jenkins LTS line to 2.580.1

Rehearsed 2026-10-21 with [jenkins-upgrade-check](https://github.com/keelinfra/jenkins-upgrade-check). Each row is a controller on the last release of an LTS line with the most installed plugins (Jenkins usage statistics, March 2026), each at the newest release that line can load. Every rehearsal starts Jenkins in a container with no network and records which plugins fail to load.

- **As it is**: the starting controller itself.
- **Core only**: the same plugin files on 2.580.1, as when only the war, the package or the image is replaced.
- **Planned**: the path `jenkins-upgrade-check plan` gives, which switches Java where no Java runs both versions, with the plugin versions it computes for each stop. Every stop is rehearsed; the cell shows the first that does not load cleanly.
- **Controllers**: share of the 234780 controllers that reported usage in March 2026 running a release of that line.

## Summary

**300 most installed plugins, 38 starting versions.** Replacing only the core: Jenkins does not start from 12 of them, and starts with plugins that fail to load from 10 more; 15.0% of all reporting controllers are on those lines. Following the plan: every stop loads cleanly from 18 of 38. 15 starting points already have plugins that do not load before any upgrade; see their rows.

**100 most installed plugins, 38 starting versions.** Replacing only the core: Jenkins does not start from 0 of them, and starts with plugins that fail to load from 9 more; 2.8% of all reporting controllers are on those lines. Following the plan: every stop loads cleanly from 31 of 38. 3 starting points already have plugins that do not load before any upgrade; see their rows.

## 300 most installed plugins

| From | Controllers | As it is | Core only → 2.580.1 | Planned path | Planned |
|---|---|---|---|---|---|
| [2.60.3](2.60.3-top300/summary.md) | <0.1% | ⚠️ 6 of 305 fail: Parameterized-Remote-Trigger, artifactory, cucumber-reports, and 3 more | ⚠️ 6 of 305 fail: docker-java-api, Parameterized-Remote-Trigger, artifactory, and 2 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 353 fail: gradle (at 2.346.3) |
| [2.73.3](2.73.3-top300/summary.md) | <0.1% | ⚠️ 6 of 305 fail: Parameterized-Remote-Trigger, artifactory, cucumber-reports, and 3 more | ⚠️ 6 of 305 fail: docker-java-api, Parameterized-Remote-Trigger, artifactory, and 2 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 351 fail: gradle (at 2.346.3) |
| [2.89.4](2.89.4-top300/summary.md) | <0.1% | ⚠️ 6 of 305 fail: Parameterized-Remote-Trigger, artifactory, cucumber-reports, and 3 more | ⚠️ 125 of 305 fail: junit, workflow-support, aws-credentials, and 20 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 348 fail: gradle (at 2.346.3) |
| [2.107.3](2.107.3-top300/summary.md) | <0.1% | ⚠️ 6 of 305 fail: Parameterized-Remote-Trigger, artifactory, cucumber-reports, and 3 more | ⚠️ 125 of 305 fail: junit, workflow-support, aws-credentials, and 20 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 348 fail: gradle (at 2.346.3) |
| [2.121.3](2.121.3-top300/summary.md) | 0.1% | ⚠️ 5 of 304 fail: Parameterized-Remote-Trigger, artifactory, docker-build-step, and 2 more | ⚠️ 125 of 304 fail: junit, workflow-support, aws-credentials, and 22 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 347 fail: gradle (at 2.346.3) |
| [2.138.4](2.138.4-top300/summary.md) | 0.2% | ⚠️ 2 of 304 fail: Parameterized-Remote-Trigger, docker-build-step | ⚠️ 127 of 304 fail: workflow-support, junit, git-client, and 21 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 340 fail: gradle (at 2.346.3) |
| [2.150.3](2.150.3-top300/summary.md) | 0.2% | ⚠️ 4 of 304 fail: Parameterized-Remote-Trigger, blueocean, blueocean-jira, and 1 more | ⚠️ 128 of 304 fail: workflow-support, junit, git-client, and 21 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 2 of 340 fail: gradle (at 2.346.3) |
| [2.164.3](2.164.3-top300/summary.md) | 0.3% | ⚠️ 4 of 304 fail: Parameterized-Remote-Trigger, blueocean, blueocean-jira, and 1 more | ❌ Jenkins does not start: blueocean-jira, Parameterized-Remote-Trigger, docker-build-step, and 1 more | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start: Office-365-Connector, Parameterized-Remote-Trigger, active-directory, and 763 more (at 2.541.3) |
| [2.176.4](2.176.4-top300/summary.md) | 0.4% | ⚠️ 2 of 304 fail: Parameterized-Remote-Trigger, docker-build-step | ❌ Jenkins does not start: Parameterized-Remote-Trigger, docker-build-step, ec2 | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start: plugin-util-api, workflow-cps, workflow-job, and 12 more (at 2.541.3) |
| [2.190.3](2.190.3-top300/summary.md) | 0.4% | ⚠️ 2 of 304 fail: Parameterized-Remote-Trigger, docker-build-step | ❌ Jenkins does not start: Parameterized-Remote-Trigger, docker-build-step | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.541.3) |
| [2.204.6](2.204.6-top300/summary.md) | 0.4% | ⚠️ 2 of 304 fail: Parameterized-Remote-Trigger, docker-build-step | ❌ Jenkins does not start: Parameterized-Remote-Trigger, docker-build-step | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start: aws-java-sdk-minimal, Office-365-Connector, Parameterized-Remote-Trigger, and 693 more (at 2.541.3) |
| [2.222.4](2.222.4-top300/summary.md) | 0.5% | ✅ all 304 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.235.5](2.235.5-top300/summary.md) | 0.5% | ✅ all 304 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.249.3](2.249.3-top300/summary.md) | 0.5% | ✅ all 303 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.263.4](2.263.4-top300/summary.md) | 0.7% | ✅ all 303 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.277.4](2.277.4-top300/summary.md) | 0.5% | ✅ all 303 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.541.3) |
| [2.289.3](2.289.3-top300/summary.md) | 0.7% | ⚠️ 97 of 302 fail: Parameterized-Remote-Trigger, amazon-ecr, ansible, and 94 more | ⚠️ 97 of 302 fail: plain-credentials, jsch, authentication-tokens, and 26 more | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.303.3](2.303.3-top300/summary.md) | 0.7% | ⚠️ 98 of 302 fail: Parameterized-Remote-Trigger, amazon-ecr, ansible, and 95 more | ⚠️ 98 of 302 fail: plain-credentials, jsch, authentication-tokens, and 26 more | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start (at 2.580.1) |
| [2.319.3](2.319.3-top300/summary.md) | 0.9% | ✅ all 302 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start: Office-365-Connector, Parameterized-Remote-Trigger, active-directory, and 680 more (at 2.580.1) |
| [2.332.4](2.332.4-top300/summary.md) | 1.4% | ✅ all 302 load | ❌ Jenkins does not start | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ❌ Jenkins does not start: Office-365-Connector, Parameterized-Remote-Trigger, active-directory, and 692 more (at 2.580.1) |
| [2.346.3](2.346.3-top300/summary.md) | 2.2% | ✅ all 302 load | ❌ Jenkins does not start | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 349 load |
| [2.361.4](2.361.4-top300/summary.md) | 0.7% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 342 load |
| [2.375.4](2.375.4-top300/summary.md) | 0.8% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 342 load |
| [2.387.3](2.387.3-top300/summary.md) | 1.3% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 341 load |
| [2.401.3](2.401.3-top300/summary.md) | 1.0% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 328 load |
| [2.414.3](2.414.3-top300/summary.md) | 1.2% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 328 load |
| [2.426.3](2.426.3-top300/summary.md) | 2.6% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 327 load |
| [2.440.3](2.440.3-top300/summary.md) | 2.1% | ✅ all 300 load | ✅ all 300 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 324 load |
| [2.452.4](2.452.4-top300/summary.md) | 2.2% | ✅ all 301 load | ✅ all 301 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 322 load |
| [2.462.3](2.462.3-top300/summary.md) | 3.7% | ✅ all 301 load | ✅ all 301 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 322 load |
| [2.479.3](2.479.3-top300/summary.md) | 3.5% | ✅ all 303 load | ✅ all 303 load | 2.580.1 (Java 21) | ✅ all 319 load |
| [2.492.3](2.492.3-top300/summary.md) | 4.2% | ⚠️ 1 of 303 fail: gitlab-plugin | ⚠️ 1 of 303 fail: gitlab-plugin | 2.580.1 (Java 21) | ✅ all 319 load |
| [2.504.3](2.504.3-top300/summary.md) | 4.7% | ✅ all 304 load | ✅ all 304 load | 2.580.1 (Java 21) | ✅ all 317 load |
| [2.516.3](2.516.3-top300/summary.md) | 5.7% | ✅ all 304 load | ✅ all 304 load | 2.580.1 (Java 21) | ✅ all 317 load |
| [2.528.3](2.528.3-top300/summary.md) | 11.6% | ✅ all 304 load | ✅ all 304 load | 2.580.1 (Java 21) | ✅ all 317 load |
| [2.541.3](2.541.3-top300/summary.md) | 23.2% | ✅ all 317 load | ✅ all 317 load | 2.580.1 (Java 21) | ✅ all 317 load |
| [2.555.3](2.555.3-top300/summary.md) | – | ✅ all 317 load | ✅ all 317 load | 2.580.1 (Java 21) | ✅ all 317 load |
| [2.568.3](2.568.3-top300/summary.md) | – | ⚠️ 75 of 308 fail: junit, github, aws-java-sdk2-apigateway, and 35 more (9 not downloadable) | ✅ all 317 load | 2.580.1 (Java 21) | ✅ all 317 load |

## 100 most installed plugins

| From | Controllers | As it is | Core only → 2.580.1 | Planned path | Planned |
|---|---|---|---|---|---|
| [2.60.3](2.60.3-top100/summary.md) | <0.1% | ✅ all 108 load | ✅ all 108 load | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 133 fail: gradle (at 2.346.3) |
| [2.73.3](2.73.3-top100/summary.md) | <0.1% | ✅ all 108 load | ✅ all 108 load | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 131 fail: gradle (at 2.346.3) |
| [2.89.4](2.89.4-top100/summary.md) | <0.1% | ✅ all 108 load | ⚠️ 42 of 108 fail: workflow-support, junit, command-launcher, and 4 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 130 fail: gradle (at 2.346.3) |
| [2.107.3](2.107.3-top100/summary.md) | <0.1% | ✅ all 108 load | ⚠️ 42 of 108 fail: workflow-support, junit, command-launcher, and 4 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 130 fail: gradle (at 2.346.3) |
| [2.121.3](2.121.3-top100/summary.md) | 0.1% | ✅ all 108 load | ⚠️ 44 of 108 fail: workflow-support, junit, config-file-provider, and 6 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 130 fail: gradle (at 2.346.3) |
| [2.138.4](2.138.4-top100/summary.md) | 0.2% | ✅ all 110 load | ⚠️ 44 of 110 fail: workflow-support, junit, git-client, and 7 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 127 fail: gradle (at 2.346.3) |
| [2.150.3](2.150.3-top100/summary.md) | 0.2% | ✅ all 110 load | ⚠️ 45 of 110 fail: workflow-support, junit, git-client, and 7 more | 2.346.3 (Java 8) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ⚠️ 1 of 127 fail: gradle (at 2.346.3) |
| [2.164.3](2.164.3-top100/summary.md) | 0.3% | ✅ all 110 load | ⚠️ 2 of 110 fail: jaxb, trilead-api | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 141 load |
| [2.176.4](2.176.4-top100/summary.md) | 0.4% | ✅ all 107 load | ⚠️ 2 of 107 fail: jaxb, trilead-api | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 138 load |
| [2.190.3](2.190.3-top100/summary.md) | 0.4% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 138 load |
| [2.204.6](2.204.6-top100/summary.md) | 0.4% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 137 load |
| [2.222.4](2.222.4-top100/summary.md) | 0.5% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 135 load |
| [2.235.5](2.235.5-top100/summary.md) | 0.5% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 135 load |
| [2.249.3](2.249.3-top100/summary.md) | 0.5% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 134 load |
| [2.263.4](2.263.4-top100/summary.md) | 0.7% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 130 load |
| [2.277.4](2.277.4-top100/summary.md) | 0.5% | ✅ all 107 load | ✅ all 107 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 129 load |
| [2.289.3](2.289.3-top100/summary.md) | 0.7% | ⚠️ 30 of 106 fail: authentication-tokens, config-file-provider, credentials, and 27 more | ⚠️ 30 of 106 fail: jsch, plain-credentials, pipeline-input-step, and 8 more | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 121 load |
| [2.303.3](2.303.3-top100/summary.md) | 0.7% | ⚠️ 31 of 104 fail: authentication-tokens, config-file-provider, credentials, and 28 more | ⚠️ 31 of 104 fail: plain-credentials, jsch, pipeline-input-step, and 8 more | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 119 load |
| [2.319.3](2.319.3-top100/summary.md) | 0.9% | ✅ all 104 load | ✅ all 104 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 119 load |
| [2.332.4](2.332.4-top100/summary.md) | 1.4% | ✅ all 103 load | ✅ all 103 load | 2.462.3 (Java 11) → 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 118 load |
| [2.346.3](2.346.3-top100/summary.md) | 2.2% | ⚠️ 2 of 103 fail: metrics, pipeline-graph-view | ✅ all 103 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 116 load |
| [2.361.4](2.361.4-top100/summary.md) | 0.7% | ✅ all 101 load | ✅ all 101 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 113 load |
| [2.375.4](2.375.4-top100/summary.md) | 0.8% | ✅ all 101 load | ✅ all 101 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 113 load |
| [2.387.3](2.387.3-top100/summary.md) | 1.3% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 112 load |
| [2.401.3](2.401.3-top100/summary.md) | 1.0% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 107 load |
| [2.414.3](2.414.3-top100/summary.md) | 1.2% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 107 load |
| [2.426.3](2.426.3-top100/summary.md) | 2.6% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 106 load |
| [2.440.3](2.440.3-top100/summary.md) | 2.1% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 106 load |
| [2.452.4](2.452.4-top100/summary.md) | 2.2% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 106 load |
| [2.462.3](2.462.3-top100/summary.md) | 3.7% | ✅ all 100 load | ✅ all 100 load | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 106 load |
| [2.479.3](2.479.3-top100/summary.md) | 3.5% | ✅ all 102 load | ✅ all 102 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.492.3](2.492.3-top100/summary.md) | 4.2% | ✅ all 103 load | ✅ all 103 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.504.3](2.504.3-top100/summary.md) | 4.7% | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.516.3](2.516.3-top100/summary.md) | 5.7% | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.528.3](2.528.3-top100/summary.md) | 11.6% | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.541.3](2.541.3-top100/summary.md) | 23.2% | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.555.3](2.555.3-top100/summary.md) | – | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |
| [2.568.3](2.568.3-top100/summary.md) | – | ✅ all 105 load | ✅ all 105 load | 2.580.1 (Java 21) | ✅ all 105 load |

## Limits

- Plugins are at the newest release each line can load. A controller whose plugins are older than that has more to change, not less.
- There is no real configuration and no network: what plugins do with jobs, credentials and outside systems is not exercised. `jenkins-upgrade-check dry-run --home` rehearses your own configuration.
- A plugin that loads can still fail at run time, in a pipeline step or a page. The rehearsal shows what fails at startup.
- The usage statistics count controllers that send them; controllers inside closed networks usually do not.
