### 2.249.3 → 2.580.2, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.249.3 | jenkins/jenkins:2.249.3 | started | 303 | 0 | 0 | 0 | 0 |
| as-is: 2.580.2, plugins unchanged | jenkins/jenkins:2.580.2-jdk21 | did not start | 303 | 0 | 0 | 7 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 342 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 365 | 0 | 0 | 0 | 0 |
| planned: 2.580.2, plugins from plan | jenkins/jenkins:2.580.2-jdk21 | started | 366 | 0 | 0 | 0 | 0 |

- as-is: loaded in part (extensions failed): copyartifact, email-ext, git-server, github, mailer, scriptler, translation
- as-is: the container exited (state: false 5); the log ends: jenkins.InitReactorRunner$1#onTaskFailed: Failed KubernetesCredentialProvider.startWatchingForSecrets — java.lang.IllegalArgumentException: Unable to inject class com.cloudbees.jenkins.plugins.kubernetes_credentials_provider.KubernetesCredentialProvider

