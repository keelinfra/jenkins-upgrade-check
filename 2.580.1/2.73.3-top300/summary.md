### 2.73.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.73.3 | jenkins/jenkins:2.73.3 | started | 305 | 6 | 6 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 305 | 6 | 5 | 0 | 0 |
| planned: 2.346.3, plugins from plan | jenkins/jenkins:2.346.3-jdk8 | started | 351 | 2 | 1 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | did not start | 402 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan |  | did not start | 0 | 0 | 0 | 0 | 0 |

- baseline: **Parameterized-Remote-Trigger 3.1.6.3** other: Credentials Plugin v2.1.19 is older than required. To fix, install v2.6.1.1 or later.
- baseline: **artifactory 3.5.0** other: Credentials Plugin v2.1.19 is older than required. To fix, install v2.2.1 or later.
- baseline: **cucumber-reports 4.11.0** other: Token Macro Plugin v2.6 is older than required. To fix, install v2.8 or later.
- baseline: **docker-build-step 2.8** other: Docker Commons Plugin v1.14 is older than required. To fix, install v1.15 or later.
- baseline: **docker-java-api 3.2.13-37.vf3411c9828b9** other: Jackson 2 API Plugin v2.9.9.1 is older than required. To fix, install v2.10.3 or later.
- baseline: **docker-plugin 1.2.2** other: Docker API Plugin v3.2.13-37.vf3411c9828b9 failed to load. Fix this plugin first.
- as-is: **docker-java-api 3.2.13-37.vf3411c9828b9** old-dependency: needs jackson2-api 2.10.3 or newer; 2.9.9.1 is installed (takes down 1 more)
- as-is: **Parameterized-Remote-Trigger 3.1.6.3** old-dependency: needs credentials 2.6.1.1 or newer; 2.1.19 is installed
- as-is: **artifactory 3.5.0** old-dependency: needs credentials 2.2.1 or newer; 2.1.19 is installed
- as-is: **cucumber-reports 4.11.0** old-dependency: needs token-macro 2.8 or newer; 2.6 is installed
- as-is: **docker-build-step 2.8** old-dependency: needs docker-commons 1.15 or newer; 1.14 is installed
- planned 2.346.3: **gradle 2.12.0.1** java-too-old: this release needs Java 11 or newer (takes down 1 more)

