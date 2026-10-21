### 2.303.3 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.303.3 | jenkins/jenkins:2.303.3-jdk11 | started | 104 | 31 | 31 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 104 | 31 | 11 | 0 | 0 |
| planned: 2.462.3, plugins from plan | jenkins/jenkins:2.462.3-jdk11 | started | 113 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 119 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 119 | 0 | 0 | 0 | 0 |

- baseline: **authentication-tokens 1.4** other: Required plugin is disabled: Credentials Plugin
- baseline: **config-file-provider 3.10.0** other: Required plugin is disabled: Credentials Plugin
- baseline: **credentials 2.6.1.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- baseline: **credentials-binding 1.27.1** other: Required plugin is disabled: Credentials Plugin
- baseline: **docker-commons 1.21** other: Failed to load: Authentication Tokens API Plugin (1.4)
- baseline: **docker-workflow 1.28** other: Failed to load: Docker Commons Plugin (1.21)
- baseline: **email-ext 2.89** other: Required plugin is disabled: Credentials Plugin
- baseline: **git 4.11.5** other: Failed to load: Credentials Binding Plugin (1.27.1)
- baseline: **git-client 3.11.2** other: Failed to load: Jenkins JSch dependency plugin (0.1.55.2)
- baseline: **git-server 1.11** other: Required plugin is disabled: SSH server
- baseline: **github 1.34.3.1** other: Required plugin is disabled: Credentials Plugin
- baseline: **github-branch-source 2.11.4** other: Failed to load: GitHub plugin (1.34.3.1)
- baseline: **gradle 2.12.0.1** other: Required plugin is disabled: Credentials Plugin
- baseline: **jsch 0.1.55.2** other: Required plugin is disabled: Trilead API Plugin
- baseline: **maven-plugin 3.16** other: Failed to load: Jenkins JSch dependency plugin (0.1.55.2)
- baseline: **mina-sshd-api-core 2.8.0-30.vf9df64641cb_d** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- baseline: **pipeline-github-lib 38.v445716ea_edda_** other: Failed to load: Git plugin (4.11.5)
- baseline: **pipeline-graph-analysis 202.va_d268e64deb_3** other: Failed to load: Pipeline: Input Step (466.v6d0a_5df34f81)
- baseline: **pipeline-graph-view 61.v5f35529d13b7** other: Failed to load: Pipeline Graph Analysis Plugin (202.va_d268e64deb_3)
- baseline: **pipeline-input-step 466.v6d0a_5df34f81** other: Required plugin is disabled: Credentials Plugin
- baseline: **pipeline-model-definition 1.9.3** other: Failed to load: Pipeline: Deprecated Groovy Libraries (609.vd95673f149b_b)
- baseline: **pipeline-model-extensions 1.9.3** other: Failed to load: Credentials Binding Plugin (1.27.1)
- baseline: **pipeline-rest-api 2.28** other: Failed to load: Pipeline Graph Analysis Plugin (202.va_d268e64deb_3)
- baseline: **pipeline-stage-view 2.28** other: Failed to load: Pipeline: REST API Plugin (2.28)
- baseline: **plain-credentials 1.8** other: Required plugin is disabled: Credentials Plugin
- baseline: **ssh-credentials 277.v95c2fec1c047** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- baseline: **ssh-slaves 1.806.v2253cedd3295** other: Required plugin is disabled: Trilead API Plugin
- baseline: **sshd 3.242.va_db_9da_b_26a_c3** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- baseline: **trilead-api 1.0.13** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:trilead-api -> Plugin:sshd
- baseline: **workflow-aggregator 596.v8c21c963d92d** other: Failed to load: Pipeline: Input Step (466.v6d0a_5df34f81)
- baseline: **workflow-cps-global-lib 609.vd95673f149b_b** other: Failed to load: Jenkins GIT server Plugin (1.11)
- as-is: **plain-credentials 1.8** other: Required plugin is disabled: Credentials Plugin (credentials) (takes down 8 more)
- as-is: **jsch 0.1.55.2** other: Required plugin is disabled: SSH Credentials Plugin (ssh-credentials) (takes down 5 more)
- as-is: **pipeline-input-step 466.v6d0a_5df34f81** other: Required plugin is disabled: Credentials Plugin (credentials) (takes down 5 more)
- as-is: **authentication-tokens 1.4** other: Required plugin is disabled: Credentials Plugin (credentials) (takes down 2 more)
- as-is: **credentials 2.6.1.1** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- as-is: **email-ext 2.89** other: Required plugin is disabled: Credentials Plugin (credentials)
- as-is: **mina-sshd-api-core 2.8.0-30.vf9df64641cb_d** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- as-is: **ssh-credentials 277.v95c2fec1c047** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- as-is: **ssh-slaves 1.806.v2253cedd3295** other: Required plugin is disabled: SSH Credentials Plugin (ssh-credentials)
- as-is: **sshd 3.242.va_db_9da_b_26a_c3** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:sshd
- as-is: **trilead-api 1.0.13** other: hudson.util.CyclicGraphDetector$CycleDetectedException: Cycle detected: Plugin:sshd -> Plugin:mina-sshd-api-core -> Plugin:ssh-credentials -> Plugin:credentials -> Plugin:trilead-api -> Plugin:sshd

