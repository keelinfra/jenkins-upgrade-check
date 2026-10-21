### 2.176.4 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.176.4 | jenkins/jenkins:2.176.4 | started | 304 | 0 | 0 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | did not start | 304 | 1 | 1 | 9 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 361 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 384 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 386 | 0 | 0 | 1 | 0 |

- as-is: **ec2 1.50.3** other: Required plugin is disabled: Trilead API Plugin (trilead-api)
- as-is: loaded in part (extensions failed): copyartifact, email-ext, git-server, github, ldap, mailer, scriptler, ssh-slaves, translation
- as-is: the container exited (state: false 5); the log ends: hudson.util.BootFailure#publish: Failed to initialize Jenkins — java.lang.NullPointerException: Cannot invoke "hudson.plugins.git.GitSCM$DescriptorImpl.getOldGitExe()" because "desc" is null
- planned 2.580.1: loaded in part (extensions failed): translation

