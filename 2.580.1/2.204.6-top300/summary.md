### 2.204.6 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.204.6 | jenkins/jenkins:2.204.6 | started | 304 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 304 | 0 | 0 | 8 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 361 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 384 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 386 | 0 | 0 | 0 | 0 |

- as-is: loaded in part (extensions failed): copyartifact, email-ext, git-server, github, ldap, mailer, scriptler, translation
- as-is: the container exited (state: false 5); the log ends: jenkins.InitReactorRunner$1#onTaskFailed: Failed KubernetesCredentialProvider.startWatchingForSecrets — java.lang.IllegalArgumentException: Unable to inject class com.cloudbees.jenkins.plugins.kubernetes_credentials_provider.KubernetesCredentialProvider

