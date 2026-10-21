### 2.263.4 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.263.4 | jenkins/jenkins:2.263.4 | started | 303 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 303 | 0 | 0 | 4 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 337 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 360 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 360 | 0 | 0 | 0 | 0 |

- as-is: loaded in part (extensions failed): copyartifact, git-server, github, scriptler
- as-is: the container exited (state: false 5); the log ends: hudson.util.BootFailure#publish: Failed to initialize Jenkins — java.lang.IllegalArgumentException: Unable to inject class com.atlassian.bitbucket.jenkins.internal.jenkins.auth.OAuthFilterRegistrar

