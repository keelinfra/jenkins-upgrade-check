### 2.164.3 → 2.580.2, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.164.3 | jenkins/jenkins:2.164.3 | started | 304 | 0 | 0 | 0 | 0 |
| as-is: 2.580.2, plugins unchanged | jenkins/jenkins:2.580.2-jdk21 | did not start | 304 | 1 | 1 | 10 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 364 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 387 | 0 | 0 | 0 | 0 |
| planned: 2.580.2, plugins from plan | jenkins/jenkins:2.580.2-jdk21 | started | 390 | 0 | 0 | 0 | 0 |

- as-is: **ec2 1.50.2.1** other: Required plugin is disabled: Trilead API Plugin (trilead-api)
- as-is: loaded in part (extensions failed): bitbucket-push-and-pull-request, copyartifact, email-ext, git-server, github, ldap, mailer, scriptler, ssh-slaves, translation
- as-is: the container exited (state: false 5); the log ends: jenkins.InitReactorRunner$1#onTaskFailed: Failed PluginImpl.init — java.lang.IllegalStateException: Expected 1 instance of hudson.PluginManager$PluginUpdateMonitor but got 0

