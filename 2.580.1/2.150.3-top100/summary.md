### 2.150.3 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.150.3 | jenkins/jenkins:2.150.3 | started | 110 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 110 | 45 | 10 | 3 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk11 | started | 128 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 145 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 145 | 0 | 0 | 0 | 0 |

- as-is: **workflow-support 3.4** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 21 more)
- as-is: **junit 1.28** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 8 more)
- as-is: **git-client 3.2.1** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 5 more)
- as-is: **config-file-provider 3.7.2** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 1 more)
- as-is: **command-launcher 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **jaxb 2.3.0** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **jdk-tool 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **pipeline-utility-steps 2.7.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **script-security 1.71** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **trilead-api 1.0.3** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jaxb -> Plugin:command-launcher
- as-is: loaded in part (extensions failed): ldap, mailer, ssh-slaves
- planned 2.346.3: Java 11 because gradle needs Java 11 or newer

