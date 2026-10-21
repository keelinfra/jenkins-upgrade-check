### 2.492.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.492.3 | jenkins/jenkins:2.492.3-jdk17 | started | 303 | 1 | 1 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 303 | 1 | 1 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 319 | 0 | 0 | 0 | 0 |

- baseline: **gitlab-plugin 1.9.10** missing-dependency: needs the jersey3-api plugin (3.1.11-2.v4a_b_db_cf4557b_ or newer), which is not installed
- as-is: **gitlab-plugin 1.9.10** missing-dependency: needs the jersey3-api plugin (3.1.11-2.v4a_b_db_cf4557b_ or newer), which is not installed

