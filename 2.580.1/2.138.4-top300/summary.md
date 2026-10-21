### 2.138.4 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.138.4 | jenkins/jenkins:2.138.4 | started | 304 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 304 | 126 | 22 | 7 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk11 | started | 340 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 389 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 391 | 0 | 0 | 1 | 0 |

- as-is: **workflow-support 3.4** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 47 more)
- as-is: **junit 1.28** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 35 more)
- as-is: **git-client 3.2.1** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 15 more)
- as-is: **aws-credentials 1.28.2** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 3 more)
- as-is: **config-file-provider 3.7.2** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 3 more)
- as-is: **badge 1.8** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 1 more)
- as-is: **command-launcher 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool
- as-is: **embeddable-build-status 2.0.4** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **extended-choice-parameter 346.vd87693c5a_86c** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **extensible-choice-parameter 1.8.0** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **folder-properties 1.2.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **gerrit-trigger 2.34.0** not-active: installed but not active after startup
- as-is: **groovy 2.4** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **http_request 1.8.26** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **jaxb 2.3.0** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool
- as-is: **jdk-tool 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool
- as-is: **job-dsl 1.76** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **permissive-script-security 0.7** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **pipeline-utility-steps 2.7.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **script-security 1.71** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool -> Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:jdk-tool
- as-is: **trilead-api 1.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jdk-tool -> Plugin:trilead-api -> Plugin:jdk-tool
- as-is: **uno-choice 2.2.2** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: loaded in part (extensions failed): configurationslicing, description-setter, flexible-publish, ldap, mailer, ssh-slaves, translation
- planned 2.346.3: Java 11 because gradle needs Java 11 or newer
- planned 2.580.1: loaded in part (extensions failed): translation

