### 2.332.4 → 2.580.2, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.332.4 | jenkins/jenkins:2.332.4-jdk11 | started | 302 | 0 | 0 | 0 | 0 |
| as-is: 2.580.2, plugins unchanged | jenkins/jenkins:2.580.2-jdk21 | did not start | 302 | 0 | 0 | 2 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 327 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 350 | 0 | 0 | 0 | 0 |
| planned: 2.580.2, plugins from plan | jenkins/jenkins:2.580.2-jdk21 | started | 351 | 0 | 0 | 0 | 0 |

- baseline: Java 11 because gradle needs Java 11 or newer
- as-is: loaded in part (extensions failed): github, sshd
- as-is: the container exited (state: false 5); the log ends: jenkins.InitReactorRunner$1#onTaskFailed: Failed Contributed.load — java.lang.NullPointerException: Cannot read field "overrides" because the return value of "hudson.ExtensionList.get(java.lang.Class)" is null

