### 2.73.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.73.3 | jenkins/jenkins:2.73.3 | started | 305 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 305 | 0 | 0 | 9 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk11 | started | 353 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 404 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 405 | 0 | 0 | 0 | 0 |

- as-is: loaded in part (extensions failed): copyartifact, email-ext, git-server, github, ldap, mailer, ssh-slaves, swarm, translation
- planned 2.346.3: Java 11 because gradle needs Java 11 or newer

