### 2.89.4 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.89.4 | jenkins/jenkins:2.89.4 | started | 108 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 108 | 42 | 7 | 4 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk11 | started | 130 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 147 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 147 | 0 | 0 | 0 | 0 |

- as-is: **workflow-support 2.20** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 22 more)
- as-is: **junit 1.27** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 13 more)
- as-is: **command-launcher 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **jaxb 2.3.0** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **pipeline-utility-steps 2.1.0** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **script-security 1.71** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **trilead-api 1.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: loaded in part (extensions failed): git-server, ldap, mailer, ssh-slaves
- planned 2.346.3: Java 11 because gradle needs Java 11 or newer

