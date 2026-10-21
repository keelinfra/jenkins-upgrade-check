### 2.121.3 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.121.3 | jenkins/jenkins:2.121.3 | started | 108 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 108 | 44 | 9 | 0 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk8 | started | 130 | 1 | 1 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 147 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 147 | 0 | 0 | 0 | 0 |

- as-is: **workflow-support 3.3** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 22 more)
- as-is: **junit 1.27** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 12 more)
- as-is: **config-file-provider 3.7.2** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 1 more)
- as-is: **command-launcher 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **jaxb 2.3.0** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **jdk-tool 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool
- as-is: **pipeline-utility-steps 2.7.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **script-security 1.71** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher
- as-is: **trilead-api 1.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool
- planned 2.346.3: **gradle 2.12.0.1** java-too-old: this release needs Java 11 or newer

