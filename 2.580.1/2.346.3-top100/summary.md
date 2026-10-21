### 2.346.3 → 2.580.1, 100 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.346.3 | jenkins/jenkins:2.346.3-jdk11 | started | 103 | 2 | 2 | 0 | 0 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 103 | 0 | 0 | 0 | 0 |
| planned: 2.541.3, plugins from plan | jenkins/jenkins:2.541.3-jdk17 | started | 116 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 116 | 0 | 0 | 0 | 0 |

- baseline: **metrics 4.2.13-420.vea_2f17932dd6** other: java.lang.NullPointerException
- baseline: **pipeline-graph-view 144.vf3924feb_7e35** missing-dependency: needs the metrics plugin (4.2.10-389.v93143621b_050 or newer), which is not installed

