# Security policy

## Reporting a vulnerability

Please do not open a public issue for a security problem in jenkins-upgrade-check. Email **hello@keelinfra.io** with the subject **"jenkins-upgrade-check security report"**, and include:

- what is affected (a command, the runbook it writes, `install.sh` or the release workflow) and the version, from `jenkins-upgrade-check version`;
- how to reproduce it, and what an attacker gains.

We answer within two business days, agree on a fix and a disclosure date with you, and credit you in the release notes unless you prefer not to be named.

## Scope

A vulnerability in Jenkins or a plugin is fixed by the Jenkins project: report it to them (https://www.jenkins.io/security/). Problems in this tool belong here, for example how it handles the API token or a copy of JENKINS_HOME, how it downloads and checks plugins, the commands the runbook runs, `install.sh` and its checksum check, or the release process.

## Supported versions

jenkins-upgrade-check is supported until at least **31 December 2031**. Security updates are free, for everyone.

Fixes ship in a new release; rerun `install.sh` to get it. When a fix ships we publish a GitHub security advisory saying who is affected and what to do.

## Checking a release

Every release carries `sha256sums.txt`, which `install.sh` checks before it installs anything, and a CycloneDX SBOM (`jenkins-upgrade-check.cdx.json`) naming the Go version the binaries were built with. Releases from a public repository also carry signed build-provenance and SBOM attestations. To check a binary came unmodified from this repository's release workflow:

```sh
gh attestation verify jenkins-upgrade-check_linux_amd64 --repo keelinfra/jenkins-upgrade-check
```

## EU Cyber Resilience Act

keelinfra sells services for Jenkins upgrades, so we treat ourselves as this tool's manufacturer under the CRA ([Regulation (EU) 2024/2847](https://eur-lex.europa.eu/eli/reg/2024/2847/oj)). jenkins-upgrade-check is a default-class product (not listed in Annex III). Actively exploited vulnerabilities are reported through ENISA's Single Reporting Platform within 24 hours, as the CRA requires since 11 September 2026. How we read the CRA is on [keelinfra.io/security](https://keelinfra.io/security/#cyber-resilience-act).
