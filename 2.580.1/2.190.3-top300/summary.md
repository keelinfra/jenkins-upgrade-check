### 2.190.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.190.3 | jenkins/jenkins:2.190.3 | started | 304 | 2 | 2 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 304 | 2 | 2 | 0 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 361 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | did not start | 384 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan |  | did not start | 0 | 0 | 0 | 0 | 0 |

- baseline: **Parameterized-Remote-Trigger 3.1.6.3** other: Credentials Plugin version 2.3.11 is older than required. To fix, install version 2.6.1.1 or later.
- baseline: **docker-build-step 2.8** other: Credentials Plugin version 2.3.11 is older than required. To fix, install version 2.3.19 or later.
- as-is: **Parameterized-Remote-Trigger 3.1.6.3** old-dependency: needs credentials 2.6.1.1 or newer; 2.3.11 is installed
- as-is: **docker-build-step 2.8** old-dependency: needs credentials 2.3.19 or newer; 2.3.11 is installed

