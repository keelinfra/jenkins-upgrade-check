### 2.176.4 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.176.4 | jenkins/jenkins:2.176.4 | started | 107 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 107 | 2 | 2 | 0 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 129 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 138 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 138 | 0 | 0 | 0 | 0 |

- as-is: **jaxb 2.3.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:trilead-api -> Plugin:jaxb -> Plugin:trilead-api
- as-is: **trilead-api 1.0.3** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:trilead-api -> Plugin:jaxb -> Plugin:trilead-api

