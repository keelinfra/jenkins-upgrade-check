### 2.346.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.346.3 | jenkins/jenkins:2.346.3-jdk11 | started | 302 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 302 | 0 | 0 | 2 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 349 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 349 | 0 | 0 | 0 | 0 |

- baseline: Java 11 because gradle needs Java 11 or newer
- as-is: loaded in part (extensions failed): github, sshd
- as-is: the container exited (state: false 5); the log ends: hudson.util.BootFailure#publish: Failed to initialize Jenkins — java.lang.IllegalArgumentException: Unable to inject class com.atlassian.bitbucket.jenkins.internal.jenkins.auth.OAuthFilterRegistrar

