### 2.107.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.107.3 | jenkins/jenkins:2.107.3 | started | 305 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 305 | 121 | 19 | 10 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk11 | started | 348 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 397 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 399 | 0 | 0 | 1 | 0 |

- as-is: **junit 1.27** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 54 more)
- as-is: **workflow-support 2.20** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 44 more)
- as-is: **aws-credentials 1.28.2** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 3 more)
- as-is: **badge 1.8** other: Required plugin is disabled: Script Security Plugin (script-security) (takes down 1 more)
- as-is: **command-launcher 1.4** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb
- as-is: **embeddable-build-status 2.0.4** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **extended-choice-parameter 346.vd87693c5a_86c** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **extensible-choice-parameter 1.8.0** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **folder-properties 1.2.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **groovy 2.4** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **http_request 1.8.26** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **jaxb 2.3.0** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb
- as-is: **job-dsl 1.71** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **permissive-script-security 0.7** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **pipeline-utility-steps 2.1.0** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **script-security 1.71** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb
- as-is: **scriptler 3.1** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: **trilead-api 1.0.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:jaxb -> Plugin:command-launcher -> Plugin:script-security -> Plugin:trilead-api -> Plugin:jaxb
- as-is: **uno-choice 2.2.2** other: Required plugin is disabled: Script Security Plugin (script-security)
- as-is: loaded in part (extensions failed): configurationslicing, description-setter, flexible-publish, git-server, ldap, mailer, postbuildscript, ssh-slaves, swarm, translation
- planned 2.346.3: Java 11 because gradle needs Java 11 or newer
- planned 2.580.1: loaded in part (extensions failed): translation

