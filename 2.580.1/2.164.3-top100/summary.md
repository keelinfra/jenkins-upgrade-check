### 2.164.3 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.164.3 | jenkins/jenkins:2.164.3 | started | 110 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 110 | 2 | 2 | 6 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 132 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 141 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 141 | 0 | 0 | 0 | 0 |

- as-is: **jaxb 2.3.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:trilead-api -> Plugin:jaxb -> Plugin:trilead-api
- as-is: **trilead-api 1.0.3** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:trilead-api -> Plugin:jaxb -> Plugin:trilead-api
- as-is: loaded in part (extensions failed): email-ext, git-server, github, ldap, mailer, ssh-slaves

