### 2.176.4 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.176.4 | jenkins/jenkins:2.176.4 | started | 304 | 2 | 2 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 304 | 3 | 3 | 0 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 361 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | did not start | 384 | 107 | 15 | 0 | 0 |
| planned: 2.580.1, plugins from plan |  | did not start | 0 | 0 | 0 | 0 | 0 |

- baseline: **Parameterized-Remote-Trigger 3.1.6.3** other: Credentials Plugin version 2.3.11 is older than required. To fix, install version 2.6.1.1 or later.
- baseline: **docker-build-step 2.8** other: Credentials Plugin version 2.3.11 is older than required. To fix, install version 2.3.19 or later.
- as-is: **Parameterized-Remote-Trigger 3.1.6.3** old-dependency: needs credentials 2.6.1.1 or newer; 2.3.11 is installed
- as-is: **docker-build-step 2.8** old-dependency: needs credentials 2.3.19 or newer; 2.3.11 is installed
- as-is: **ec2 1.50.3** other: Required plugin is disabled: Trilead API Plugin (trilead-api)
- planned 2.541.3: **plugin-util-api 6.1192.v30fe6e2837ff** missing-dependency: needs the workflow-support plugin (989.va_20a_1a_57710a_ or newer), which is not installed (takes down 49 more)
- planned 2.541.3: **workflow-cps 4395.v738456a_5da_6b_** missing-dependency: needs the workflow-support plugin (1015.v785e5a_b_b_8b_22 or newer), which is not installed (takes down 26 more)
- planned 2.541.3: **workflow-job 1571.1585.v21fce07cd04d** missing-dependency: needs the workflow-support plugin (976.vb_d9493c2eb_09 or newer), which is not installed (takes down 12 more)
- planned 2.541.3: **workflow-durable-task-step 1479.v56e587f413a_7** missing-dependency: needs the workflow-support plugin (944.v5a_859593b_98a_ or newer), which is not installed (takes down 5 more)
- planned 2.541.3: **git-parameter 462.463.v496a_59f698e5** missing-dependency: needs the workflow-support plugin (968.v8f17397e87b_8 or newer), which is not installed
- planned 2.541.3: **github-branch-source 1986.v801e40614d00** missing-dependency: needs the workflow-support plugin (1010.vb_b_39488a_9841 or newer), which is not installed
- planned 2.541.3: **lockable-resources 1576.v217259b_f9976** missing-dependency: needs the workflow-support plugin (1015.v785e5a_b_b_8b_22 or newer), which is not installed
- planned 2.541.3: **pipeline-build-step 601.v6d4c6d1a_9dc7** missing-dependency: needs the workflow-support plugin (944.v5a_859593b_98a_ or newer), which is not installed
- planned 2.541.3: **pipeline-input-step 560.v56198a_642157** missing-dependency: needs the workflow-support plugin (936.v9fa_77211ca_e1 or newer), which is not installed
- planned 2.541.3: **workflow-support 1015.v785e5a_b_b_8b_22** not-active: installed but not active after startup
- planned 2.541.3: **workflow-support.jpi ** other: java.io.IOException: Failed to expand /var/jenkins_home/plugins/workflow-support.jpi
- planned 2.541.3: **ws-cleanup 0.49** not-active: installed but not active after startup
- planned 2.541.3: **ws-cleanup.jpi ** other: java.nio.file.FileSystemException: /var/jenkins_home/plugins/ws-cleanup: No space left on device
- planned 2.541.3: **xunit 3.1.7** not-active: installed but not active after startup
- planned 2.541.3: **xunit.jpi ** other: java.nio.file.FileSystemException: /var/jenkins_home/plugins/xunit: No space left on device

