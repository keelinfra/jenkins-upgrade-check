### 2.249.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.249.3 | jenkins/jenkins:2.249.3 | started | 303 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 303 | 0 | 0 | 7 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 341 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 364 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 364 | 0 | 0 | 1 | 0 |

- as-is: loaded in part (extensions failed): copyartifact, email-ext, git-server, github, mailer, scriptler, translation
- as-is: the container exited (state: false 5); the log ends: hudson.util.BootFailure#publish: Failed to initialize Jenkins — java.lang.IllegalArgumentException: Unable to inject class com.atlassian.bitbucket.jenkins.internal.jenkins.auth.OAuthFilterRegistrar
- planned 2.580.1: loaded in part (extensions failed): translation

