# Jenkins upgrade plan: 2.138.4 → 2.580.1

Source: results/2.138.4-top300/plugins-2.138.4.txt · generated 2026-10-21 12:52 UTC by jenkins-upgrade-check dev

**3 core upgrade(s)**, 391 plugins: 291 to upgrade, 87 new dependencies; 120 with security warnings today, 4 still with warnings after the upgrade, 3 deprecated, 3 not from the update center, 4 dependency problem(s).

Jenkins 2.138.4 runs on Java 8; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

- [core-2_160](https://jenkins.io/security/advisory/2019-01-16/): Multiple security vulnerabilities in Jenkins 2.159 and earlier, and LTS 2.150.1 and earlier
- [core-2_172](https://jenkins.io/security/advisory/2019-04-10/): Multiple security vulnerabilities in Jenkins 2.171 and earlier, and LTS 2.164.1 and earlier
- [core-2_186](https://jenkins.io/security/advisory/2019-07-17/): Multiple security vulnerabilities in Jenkins 2.185 and earlier, and LTS 2.176.1 and earlier
- [core-2_192](https://jenkins.io/security/advisory/2019-08-28/): Multiple security vulnerabilities in Jenkins 2.191 and earlier, and LTS 2.176.2 and earlier
- [core-2_197](https://jenkins.io/security/advisory/2019-09-25/): Multiple security vulnerabilities in Jenkins 2.196 and earlier, and LTS 2.176.3 and earlier
- [core-2_219](https://jenkins.io/security/advisory/2020-01-29/): Multiple vulnerabilities in Jenkins 2.218 and earlier, LTS 2.204.1 and earlier
- [core-2_228](https://jenkins.io/security/advisory/2020-03-25/): Multiple security vulnerabilities in Jenkins 2.227 and earlier, LTS 2.204.5 and earlier
- [core-2_245](https://jenkins.io/security/advisory/2020-07-15/): Multiple security vulnerabilities in Jenkins 2.244 and earlier, LTS 2.235.1 and earlier
- [core-2_252](https://jenkins.io/security/advisory/2020-08-12/): Multiple security vulnerabilities in Jenkins 2.251 and earlier, LTS 2.235.3 and earlier
- [core-2_275](https://www.jenkins.io/security/advisory/2021-01-13/): Multiple security vulnerabilities in Jenkins 2.274 and earlier, LTS 2.263.1 and earlier
- [core-2_287](https://www.jenkins.io/security/advisory/2021-04-07/): Multiple security vulnerabilities in Jenkins 2.286 and earlier, LTS 2.277.1 and earlier
- [core-2_286](https://www.jenkins.io/security/advisory/2021-04-20/): Denial of service vulnerability in bundled Jetty
- [core-2_300](https://www.jenkins.io/security/advisory/2021-06-30/): Multiple security vulnerabilities in Jenkins 2.299 and earlier, LTS 2.289.1 and earlier
- [core-2_315](https://www.jenkins.io/security/advisory/2021-10-06/): Multiple security vulnerabilities in Jenkins 2.314 and earlier, LTS 2.303.1 and earlier
- [core-2_319](https://www.jenkins.io/security/advisory/2021-11-04/): Multiple security vulnerabilities in Jenkins 2.318 and earlier, LTS 2.303.2 and earlier
- [core-2_330](https://www.jenkins.io/security/advisory/2022-01-12/): CSRF vulnerability in build triggers
- [core-2_334](https://www.jenkins.io/security/advisory/2022-02-09/): DoS vulnerability in bundled XStream library
- [core-2_356](https://www.jenkins.io/security/advisory/2022-06-22/): Multiple security vulnerabilities in Jenkins 2.355 and earlier, LTS 2.332.3 and earlier
- [core-2_263](https://www.jenkins.io/security/advisory/2022-09-09/): HTTP/2 denial of service vulnerability in bundled Jetty
- [core-2_394](https://www.jenkins.io/security/advisory/2023-03-08/): Multiple security vulnerabilities in Jenkins 2.393 and earlier, LTS 2.375.3 and earlier
- [core-2_400](https://www.jenkins.io/security/advisory/2023-06-14/): CSRF bypass vulnerability
- [core-2_416](https://www.jenkins.io/security/advisory/2023-07-26/): Stored XSS vulnerability
- [core-2_424](https://www.jenkins.io/security/advisory/2023-09-20/): Multiple security vulnerabilities in Jenkins 2.423 and earlier, LTS 2.414.1 and earlier
- [core-2_428](https://www.jenkins.io/security/advisory/2023-10-18/): HTTP/2 denial of service vulnerabilities in bundled Jetty
- [core-2_442](https://www.jenkins.io/security/advisory/2024-01-24/): Multiple security vulnerabilities in Jenkins 2.441 and earlier, LTS 2.426.2 and earlier
- [core-2_444](https://www.jenkins.io/security/advisory/2024-03-20/): HTTP/2 denial of service vulnerability in bundled Jetty
- [core-2_452](https://www.jenkins.io/security/advisory/2024-04-17/): Terrapin SSH vulnerability in Jenkins CLI client
- [core-2_471](https://www.jenkins.io/security/advisory/2024-08-07/): Multiple security vulnerabilities in Jenkins 2.470 and earlier, LTS 2.452.3 and earlier
- [core-2_479](https://www.jenkins.io/security/advisory/2024-10-02/): Multiple security vulnerabilities in Jenkins 2.478 and earlier, LTS 2.462.2 and earlier
- [core-2_487](https://www.jenkins.io/security/advisory/2024-11-27/): Denial of service vulnerability in bundled json-lib
- [core-2_500](https://www.jenkins.io/security/advisory/2025-03-05/): Multiple security vulnerabilities in Jenkins 2.499 and earlier, LTS 2.492.1 and earlier
- [core-2_504](https://www.jenkins.io/security/advisory/2025-04-02/): Multiple security vulnerabilities in Jenkins 2.503 and earlier, LTS 2.492.2 and earlier
- [core-2_528](https://www.jenkins.io/security/advisory/2025-09-17/): Multiple security vulnerabilities in Jenkins 2.527 and earlier, LTS 2.516.2 and earlier
- [core-2_541](https://www.jenkins.io/security/advisory/2025-12-10/): Multiple security vulnerabilities in Jenkins 2.540 and earlier, LTS 2.528.2 and earlier
- [core-2_551](https://www.jenkins.io/security/advisory/2026-02-18/): Multiple security vulnerabilities in Jenkins 2.550 and earlier, LTS 2.541.1 and earlier
- [core-2_555](https://www.jenkins.io/security/advisory/2026-03-18/): Multiple security vulnerabilities in Jenkins 2.554 and earlier, LTS 2.541.2 and earlier
- [core-2_568](https://www.jenkins.io/security/advisory/2026-06-10/): Multiple security vulnerabilities in Jenkins 2.567 and earlier, LTS 2.555.2 and earlier
- [core-2_576](https://www.jenkins.io/security/advisory/2026-08-05/): Multiple security vulnerabilities in Jenkins 2.575 and earlier, LTS 2.568.1 and earlier
- [core-2_580](https://www.jenkins.io/security/advisory/2026-09-02/): Multiple security vulnerabilities in Jenkins 2.579 and earlier, LTS 2.568.2 and earlier

## Steps

1. **On 2.138.4, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
   - ⚠ Parameterized-Remote-Trigger 3.1.6.3 needs credentials 2.6.1.1 or newer, but the newest credentials that runs on Jenkins 2.138.4 is 2.3.11
   - ⚠ Parameterized-Remote-Trigger 3.1.6.3 needs script-security 1.75 or newer, but the newest script-security that runs on Jenkins 2.138.4 is 1.71
   - ⚠ docker-build-step 2.8 needs credentials 2.3.19 or newer, but the newest credentials that runs on Jenkins 2.138.4 is 2.3.11
   - ⚠ docker-build-step 2.8 needs docker-commons 1.15 or newer, but the newest docker-commons that runs on Jenkins 2.138.4 is 1.14
2. **Run the controller and every agent on Java 8**, upgrade Jenkins to **2.346.3**, start it, then update plugins (293 change(s)).
3. **Run the controller and every agent on Java 17**, upgrade Jenkins to **2.541.3**, start it, then update plugins (361 change(s)).
4. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (34 change(s)).

No Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| Parameterized-Remote-Trigger | 3.1.6.3 | 3.2.2 | [SECURITY-3768](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3768) **no fix yet**<br>[SECURITY-3800](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3800) **no fix yet** |
| ace-editor | 1.0.1 | 1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| active-directory | 2.16.1 | 2.971.v653d8b_6e5548 | [SECURITY-1389](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-1389) fixed by the upgrade<br>[SECURITY-3059](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3059) fixed by the upgrade<br>[SECURITY-3659](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3659) fixed by the upgrade<br>[SECURITY-3651](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3651) fixed by the upgrade |
| allure-jenkins-plugin | 2.30.3 | 3.0.0 | [SECURITY-3645](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3645) fixed by the upgrade |
| analysis-model-api | 8.0.1 | 14.18.0-1028.v9c938eeda_ef2 | [SECURITY-3327](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3327) fixed by the upgrade |
| ansible | 1.1 | 635.v34b_48e979c86 | [SECURITY-3017](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3017) fixed by the upgrade |
| artifactory | 3.5.0 | 4.0.8 | [SECURITY-1542-1](https://jenkins.io/security/advisory/2020-03-25/#SECURITY-1542%20(1)) fixed by the upgrade<br>[SECURITY-1542-2](https://jenkins.io/security/advisory/2020-03-25/#SECURITY-1542%20(2)) fixed by the upgrade |
| audit-trail | 3.4 | 464.v78b_845ffc5d8 | [SECURITY-1815](https://www.jenkins.io/security/advisory/2020-10-08/#SECURITY-1815) fixed by the upgrade<br>[SECURITY-1846](https://www.jenkins.io/security/advisory/2020-10-08/#SECURITY-1846) fixed by the upgrade |
| authorize-project | 1.3.0 | 534.v2f208c45e11c | [SECURITY-3010](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3010) fixed by the upgrade |
| azure-ad | 1.2.1 | 720.v2a_e3d48ffeb_7 | [SECURITY-2470](https://www.jenkins.io/security/advisory/2021-08-31/#SECURITY-2470) fixed by the upgrade<br>[SECURITY-2980](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2980) fixed by the upgrade<br>[SECURITY-3227](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3227) fixed by the upgrade<br>[SECURITY-3760](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3760) fixed by the upgrade<br>[SECURITY-3935](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3935) fixed by the upgrade |
| azure-credentials | 2.0.2 | 426.v5627e6a_4c7a_7 | [SECURITY-1756](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-1756) fixed by the upgrade<br>[SECURITY-1757](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-1757) fixed by the upgrade |
| badge | 1.8 | 3.615.vb_ed967ca_b_e93 | [SECURITY-2547](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2547) fixed by the upgrade |
| bitbucket-push-and-pull-request | 2.8.4 | 4.1.0 | [SECURITY-3856](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3856) fixed by the upgrade<br>[SECURITY-3980](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3980) fixed by the upgrade |
| blueocean | 1.22.0 | 1.27.26 | [SECURITY-1956](https://www.jenkins.io/security/advisory/2020-09-16/#SECURITY-1956) fixed by the upgrade<br>[SECURITY-1961](https://www.jenkins.io/security/advisory/2020-09-16/#SECURITY-1961) fixed by the upgrade<br>[SECURITY-2502](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2502) fixed by the upgrade<br>[SECURITY-3116](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3116) fixed by the upgrade |
| blueocean-pipeline-scm-api | 1.22.0 | 1.27.26 | [SECURITY-714](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-714) fixed by the upgrade |
| build-failure-analyzer | 1.26.0 | 3.889.vda_d3b_30435fc | [SECURITY-1770](https://jenkins.io/security/advisory/2020-09-01/#SECURITY-1770) fixed by the upgrade<br>[SECURITY-3226](https://www.jenkins.io/security/advisory/2023-09-20/#SECURITY-3226) fixed by the upgrade<br>[SECURITY-3239](https://www.jenkins.io/security/advisory/2023-09-20/#SECURITY-3239) fixed by the upgrade<br>[SECURITY-3244](https://www.jenkins.io/security/advisory/2023-09-20/#SECURITY-3244) fixed by the upgrade |
| build-monitor-plugin | 1.12+build.201809061734 | 1131.va_2d8e68fc912 | [SECURITY-3280](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3280) fixed by the upgrade |
| build-pipeline-plugin | 1.5.8 | 2.2.0 | [SECURITY-879](https://jenkins.io/security/advisory/2019-08-07/#SECURITY-879) fixed by the upgrade |
| cloudbees-bitbucket-branch-source | 2.9.1 | 937.3.10 | [SECURITY-2033](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2033) fixed by the upgrade<br>[SECURITY-2467](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2467) fixed by the upgrade<br>[SECURITY-3300](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3300) fixed by the upgrade<br>[SECURITY-3363](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3363) fixed by the upgrade |
| cloudbees-folder | 6.12 | 6.1106.v3a_d9a_6d2465e | [SECURITY-3105](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3105) fixed by the upgrade<br>[SECURITY-3106](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3106) fixed by the upgrade<br>[SECURITY-3109](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3109) fixed by the upgrade |
| code-coverage-api | 1.1.6 | 4.99.0 | [SECURITY-2376](https://www.jenkins.io/security/advisory/2021-08-31/#SECURITY-2376) fixed by the upgrade |
| config-file-provider | 3.7.2 | 1013.v73c323e52b_1f | [SECURITY-3090](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3090) fixed by the upgrade |
| configuration-as-code | 1.36.2 | 2136.veb_3a_554550b_f | [SECURITY-2141](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2141) fixed by the upgrade |
| configurationslicing | 1.47 | 715.vcc64f669003e | [SECURITY-2003](https://www.jenkins.io/security/advisory/2021-02-24/#SECURITY-2003) fixed by the upgrade |
| credentials | 2.3.11 | 1514.vb_01eb_014e7d4 | [SECURITY-2349](https://www.jenkins.io/security/advisory/2021-05-11/#SECURITY-2349) fixed by the upgrade<br>[SECURITY-2617-credentials](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-3373](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3373) fixed by the upgrade |
| credentials-binding | 1.18 | 728.v902a_273b_8947 | [SECURITY-1374](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-1374) fixed by the upgrade<br>[SECURITY-1835](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-1835) fixed by the upgrade<br>[SECURITY-2342](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2342) fixed by the upgrade<br>[SECURITY-3499](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3499) fixed by the upgrade<br>[SECURITY-3672](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3672) fixed by the upgrade<br>[SECURITY-3790](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3790) fixed by the upgrade |
| cvs | 2.17 | 502.v02d5a_0c6b_9a_3 | [SECURITY-2617-cvs](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade |
| dashboard-view | 2.12.1 | 2.569.v5ddb_c9b_a_7264 | [SECURITY-2559](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2559) fixed by the upgrade |
| dependency-check-jenkins-plugin | 5.1.2 | 5.6.5 | [SECURITY-3344](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3344) fixed by the upgrade<br>[SECURITY-3992](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3992) fixed by the upgrade |
| dingding-notifications | 1.9 | 2.9.1 | [SECURITY-1423](https://jenkins.io/security/advisory/2019-10-01/#SECURITY-1423) fixed by the upgrade<br>[SECURITY-3353](https://www.jenkins.io/security/advisory/2025-05-14/#SECURITY-3353) fixed by the upgrade |
| docker-build-step | 2.8 | 2.14 | [SECURITY-3200](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3200) fixed by the upgrade |
| docker-commons | 1.14 | 477.v289085a_b_6896 | [SECURITY-1878](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-1878) fixed by the upgrade |
| ec2 | 1.44.1 | 2064.v79de4fd2d03b_ | [SECURITY-1004](https://jenkins.io/security/advisory/2020-01-15/#SECURITY-1004) fixed by the upgrade<br>[SECURITY-1408](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-1408) fixed by the upgrade<br>[SECURITY-1528](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-1528) fixed by the upgrade<br>[SECURITY-1844](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-1844) fixed by the upgrade<br>[SECURITY-381](https://jenkins.io/security/advisory/2020-05-06/#SECURITY-381) fixed by the upgrade |
| echarts-api | 4.6.0-8 | 6.1.0-1321.v86c4c5114494 | [SECURITY-1841](https://jenkins.io/security/advisory/2020-06-03/#SECURITY-1841) fixed by the upgrade<br>[SECURITY-1842](https://jenkins.io/security/advisory/2020-06-03/#SECURITY-1842) fixed by the upgrade |
| email-ext | 2.69.2 | 2038.v7b_8817a_499d9 | [SECURITY-2931](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2931) fixed by the upgrade<br>[SECURITY-2934](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2934) fixed by the upgrade<br>[SECURITY-2939](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2939) fixed by the upgrade<br>[SECURITY-3088-1](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(1)) fixed by the upgrade<br>[SECURITY-3088-2](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(2)) fixed by the upgrade<br>[SECURITY-3705](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3705) fixed by the upgrade |
| extended-choice-parameter | 346.vd87693c5a_86c | 388.ve7b_d0b_920e10 | [SECURITY-1350](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1350) **no fix yet**<br>[SECURITY-1351](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1351) **no fix yet**<br>[SECURITY-2232](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2232) **no fix yet**<br>[SECURITY-2617-extended-choice-parameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) **no fix yet**<br>[deprecated](https://github.com/jenkinsci/extended-choice-parameter-plugin?tab=readme-ov-file#end-of-life) |
| extensible-choice-parameter | 1.8.0 | 255.vfa_41f46828ef | [SECURITY-3583](https://www.jenkins.io/security/advisory/2025-10-29/#SECURITY-3583) fixed by the upgrade |
| external-monitor-job | 1.7.1 | 223.vb_fddcf42c9b_3 | [SECURITY-3133](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3133) fixed by the upgrade |
| favorite | 2.3.2 | 2.267.vb_90d08408081 | [SECURITY-2557](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2557) fixed by the upgrade |
| generic-webhook-trigger | 1.76 | 2.4.3 | [SECURITY-2592](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2592) fixed by the upgrade<br>[SECURITY-2874](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2874) fixed by the upgrade |
| gerrit-trigger | 2.34.0 | 3.1998.v9487a_a_89a_eb_1 | [SECURITY-2617-gerrit-trigger](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-2137](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2137) fixed by the upgrade |
| ghprb | 1.42.2 | 1.42.2 | [SECURITY-2789-1](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2789%20(1)) **no fix yet**<br>[SECURITY-2789-2](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2789%20(2)) **no fix yet** |
| git | 4.2.2 | 5.11.0 | [SECURITY-2499](https://www.jenkins.io/security/advisory/2021-10-06/#SECURITY-2499) fixed by the upgrade<br>[SECURITY-2478-git](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2478) fixed by the upgrade<br>[SECURITY-284](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-284) fixed by the upgrade<br>[SECURITY-2796](https://www.jenkins.io/security/advisory/2022-08-23/#SECURITY-2796) fixed by the upgrade |
| git-client | 3.2.1 | 6.6.1 | [SECURITY-1468](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-1468) fixed by the upgrade<br>[SECURITY-3590](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3590) fixed by the upgrade<br>[SECURITY-3614](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3614) fixed by the upgrade<br>[SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| git-parameter | 0.9.13 | 462.463.v496a_59f698e5 | [SECURITY-2617-git-parameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-3419](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3419) fixed by the upgrade<br>[SECURITY-3745](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3745) fixed by the upgrade |
| git-server | 1.9 | 137.ve0060b_432302 | [SECURITY-3319](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3319) fixed by the upgrade<br>[SECURITY-3342](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3342) fixed by the upgrade |
| gitea | 1.3.0 | 286.v0740a_413c1eb_ | [SECURITY-2661](https://www.jenkins.io/security/advisory/2022-12-07/#SECURITY-2661) fixed by the upgrade |
| github | 1.30.0 | 1.47.0 | [SECURITY-1849](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-1849) fixed by the upgrade<br>[SECURITY-3246](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3246) fixed by the upgrade<br>[SECURITY-3704](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3704) fixed by the upgrade |
| github-branch-source | 2.5.8 | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| github-pullrequest | 0.2.8 | 0.7.5 | [SECURITY-3776](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3776) fixed by the upgrade |
| gitlab-oauth | 1.11 | 1.30 | [SECURITY-1833](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-1833) fixed by the upgrade<br>[SECURITY-1891](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1891) fixed by the upgrade<br>[SECURITY-2696](https://www.jenkins.io/security/advisory/2023-07-26/#SECURITY-2696) fixed by the upgrade |
| gitlab-plugin | 1.5.13 | 1.2158.v49db_f64b_fc28 | [SECURITY-2753](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2753) fixed by the upgrade<br>[SECURITY-2316](https://www.jenkins.io/security/advisory/2022-06-30/#SECURITY-2316) fixed by the upgrade<br>[SECURITY-2877](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2877) fixed by the upgrade<br>[SECURITY-3260](https://www.jenkins.io/security/advisory/2025-01-22/#SECURITY-3260) fixed by the upgrade<br>[SECURITY-4097](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4097) fixed by the upgrade<br>[SECURITY-3851](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3851) fixed by the upgrade |
| global-build-stats | 1.5 | 360.vd08338d06c53 | [SECURITY-1886](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1886) fixed by the upgrade<br>[SECURITY-3535](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3535) fixed by the upgrade |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| greenballs | 1.15 | 1.15.1 | [deprecated](https://github.com/jenkins-infra/update-center2/pull/735) |
| handlebars | 1.1.1 | 1.1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| hashicorp-vault-plugin | 3.4.0 | 384.vda_86ec66c537 | [SECURITY-2213](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2213) fixed by the upgrade<br>[SECURITY-2429](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2429) fixed by the upgrade<br>[SECURITY-2521](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2521) fixed by the upgrade<br>[SECURITY-2593](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-2593) fixed by the upgrade<br>[SECURITY-3077](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3077) fixed by the upgrade<br>[SECURITY-3045](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3045) fixed by the upgrade |
| hidden-parameter | 0.0.4 | 546.vdeec4ce8b_26f | [SECURITY-2784-hidden-parameter](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2784) fixed by the upgrade |
| htmlpublisher | 1.21 | 429 | [SECURITY-3301](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3301) fixed by the upgrade<br>[SECURITY-3302](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3302) fixed by the upgrade<br>[SECURITY-3303](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3303) fixed by the upgrade<br>[SECURITY-3547](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3547) fixed by the upgrade<br>[SECURITY-3706](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3706) fixed by the upgrade |
| http_request | 1.8.26 | 1.659.v1b_9d9e942a_e6 | [SECURITY-2053](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-2053) fixed by the upgrade |
| ivy | 1.28 | 609.v66779e68503d | [SECURITY-2924](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-2924) fixed by the upgrade<br>[SECURITY-3093](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3093) fixed by the upgrade |
| jacoco | 3.0.4 | 3.3.7 | [SECURITY-3061](https://www.jenkins.io/security/advisory/2023-03-21/#SECURITY-3061) fixed by the upgrade |
| jenkins-multijob-plugin | 1.36 | 682.v4ed91b_52d621 | [SECURITY-3781](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3781) fixed by the upgrade<br>[SECURITY-3823-1](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3823%20(1)) fixed by the upgrade<br>[SECURITY-3823-2](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3823%20(2)) fixed by the upgrade |
| jira | 3.0.13 | 3.22 | [SECURITY-2617-jira](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-3225](https://www.jenkins.io/security/advisory/2023-11-29/#SECURITY-3225) fixed by the upgrade |
| jira-steps | 1.5.5 | 2.1.197.vdd1b_3987e040 | [SECURITY-2774](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2774) fixed by the upgrade<br>[SECURITY-2786](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2786) fixed by the upgrade |
| job-import-plugin | 3.6 | 143.145.v48f9a_a_6ff384 | [SECURITY-3783](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3783) fixed by the upgrade |
| jobConfigHistory | 2.25 | 1380.v762185b_9a_793 | [SECURITY-2766](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-2766) fixed by the upgrade<br>[SECURITY-2765](https://www.jenkins.io/security/advisory/2022-08-23/#SECURITY-2765) fixed by the upgrade<br>[SECURITY-3233](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3233) fixed by the upgrade<br>[SECURITY-3235](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3235) fixed by the upgrade<br>[SECURITY-3742](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3742) fixed by the upgrade<br>[SECURITY-4098](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4098) fixed by the upgrade |
| junit | 1.28 | 1442.vb_b_a_467df2767 | [SECURITY-2760](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2760) fixed by the upgrade<br>[SECURITY-2888](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2888) fixed by the upgrade<br>[SECURITY-3032](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-3032) fixed by the upgrade |
| kubernetes | 1.16.7 | 4560.vc25b_83890c34 | [SECURITY-1646](https://www.jenkins.io/security/advisory/2020-11-04/#SECURITY-1646) fixed by the upgrade<br>[SECURITY-2102](https://www.jenkins.io/security/advisory/2020-11-04/#SECURITY-2102) fixed by the upgrade<br>[SECURITY-2103](https://www.jenkins.io/security/advisory/2020-11-04/#SECURITY-2103) fixed by the upgrade<br>[SECURITY-3075-kubernetes](https://www.jenkins.io/security/advisory/2023-04-12/#SECURITY-3075) fixed by the upgrade |
| kubernetes-cli | 1.7.0 | 1.375.v8001982544ba | [SECURITY-2370](https://www.jenkins.io/security/advisory/2021-06-10/#SECURITY-2370) fixed by the upgrade |
| kubernetes-credentials-provider | 0.12.1 | 1.315.v92008589c044 | [SECURITY-3022](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-3022) fixed by the upgrade |
| ldap | 1.25 | 825.v2fca_37dd5b_cb_ | [SECURITY-3046](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3046) fixed by the upgrade<br>[SECURITY-3654](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3654) fixed by the upgrade<br>[SECURITY-3678](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3678) fixed by the upgrade |
| mailer | 1.23 | 534.v1b_36f5864073 | [SECURITY-1813](https://www.jenkins.io/security/advisory/2020-09-16/#SECURITY-1813) fixed by the upgrade<br>[SECURITY-2163](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2163) fixed by the upgrade |
| mask-passwords | 2.13 | 220.v95819055b_265 | [SECURITY-2617-mask-passwords](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade |
| matrix-auth | 2.5.2 | 3.3 | [SECURITY-3676](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3676) fixed by the upgrade |
| matrix-project | 1.14.1 | 928.v1effc3f3a_236 | [SECURITY-2017](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2017) fixed by the upgrade<br>[SECURITY-3289](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3289) fixed by the upgrade |
| mercurial | 2.8.1 | 1323.ve69d2a_db_8a_b_d | [SECURITY-2478-mercurial](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2478) fixed by the upgrade<br>[SECURITY-2831](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2831) fixed by the upgrade |
| metrics | 4.0.2.6 | 4.2.37-494.v06f9a_939d33a_ | [SECURITY-1624](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-1624) fixed by the upgrade |
| momentjs | 1.1 | 1.1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| multibranch-scan-webhook-trigger | 1.0.5 | 1.0.11 | [SECURITY-2875](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-2875) fixed by the upgrade |
| nested-view | 1.17 | 241.v178f0b_a_cd76a_ | [SECURITY-2411](https://www.jenkins.io/security/advisory/2021-08-31/#SECURITY-2411) fixed by the upgrade<br>[SECURITY-2768](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2768) fixed by the upgrade |
| nodejs | 1.3.5 | 1.6.6 | [SECURITY-3196](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3196) fixed by the upgrade |
| nodelabelparameter | 1.7.2 | 851.vd94e5048d321 | [SECURITY-2617-nodelabelparameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade |
| oic-auth | 1.8 | 4.731.v4932b_506e767 | [SECURITY-2978](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2978) fixed by the upgrade<br>[SECURITY-2979](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-2979) fixed by the upgrade<br>[SECURITY-3168](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3168) fixed by the upgrade<br>[SECURITY-3441-1](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3441%20(1)) fixed by the upgrade<br>[SECURITY-3441-2](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3441%20(2)) fixed by the upgrade<br>[SECURITY-3473](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3473) fixed by the upgrade<br>[SECURITY-3461](https://www.jenkins.io/security/advisory/2025-01-22/#SECURITY-3461) fixed by the upgrade |
| parameterized-trigger | 2.36 | 893.va_383a_9b_a_4c11 | [SECURITY-2185](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2185) fixed by the upgrade |
| performance | 3.20 | 1022.v48fce959f0c3 | [SECURITY-2394](https://www.jenkins.io/security/advisory/2021-11-12/#SECURITY-2394) fixed by the upgrade<br>[SECURITY-4026](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4026) fixed by the upgrade |
| pipeline-build-step | 2.12 | 601.v6d4c6d1a_9dc7 | [SECURITY-2519](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2519) fixed by the upgrade<br>[SECURITY-3019](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-3019) fixed by the upgrade<br>[SECURITY-3870](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3870) fixed by the upgrade |
| pipeline-input-step | 2.11 | 560.v56198a_642157 | [SECURITY-2705](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2705) fixed by the upgrade<br>[SECURITY-2880](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2880) fixed by the upgrade |
| pipeline-maven | 3.8.3 | 1760.v9a_a_e6dcb_0444 | [SECURITY-1976](https://www.jenkins.io/security/advisory/2020-09-16/#SECURITY-1976) fixed by the upgrade<br>[SECURITY-3257](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3257) fixed by the upgrade |
| pipeline-model-definition | 1.3.8 | 2.2293.v6e7193cec599 | [SECURITY-3361](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3361) fixed by the upgrade |
| pipeline-stage-view | 2.19 | 2.41 | [SECURITY-2828](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2828) fixed by the upgrade |
| pipeline-utility-steps | 2.7.1 | 3.810.va_7672d206740 | [SECURITY-2948](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2948) fixed by the upgrade<br>[SECURITY-2949](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2949) fixed by the upgrade<br>[SECURITY-2196](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-2196) fixed by the upgrade |
| plain-credentials | 1.7 | 199.v9f8e1f741799 | [SECURITY-2495](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-2495) fixed by the upgrade |
| plot | 2.1.8 | 698.v42408994b_b_b_a_ | [SECURITY-2220](https://www.jenkins.io/security/advisory/2022-06-30/#SECURITY-2220) fixed by the upgrade<br>[SECURITY-2940](https://www.jenkins.io/security/advisory/2022-12-07/#SECURITY-2940) fixed by the upgrade |
| promoted-builds | 3.11 | 1029.v979cc7835e39 | [SECURITY-2617-promoted-builds](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-2655](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2655) fixed by the upgrade |
| publish-over-ftp | 1.16 | 261.v87f1561f44a_8 | [SECURITY-2321](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2321) fixed by the upgrade |
| publish-over-ssh | 1.20.1 | 409.v73ca_23a_5b_48e | [SECURITY-2287](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2287) fixed by the upgrade<br>[SECURITY-2290](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2290) fixed by the upgrade<br>[SECURITY-2291](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2291) fixed by the upgrade<br>[SECURITY-2307](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2307) fixed by the upgrade |
| rebuild | 1.31 | 344.v805899551b_13 | [SECURITY-3033](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3033) fixed by the upgrade |
| robot | 2.0.1 | 6.3.0 | [SECURITY-3971](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3971) fixed by the upgrade |
| role-strategy | 2.16 | 927.v9cf5527c4085 | [SECURITY-1767](https://www.jenkins.io/security/advisory/2020-10-08/#SECURITY-1767) fixed by the upgrade<br>[SECURITY-2182](https://www.jenkins.io/security/advisory/2021-03-18/#SECURITY-2182) fixed by the upgrade<br>[SECURITY-3053](https://www.jenkins.io/security/advisory/2023-03-21/#SECURITY-3053) fixed by the upgrade |
| saml | 1.1.2 | 4.623.v7875d61cd9f5 | [SECURITY-2469](https://www.jenkins.io/security/advisory/2021-08-31/#SECURITY-2469) fixed by the upgrade<br>[SECURITY-3613](https://www.jenkins.io/security/advisory/2025-10-29/#SECURITY-3613) fixed by the upgrade<br>[SECURITY-4065](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4065) fixed by the upgrade |
| script-security | 1.71 | 1436.v1082791135e4 | [SECURITY-1866](https://jenkins.io/security/advisory/2020-06-03/#SECURITY-1866) fixed by the upgrade<br>[SECURITY-2020](https://www.jenkins.io/security/advisory/2020-09-23/#SECURITY-2020) fixed by the upgrade<br>[SECURITY-2116](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2116) fixed by the upgrade<br>[SECURITY-2824-1-script-security](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-2564](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2564) fixed by the upgrade<br>[SECURITY-3016](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-3016) fixed by the upgrade<br>[SECURITY-3341](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3341) fixed by the upgrade<br>[SECURITY-3447](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3447) fixed by the upgrade<br>[SECURITY-3662](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3662) fixed by the upgrade<br>[SECURITY-3792](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3792) fixed by the upgrade<br>[SECURITY-3793](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3793) fixed by the upgrade<br>[SECURITY-3986](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3986) fixed by the upgrade<br>[SECURITY-3897](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3897) fixed by the upgrade<br>[SECURITY-3923](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3923) fixed by the upgrade<br>[SECURITY-3925-1](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(1)) fixed by the upgrade<br>[SECURITY-3925-2](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(2)) fixed by the upgrade<br>[SECURITY-3929](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3929) fixed by the upgrade<br>[SECURITY-3931](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3931) fixed by the upgrade<br>[SECURITY-3932](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3932) fixed by the upgrade<br>[SECURITY-3977](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3977) fixed by the upgrade |
| scriptler | 3.4 | 463.v0a_1f378da_06f | [SECURITY-3205](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3205) fixed by the upgrade<br>[SECURITY-3206](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3206) fixed by the upgrade |
| sidebar-link | 1.11.0 | 2.5.0 | [SECURITY-3125](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3125) fixed by the upgrade |
| sonar | 2.9 | 2.19.0 | [SECURITY-3989](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3989) fixed by the upgrade |
| sonar-quality-gates | 1.3.1 | 365.v76b_3ef6cf9b_e | [SECURITY-1523](https://jenkins.io/security/advisory/2020-03-09/#SECURITY-1523) fixed by the upgrade |
| ssh | 2.6.1 | 159.v496ca_25e5e82 | [SECURITY-2093](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2093) **no fix yet**<br>[SECURITY-2315](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2315) **no fix yet** |
| ssh-agent | 1.20 | 433.v9e73d67c3f74 | [SECURITY-2189](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2189) fixed by the upgrade |
| ssh2easy | 1.4 | 1.6 | [SECURITY-3064](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3064) fixed by the upgrade |
| structs | 1.21 | 362.va_b_695ef4fdf9 | [SECURITY-3371](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3371) fixed by the upgrade |
| subversion | 2.13.2 | 1327.vc36e0019a_653 | [SECURITY-2506](https://www.jenkins.io/security/advisory/2021-11-04/#SECURITY-2506) fixed by the upgrade<br>[SECURITY-2075](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2075) fixed by the upgrade<br>[SECURITY-2617-subversion](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade |
| swarm | 3.19 | 1258.vda_14ddda_0613 | [SECURITY-1200](https://jenkins.io/security/advisory/2020-06-03/#SECURITY-1200) fixed by the upgrade |
| testng-plugin | 1.15 | 981.v1dc64d227855 | [SECURITY-2788](https://www.jenkins.io/security/advisory/2022-06-30/#SECURITY-2788) fixed by the upgrade<br>[SECURITY-3047](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3047) fixed by the upgrade |
| thinBackup | 1.10 | 2.1.5 | [SECURITY-4099](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4099) fixed by the upgrade |
| timestamper | 1.9 | 1.30 | [SECURITY-1784](https://jenkins.io/security/advisory/2020-03-09/#SECURITY-1784) fixed by the upgrade |
| translation | 1.16 | 1.16 | [deprecated](https://github.com/jenkins-infra/update-center2/pull/719) |
| trilead-api | 1.0.1 | 2.284.v1974ea_324382 | [SECURITY-3333](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3333) fixed by the upgrade |
| uno-choice | 2.2.2 | 2.8.11 | [SECURITY-1954](https://www.jenkins.io/security/advisory/2020-10-08/#SECURITY-1954) fixed by the upgrade<br>[SECURITY-2008](https://www.jenkins.io/security/advisory/2020-10-08/#SECURITY-2008) fixed by the upgrade<br>[SECURITY-2192](https://www.jenkins.io/security/advisory/2021-02-24/#SECURITY-2192) fixed by the upgrade<br>[SECURITY-2219](https://www.jenkins.io/security/advisory/2021-11-12/#SECURITY-2219) fixed by the upgrade |
| warnings-ng | 8.1.0 | 13.10305.v3231d55c5755 | [SECURITY-2041](https://www.jenkins.io/security/advisory/2021-03-18/#SECURITY-2041) fixed by the upgrade<br>[SECURITY-2090](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2090) fixed by the upgrade<br>[SECURITY-3265](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3265) fixed by the upgrade<br>[SECURITY-3937](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3937) fixed by the upgrade |
| workflow-cps | 2.80 | 4395.v738456a_5da_6b_ | [SECURITY-2443](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2443) fixed by the upgrade<br>[SECURITY-2463-workflow-cps](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2463) fixed by the upgrade<br>[SECURITY-2613-workflow-cps](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2613) fixed by the upgrade<br>[SECURITY-359](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-359) fixed by the upgrade<br>[SECURITY-2824-1-workflow-cps](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-3362](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3362) fixed by the upgrade<br>[SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-cps-global-lib | 2.16 | 636.v3d2c06d018a_6 | [SECURITY-2422](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2422) fixed by the upgrade<br>[SECURITY-2441](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2441) fixed by the upgrade<br>[SECURITY-2463-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2463) fixed by the upgrade<br>[SECURITY-2586](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2586) fixed by the upgrade<br>[SECURITY-2613-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2613) fixed by the upgrade<br>[SECURITY-1951](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-1951) fixed by the upgrade<br>[SECURITY-2824-2-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(2)) fixed by the upgrade |
| workflow-job | 2.39 | 1609.v8926e7962024 | [SECURITY-3042](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3042) fixed by the upgrade |
| workflow-multibranch | 2.21 | 842.v3a_b_59b_57b_e6e | [SECURITY-2463-workflow-multibranch](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2463) fixed by the upgrade<br>[SECURITY-2613-workflow-multibranch](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2613) fixed by the upgrade<br>[SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |
| workflow-support | 3.4 | 1015.v785e5a_b_b_8b_22 | [SECURITY-2881](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2881) fixed by the upgrade |
| xunit | 2.3.9 | 3.1.7 | [SECURITY-2549](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2549) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| Office-365-Connector | 4.13.2 | 5.3.1 |
| Parameterized-Remote-Trigger | 3.1.6.3 | 3.2.2 |
| ace-editor | 1.0.1 | 1.1 |
| active-directory | 2.16.1 | 2.971.v653d8b_6e5548 |
| allure-jenkins-plugin | 2.30.3 | 3.0.0 |
| amazon-ecr | 1.6 | 1.175.ve233e80e56fd |
| analysis-model-api | 8.0.1 | 14.18.0-1028.v9c938eeda_ef2 |
| ansible | 1.1 | 635.v34b_48e979c86 |
| ansicolor | 0.5.3 | 542.v03d235fee02d |
| ant | 1.9 | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 2.1 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.13-1.0 | 4.5.14-269.vfa_2321039a_83 |
| apache-httpcomponents-client-5-api | – | 5.6.4-207.vd5b_a_c15e4de5 (new, for Office-365-Connector, analysis-model-api, audit-trail, cloudbees-bitbucket-branch-source, docker-java-api, file-operations, http_request, saml, slack, sonar-quality-gates) |
| artifactory | 3.5.0 | 4.0.8 |
| asm-api | – | 9.11.1-226.vff833a_e4c58e (new, for analysis-model-api, cvs, jacoco, jnr-posix-api, json-path-api, oic-auth, scm-api) |
| audit-trail | 3.4 | 464.v78b_845ffc5d8 |
| authentication-tokens | 1.3 | 1.144.v5ff4a_5ec5c33 |
| authorize-project | 1.3.0 | 534.v2f208c45e11c |
| aws-credentials | 1.28.2 | 265.v8422a_f384cd9 |
| aws-java-sdk | 1.11.723 | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-api-gateway | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-autoscaling | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-cloudformation | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-cloudfront | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-cloudwatch | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-codebuild | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-codedeploy | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-ec2 | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-credentials, aws-java-sdk, ec2) |
| aws-java-sdk-ecr | – | 1.12.780-480.v4a_0819121a_9e (new, for amazon-ecr, aws-java-sdk) |
| aws-java-sdk-ecs | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-efs | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-elasticbeanstalk | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-elasticloadbalancingv2 | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-iam | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-kinesis | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-lambda | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-logs | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-minimal | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, aws-java-sdk-api-gateway, aws-java-sdk-autoscaling, aws-java-sdk-cloudformation, aws-java-sdk-cloudfront, aws-java-sdk-cloudwatch, aws-java-sdk-codebuild, aws-java-sdk-codedeploy, aws-java-sdk-ec2, aws-java-sdk-ecr, aws-java-sdk-ecs, aws-java-sdk-efs, aws-java-sdk-elasticbeanstalk, aws-java-sdk-elasticloadbalancingv2, aws-java-sdk-iam, aws-java-sdk-kinesis, aws-java-sdk-lambda, aws-java-sdk-logs, aws-java-sdk-organizations, aws-java-sdk-secretsmanager, aws-java-sdk-sns, aws-java-sdk-sqs, aws-java-sdk-ssm, hashicorp-vault-plugin, s3) |
| aws-java-sdk-organizations | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-secretsmanager | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-sns | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-sqs | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, aws-java-sdk-sns) |
| aws-java-sdk-ssm | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk2-apigateway | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-cloudformation | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-cloudfront | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-codedeploy | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-core | – | 2.42.33-70.vd69c0763fa_60 (new, for aws-java-sdk2-apigateway, aws-java-sdk2-cloudformation, aws-java-sdk2-cloudfront, aws-java-sdk2-codedeploy, aws-java-sdk2-ec2, aws-java-sdk2-ecr, aws-java-sdk2-elasticbeanstalk, aws-java-sdk2-elasticloadbalancingv2, aws-java-sdk2-iam, aws-java-sdk2-lambda, aws-java-sdk2-netty-nio-client, aws-java-sdk2-organizations, aws-java-sdk2-s3, aws-java-sdk2-sns, ec2, pipeline-aws, s3) |
| aws-java-sdk2-ec2 | – | 2.42.33-70.vd69c0763fa_60 (new, for aws-credentials, ec2, pipeline-aws) |
| aws-java-sdk2-ecr | – | 2.42.33-70.vd69c0763fa_60 (new, for amazon-ecr, pipeline-aws) |
| aws-java-sdk2-elasticbeanstalk | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-elasticloadbalancingv2 | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-iam | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-lambda | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-netty-nio-client | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws, s3) |
| aws-java-sdk2-organizations | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| aws-java-sdk2-s3 | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws, s3) |
| aws-java-sdk2-sns | – | 2.42.33-70.vd69c0763fa_60 (new, for pipeline-aws) |
| azure-ad | 1.2.1 | 720.v2a_e3d48ffeb_7 |
| azure-commons | 1.0.5 | 1.1.3 |
| azure-credentials | 2.0.2 | 426.v5627e6a_4c7a_7 |
| azure-sdk | – | 299.v3996878fe83a_ (new, for azure-ad, azure-credentials, microsoft-graph-api) |
| badge | 1.8 | 3.615.vb_ed967ca_b_e93 |
| basic-branch-build-strategies | 1.3.2 | 317.v85b_331d6cc42 |
| bitbucket | 1.1.13 | 322.v882c6a_b_82c52 |
| bitbucket-push-and-pull-request | 2.8.4 | 4.1.0 |
| blueocean | 1.22.0 | 1.27.26 |
| blueocean-autofavorite | 1.2.4 | 1.2.5 |
| blueocean-bitbucket-pipeline | 1.22.0 | 1.27.26 |
| blueocean-commons | 1.22.0 | 1.27.26 |
| blueocean-config | 1.22.0 | 1.27.26 |
| blueocean-core-js | 1.22.0 | 1.27.26 |
| blueocean-dashboard | 1.22.0 | 1.27.26 |
| blueocean-display-url | 2.3.1 | 2.4.4 |
| blueocean-events | 1.22.0 | 1.27.26 |
| blueocean-git-pipeline | 1.22.0 | 1.27.26 |
| blueocean-github-pipeline | 1.22.0 | 1.27.26 |
| blueocean-i18n | 1.22.0 | 1.27.26 |
| blueocean-jira | 1.22.0 | 1.27.26 |
| blueocean-jwt | 1.22.0 | 1.27.26 |
| blueocean-personalization | 1.22.0 | 1.27.26 |
| blueocean-pipeline-api-impl | 1.22.0 | 1.27.26 |
| blueocean-pipeline-editor | 1.22.0 | 1.27.26 |
| blueocean-pipeline-scm-api | 1.22.0 | 1.27.26 |
| blueocean-rest | 1.22.0 | 1.27.26 |
| blueocean-rest-impl | 1.22.0 | 1.27.26 |
| blueocean-web | 1.22.0 | 1.27.26 |
| bootstrap4-api | 4.4.1-12 | 4.6.0-6 |
| bootstrap5-api | – | 5.3.8-1048.va_c299057e35c (new, for code-coverage-api, coverage, data-tables-api, echarts-api, forensics-api, junit, prism-api, warnings-ng) |
| bouncycastle-api | 2.18 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.5.6 | 2.1303.v9f3b_95dc329d |
| build-blocker-plugin | 1.7.3 | 177.ve433d85b_3fd7 |
| build-failure-analyzer | 1.26.0 | 3.889.vda_d3b_30435fc |
| build-monitor-plugin | 1.12+build.201809061734 | 1131.va_2d8e68fc912 |
| build-name-setter | 2.1.0 | 2.5.1 |
| build-pipeline-plugin | 1.5.8 | 2.2.0 |
| build-timeout | 1.19.1 | 1.41 |
| build-timestamp | 1.0.3 | 1.1.1 |
| build-token-root | 1.6 | 365.v717f8685a_09e |
| build-user-vars-plugin | 1.5 | 214.va_eed2ed849ca_ |
| build-with-parameters | 1.6 | 81.ve4a_9c2499d9a |
| built-on-column | 1.1 | 1.5 |
| caffeine-api | – | 3.3.1-214.v18eb_56b_34920 (new, for active-directory, azure-ad, azure-credentials, blueocean-events, blueocean-pipeline-api-impl, blueocean-rest-impl, branch-api, configuration-as-code, gerrit-trigger, github, github-branch-source, github-oauth, gitlab-oauth, gitlab-plugin, jira, kubernetes, lockable-resources, role-strategy, saml, script-security, slack, sonar, workflow-support) |
| checks-api | – | 422.vc7e8b_51488c6 (new, for code-coverage-api, coverage, junit, warnings-ng) |
| cloud-stats | 0.27 | 423.v63f017a_f582c |
| cloudbees-bitbucket-branch-source | 2.9.1 | 937.3.10 |
| cloudbees-disk-usage-simple | 0.9 | 275.v56e74c9dd0d3 |
| cloudbees-folder | 6.12 | 6.1106.v3a_d9a_6d2465e |
| cobertura | 1.16 | 1.17 |
| code-coverage-api | 1.1.6 | 4.99.0 |
| command-launcher | 1.4 | 141.v9578cf36b_311 |
| commons-collections4-api | – | 4.5.0-8.va_d5448ef9011 (new, for cloudbees-bitbucket-branch-source, configurationslicing, cucumber-reports, jira) |
| commons-compress-api | – | 1.28.0-96.v8da_8cfa_57e2b_ (new, for docker-java-api, pipeline-utility-steps) |
| commons-httpclient3-api | – | 3.1-3 (new, for Office-365-Connector, nexus-artifact-uploader) |
| commons-lang3-api | – | 3.21.0-122.v6a_ea_4fd34049 (new, for Office-365-Connector, allure-jenkins-plugin, amazon-ecr, analysis-model-api, ansible, ant, audit-trail, authentication-tokens, azure-ad, azure-credentials, basic-branch-build-strategies, bitbucket-push-and-pull-request, blueocean-commons, bootstrap5-api, branch-api, build-blocker-plugin, build-failure-analyzer, build-name-setter, build-user-vars-plugin, checks-api, cloudbees-bitbucket-branch-source, code-coverage-api, commons-compress-api, commons-text-api, configuration-as-code, configurationslicing, copyartifact, coverage, credentials, cucumber-reports, data-tables-api, dependency-check-jenkins-plugin, deploy, dingding-notifications, docker-build-step, docker-java-api, docker-plugin, echarts-api, email-ext, extended-choice-parameter, font-awesome-api, forensics-api, git, git-client, git-parameter, gitea, github, github-api, github-branch-source, github-oauth, github-pullrequest, gitlab-oauth, gitlab-plugin, gradle, hashicorp-vault-plugin, htmlpublisher, ionicons-api, ivy, jira, job-dsl, job-import-plugin, mask-passwords, maven-plugin, mercurial, metrics, nexus-artifact-uploader, nodejs, nodelabelparameter, oic-auth, p4, pipeline-graph-analysis, pipeline-maven, pipeline-maven-api, pipeline-model-api, pipeline-model-definition, pipeline-model-extensions, pipeline-npm, pipeline-rest-api, pipeline-utility-steps, plot, plugin-util-api, popper2-api, powershell, prism-api, prometheus, publish-over, publish-over-ftp, publish-over-ssh, robot, role-strategy, saml, scm-api, slack, sonar, ssh-slaves, stashNotifier, subversion, swarm, testng-plugin, throttle-concurrents, timestamper, token-macro, uno-choice, warnings-ng, workflow-basic-steps, xunit) |
| commons-text-api | – | 1.15.0-218.va_61573470393 (new, for analysis-model-api, ansicolor, blueocean-commons, bootstrap5-api, build-failure-analyzer, checks-api, code-coverage-api, configuration-as-code, coverage, credentials, cucumber-reports, data-tables-api, dependency-check-jenkins-plugin, echarts-api, email-ext, embeddable-build-status, extended-choice-parameter, font-awesome-api, forensics-api, github-pullrequest, groovy-postbuild, hashicorp-vault-plugin, job-dsl, job-import-plugin, log-parser, oic-auth, p4, parameter-separator, pipeline-model-api, pipeline-utility-steps, plot, plugin-util-api, popper2-api, prism-api, publish-over-ssh, saml, stashNotifier, token-macro, uno-choice, warnings-ng) |
| conditional-buildstep | 1.3.6 | 1.5.0 |
| config-file-provider | 3.7.2 | 1013.v73c323e52b_1f |
| configuration-as-code | 1.36.2 | 2136.veb_3a_554550b_f |
| configuration-as-code-groovy | 1.1 | 1.1 |
| configurationslicing | 1.47 | 715.vcc64f669003e |
| copyartifact | 1.44 | 805.v9d1393350780 |
| coverage | – | 3.3497.v3fee733ca_b_35 (new, for code-coverage-api, pipeline-maven, prometheus) |
| credentials | 2.3.11 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 1.18 | 728.v902a_273b_8947 |
| cucumber-reports | 5.4.0 | 5.11.0 |
| cvs | 2.17 | 502.v02d5a_0c6b_9a_3 |
| dashboard-view | 2.12.1 | 2.569.v5ddb_c9b_a_7264 |
| data-tables-api | 1.10.20-14 | 2.3.8-1570.v1cb_1cd2a_0fb_c |
| delivery-pipeline-plugin | 1.3.2 | 1125.vb_541b_47dfd78 |
| dependency-check-jenkins-plugin | 5.1.2 | 5.6.5 |
| deploy | 1.16 | 329.vd9ec45431fb_e |
| description-setter | 1.10 | 264.v1957f215dcd5 |
| dingding-notifications | 1.9 | 2.9.1 |
| discard-old-build | 1.05 | 1.08 |
| disk-usage | 0.28 | 573.v5269a_8d16a_db_ |
| display-url-api | 2.3.2 | 2.217.va_6b_de84cc74b_ |
| docker-build-publish | 1.3.2 | 1.4.0 |
| docker-build-step | 2.8 | 2.14 |
| docker-commons | 1.14 | 477.v289085a_b_6896 |
| docker-java-api | 3.2.13-37.vf3411c9828b9 | 3.7.1-136.v5d70f77a_c3d6 |
| docker-plugin | 1.2.2 | 1336.v3711dec712de |
| docker-workflow | 1.21 | 653.v2f2c08eff0ec |
| dtkit-api | 2.1.2 | 3.0.3 |
| durable-task | 1.30 | 686.v80ff80875b_82 |
| ec2 | 1.44.1 | 2064.v79de4fd2d03b_ |
| echarts-api | 4.6.0-8 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | – | 0.3.0.1-29.v67e9a_1c969b_b_ (new, for mina-sshd-api-common, p4, subversion, trilead-api) |
| email-ext | 2.69.2 | 2038.v7b_8817a_499d9 |
| emailext-template | 1.1 | 249.v81d6ea_f33f2e |
| embeddable-build-status | 2.0.4 | 666.vcb_0de9eb_43d8 |
| emoji-symbols-api | – | 17.0-57.v8d44b_9a_b_d5ea_ (new, for badge) |
| envinject | 2.3.0 | 2.941.v351a_20c0a_3ca_ |
| envinject-api | 1.8 | 1.241.vdd714803b_403 |
| extended-choice-parameter | 346.vd87693c5a_86c | 388.ve7b_d0b_920e10 |
| extended-read-permission | 3.2 | 92.v1b_f0c20b_75db_ |
| extensible-choice-parameter | 1.8.0 | 255.vfa_41f46828ef |
| external-monitor-job | 1.7.1 | 223.vb_fddcf42c9b_3 |
| extra-columns | 1.23 | 1.29 |
| favorite | 2.3.2 | 2.267.vb_90d08408081 |
| file-operations | 1.11 | 439.vdc9d3b_d74b_d9 |
| flatpickr-api | – | 4.6.13-32.v60a_51029c136 (new, for schedule-build) |
| flexible-publish | 0.15.2 | 0.17 |
| folder-properties | 1.2.1 | 97.v0c61a_c452eb_b_ |
| font-awesome-api | 5.13.0-1 | 7.3.1-1013.v0835a_879ec6d |
| forensics-api | 0.7.0 | 4.1923.v445c56df9b_2e |
| generic-webhook-trigger | 1.76 | 2.4.3 |
| gerrit-trigger | 2.34.0 | 3.1998.v9487a_a_89a_eb_1 |
| ghprb | 1.42.2 | 1.42.2 |
| git | 4.2.2 | 5.11.0 |
| git-changelog | 3.31 | 3.45 |
| git-client | 3.2.1 | 6.6.1 |
| git-parameter | 0.9.13 | 462.463.v496a_59f698e5 |
| git-server | 1.9 | 137.ve0060b_432302 |
| gitea | 1.3.0 | 286.v0740a_413c1eb_ |
| github | 1.30.0 | 1.47.0 |
| github-api | 1.95 | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 2.5.8 | 1986.v801e40614d00 |
| github-oauth | 0.33 | 685.v53b_070455063 |
| github-organization-folder | 1.6 | 1.6 |
| github-pullrequest | 0.2.8 | 0.7.5 |
| gitlab-api | 1.0.6 | 6.3.0-126.vd066b_d896159 |
| gitlab-merge-request-jenkins | 2.0.0 | 2.0.0 |
| gitlab-oauth | 1.11 | 1.30 |
| gitlab-plugin | 1.5.13 | 1.2158.v49db_f64b_fc28 |
| global-build-stats | 1.5 | 360.vd08338d06c53 |
| global-slack-notifier | 1.5 | 1.5 |
| google-oauth-plugin | 1.0.3 | 1.346.v712b_a_e0d1e60 |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ |
| greenballs | 1.15 | 1.15.1 |
| groovy | 2.4 | 537.v741a_5a_f1b_581 |
| groovy-postbuild | 2.5 | 335.vf10eb_c1b_5e14 |
| gson-api | – | 2.14.0-201.v8eefe5515533 (new, for Office-365-Connector, analysis-model-api, azure-ad, bitbucket-push-and-pull-request, dingding-notifications, generic-webhook-trigger, gerrit-trigger, git-client, google-oauth-plugin, microsoft-graph-api, p4, performance, sonar-quality-gates, trilead-api, warnings-ng) |
| h2-api | 1.4.199 | 11.1.4.199-36.vb_ee07e965744 |
| handlebars | 1.1.1 | 1.1.1 |
| handy-uri-templates-2-api | 2.1.8-1.0 | 2.1.8-38.vcea_5d521d5f3 |
| hashicorp-vault-plugin | 3.4.0 | 384.vda_86ec66c537 |
| hidden-parameter | 0.0.4 | 546.vdeec4ce8b_26f |
| htmlpublisher | 1.21 | 429 |
| http_request | 1.8.26 | 1.659.v1b_9d9e942a_e6 |
| icon-shim | 2.0.3 | 3.0.0 |
| instance-identity | – | 203.v15e81a_1b_7a_38 (new, for docker-plugin, git-server, github, mailer, pubsub-light, sshd) |
| ionicons-api | – | 94.vcc3065403257 (new, for active-directory, badge, branch-api, build-failure-analyzer, cloudbees-disk-usage-simple, cloudbees-folder, code-coverage-api, config-file-provider, coverage, credentials, dashboard-view, docker-plugin, emailext-template, embeddable-build-status, favorite, gerrit-trigger, git, github-branch-source, global-build-stats, jobConfigHistory, junit, matrix-auth, matrix-project, metrics, naginator, nested-view, next-build-number, pipeline-model-definition, role-strategy, schedule-build, scriptler, test-results-analyzer, thinBackup, workflow-cps, workflow-job, workflow-multibranch) |
| ivy | 1.28 | 609.v66779e68503d |
| jackson-annotations2-api | – | 2.22-19.v10a_a_582ea_26e (new, for jackson2-api, jackson3-api) |
| jackson2-api | 2.10.3 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | – | 3.2.4-103.vf5a_b_3f4ea_cec (new, for coverage, data-tables-api, echarts-api, file-operations, forensics-api, junit, pipeline-utility-steps) |
| jacoco | 3.0.4 | 3.3.7 |
| jakarta-activation-api | – | 2.1.4-1 (new, for configurationslicing, gitlab-api, jackson2-api, jakarta-mail-api, jakarta-xml-bind-api, jersey3-api, jobConfigHistory) |
| jakarta-mail-api | – | 2.1.5-1 (new, for configurationslicing, disk-usage, email-ext, mailer, p4, workflow-basic-steps) |
| jakarta-xml-bind-api | – | 4.0.9-19.v2b_a_5b_44d9a_1c (new, for gitlab-oauth, gitlab-plugin, jackson2-api, jersey3-api, jobConfigHistory) |
| javadoc | 1.6 | 354.vee1a_660b_4990 |
| javax-activation-api | – | 1.2.0-8 (new, for gitlab-api, gitlab-plugin, jackson2-api, jaxb, jersey2-api) |
| jaxb | 2.3.0 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 1.4 | 83.v417146707a_3d |
| jenkins-design-language | 1.22.0 | 1.27.26 |
| jenkins-multijob-plugin | 1.36 | 682.v4ed91b_52d621 |
| jersey2-api | – | 2.48-180.ve47b_264f849b_ (new, for artifactory, gitlab-api, gitlab-plugin, jira) |
| jersey3-api | – | 3.1.11-4.v77818819c2e1 (new, for gitlab-api, gitlab-plugin) |
| jira | 3.0.13 | 3.22 |
| jira-steps | 1.5.5 | 2.1.197.vdd1b_3987e040 |
| jjwt-api | – | 0.13.0-141.vd58b_a_9592b_6c (new, for github-api, github-branch-source) |
| jnr-posix-api | 3.0.45 | 3.2.4-223.v4f0cfd7fa_006 |
| job-dsl | 1.76 | 3732.v9a_c49a_61a_313 |
| job-import-plugin | 3.6 | 143.145.v48f9a_a_6ff384 |
| job-restrictions | 0.8 | 242.v6edda_c9e4ca_f |
| jobConfigHistory | 2.25 | 1380.v762185b_9a_793 |
| joda-time-api | – | 2.15.0-208.v6b_084e7a_d9f0 (new, for aws-java-sdk-minimal, cucumber-reports, delivery-pipeline-plugin, gitlab-plugin, handy-uri-templates-2-api, jira, pipeline-model-api) |
| jquery | 1.12.4-1 | 1.12.4-3 |
| jquery-detached | 1.2.1 | 1.2.1 |
| jquery-ui | 1.0.2 | 1.0.2 |
| jquery3-api | 3.4.1-10 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.1.55.1 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | – | 20260814-226.v20f9685d642c (new, for analysis-model-api, configuration-as-code, jackson2-api, jira, p4, sonar-quality-gates, warnings-ng) |
| json-path-api | – | 3.0.0-218.vcd4dd1355de2 (new, for analysis-model-api, azure-sdk, build-timeout, generic-webhook-trigger, oic-auth, token-macro, warnings-ng) |
| jsoup | 1.10.3 | 1.23.2-106.v1161f282b_eb_f |
| junit | 1.28 | 1442.vb_b_a_467df2767 |
| kubernetes | 1.16.7 | 4560.vc25b_83890c34 |
| kubernetes-cli | 1.7.0 | 1.375.v8001982544ba |
| kubernetes-client-api | – | 7.3.1-256.v788a_0b_787114 (new, for kubernetes, kubernetes-credentials, kubernetes-credentials-provider) |
| kubernetes-credentials | 0.4.1 | 207.v492f58828b_ed |
| kubernetes-credentials-provider | 0.12.1 | 1.315.v92008589c044 |
| ldap | 1.25 | 825.v2fca_37dd5b_cb_ |
| locale | 1.4 | 641.v84f22d75fd8b_ |
| localization-support | – | 1.41.v186a_c1569458 (new, for localization-zh-cn) |
| localization-zh-cn | 0.0.15 | 371.v23851f835d6b_ |
| lockable-resources | 2.10 | 1576.v217259b_f9976 |
| log-parser | 2.2 | 3.0.7 |
| login-theme | 1.1 | 339.v8f0b_a_5695303 |
| m2release | 0.16.3 | 547.v78681a_45411a_ |
| mailer | 1.23 | 534.v1b_36f5864073 |
| mapdb-api | 1.0.9.0 | 1.0.9-44.va_1e1310c9118 |
| mask-passwords | 2.13 | 220.v95819055b_265 |
| matrix-auth | 2.5.2 | 3.3 |
| matrix-project | 1.14.1 | 928.v1effc3f3a_236 |
| maven-plugin | 3.6.1 | 3.28 |
| mercurial | 2.8.1 | 1323.ve69d2a_db_8a_b_d |
| metrics | 4.0.2.6 | 4.2.37-494.v06f9a_939d33a_ |
| microsoft-graph-api | – | 1.225181-2.v4b_0434d7d9c5 (new, for azure-ad) |
| mina-sshd-api-common | – | 2.19.0-192.v2b_a_7b_2c1dc71 (new, for git-client, mina-sshd-api-core, sshd, subversion) |
| mina-sshd-api-core | – | 2.19.0-192.v2b_a_7b_2c1dc71 (new, for ec2, git-client, mina-sshd-api-scp, sshd, subversion) |
| mina-sshd-api-scp | – | 2.19.0-192.v2b_a_7b_2c1dc71 (new, for ec2) |
| momentjs | 1.1 | 1.1.1 |
| monitoring | 1.95.0 | 2.9.1 |
| msbuild | 1.29 | 1.45 |
| multibranch-scan-webhook-trigger | 1.0.5 | 1.0.11 |
| multiple-scms | 0.6 | 0.8 |
| naginator | 1.18.2 | 1.556.v14d723a_109a_c |
| nested-view | 1.17 | 241.v178f0b_a_cd76a_ |
| netty-api | – | 4.2.18.Final-34.v861b_21a_78c0f (new, for aws-java-sdk2-netty-nio-client, azure-sdk) |
| next-build-number | 1.6 | 66.v4b_4762172d53 |
| nexus-artifact-uploader | 2.13 | 2.16 |
| nimbus-jose-jwt-api | – | 10.10-12.vd4e6d651b_7f7 (new, for oic-auth) |
| node-iterator-api | 1.5.1 | 72.vc90e81737df1 |
| nodejs | 1.3.5 | 1.6.6 |
| nodelabelparameter | 1.7.2 | 851.vd94e5048d321 |
| oauth-credentials | 0.4 | 0.668.vb_431feb_d8a_4e |
| oic-auth | 1.8 | 4.731.v4932b_506e767 |
| okhttp-api | – | 5.5.0-242.v1482d48d3ea_3 (new, for azure-ad, git-changelog, github, github-api, github-branch-source, gradle, kubernetes-client-api, microsoft-graph-api, opentelemetry-api, sonar) |
| opentelemetry-api | – | 1.54.1.102.v276c93ee966d (new, for microsoft-graph-api) |
| oss-symbols-api | – | 477.v7197b_eb_3d8b_8 (new, for javadoc, maven-plugin, pipeline-maven-api) |
| p4 | 1.13.3 | 1.18.1 |
| pam-auth | 1.6 | 1.12 |
| parameter-separator | 1.0 | 334.v1fc0e534c71d |
| parameterized-scheduler | 0.8 | 379.v95c73f233a_df |
| parameterized-trigger | 2.36 | 893.va_383a_9b_a_4c11 |
| performance | 3.20 | 1022.v48fce959f0c3 |
| permissive-script-security | 0.7 | 0.7 |
| pipeline-aws | 1.43 | 767.v48d5c10b_1899 |
| pipeline-build-step | 2.12 | 601.v6d4c6d1a_9dc7 |
| pipeline-github | 2.5 | 2.8-191.cb_1a_07f10b_a_f |
| pipeline-github-lib | 1.0 | 65.v203688e7727e |
| pipeline-githubnotify-step | 1.0.5 | 49.vf37bf92d2bc8 |
| pipeline-graph-analysis | 1.10 | 254.v0f63a_a_447dca_ |
| pipeline-groovy-lib | – | 822.vc95f70dea_db_7 (new, for groovy, pipeline-github-lib, pipeline-model-definition, subversion, workflow-aggregator, workflow-cps-global-lib) |
| pipeline-input-step | 2.11 | 560.v56198a_642157 |
| pipeline-maven | 3.8.3 | 1760.v9a_a_e6dcb_0444 |
| pipeline-maven-api | – | 1760.v9a_a_e6dcb_0444 (new, for pipeline-maven) |
| pipeline-milestone-step | 1.3.2 | 152.v6e22b_8cfc66c |
| pipeline-model-api | 1.3.8 | 2.2293.v6e7193cec599 |
| pipeline-model-declarative-agent | 1.1.1 | 1.1.1 |
| pipeline-model-definition | 1.3.8 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 1.3.8 | 2.2293.v6e7193cec599 |
| pipeline-multibranch-defaults | 2.1 | 2.1 |
| pipeline-npm | 0.9.2 | 503.vf29ef1b_913ce |
| pipeline-rest-api | 2.19 | 2.41 |
| pipeline-stage-step | 2.4 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 1.3.8 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.19 | 2.41 |
| pipeline-utility-steps | 2.7.1 | 3.810.va_7672d206740 |
| plain-credentials | 1.7 | 199.v9f8e1f741799 |
| plot | 2.1.8 | 698.v42408994b_b_b_a_ |
| plugin-usage-plugin | 2.2 | 418.v308c3c863d25 |
| plugin-util-api | 1.1.0 | 7.1403.vd380b_c42e90c |
| popper-api | 1.16.0-6 | 1.16.1-3 |
| popper2-api | – | 2.11.6-5 (new, for bootstrap5-api) |
| postbuild-task | 1.8 | 78.v24529f1f5cdb_ |
| postbuildscript | 2.9.1 | 3.4.1.821.v08352159216f |
| powershell | 1.8 | 185.v7a_026da_c54ee |
| preSCMbuildstep | 0.3 | 115.vf9808414429d |
| prism-api | – | 1.30.0-741.v034eb_0b_0a_a_fa_ (new, for code-coverage-api, configuration-as-code, coverage, junit, warnings-ng) |
| prometheus | 2.0.6 | 860.v532442b_44e9a_ |
| promoted-builds | 3.11 | 1029.v979cc7835e39 |
| publish-over | 0.22 | 242.vf56f6892cb_b_a_ |
| publish-over-ftp | 1.16 | 261.v87f1561f44a_8 |
| publish-over-ssh | 1.20.1 | 409.v73ca_23a_5b_48e |
| pubsub-light | 1.13 | 1.19 |
| python | 1.3 | 1.3 |
| rebuild | 1.31 | 344.v805899551b_13 |
| resource-disposer | 0.14 | 0.25 |
| robot | 2.0.1 | 6.3.0 |
| role-strategy | 2.16 | 927.v9cf5527c4085 |
| run-condition | 1.3 | 276.v97298f3a_cd51 |
| s3 | 0.11.8.1 | 522.vc73d01440c29 |
| saferestart | 0.3 | 125.vdb_73765c38a_e |
| saml | 1.1.2 | 4.623.v7875d61cd9f5 |
| schedule-build | 0.5.1 | 696.v0a_b_ec7334c89 |
| scm-api | 2.6.3 | 728.vc30dcf7a_0df5 |
| script-security | 1.71 | 1436.v1082791135e4 |
| scriptler | 3.4 | 463.v0a_1f378da_06f |
| sidebar-link | 1.11.0 | 2.5.0 |
| simple-theme-plugin | 0.6 | 230.v8b_fd91b_b_800c |
| skip-certificate-check | 1.0 | 50.vcfa_cdb_b_27312 |
| slack | 2.29 | 795.v4b_9705b_e6d47 |
| snakeyaml-api | – | 2.7-152.vc672728cf536 (new, for configuration-as-code, jackson2-api, kubernetes-client-api, pipeline-aws, pipeline-utility-steps, robot) |
| snakeyaml-engine-api | – | 3.2-18.v885ed7ff698f (new, for jackson3-api) |
| sonar | 2.9 | 2.19.0 |
| sonar-quality-gates | 1.3.1 | 365.v76b_3ef6cf9b_e |
| sse-gateway | 1.23 | 1.29 |
| ssh | 2.6.1 | 159.v496ca_25e5e82 |
| ssh-agent | 1.20 | 433.v9e73d67c3f74 |
| ssh-credentials | 1.17.4 | 372.va_250881b_08cd |
| ssh-slaves | 1.29.4 | 3.1097.v868116049892 |
| ssh-steps | 2.0.0 | 2.0.92.vb_a_0583935f9b_ |
| ssh2easy | 1.4 | 1.6 |
| sshd | – | 3.384.vc89b_5e138cf9 (new, for git-server, scriptler) |
| startup-trigger-plugin | 2.9.3 | 2.9.4 |
| stashNotifier | 1.24 | 1.557.vb_d02a_d29f87f |
| strict-crumb-issuer | 2.1.0 | 2.58.va_0da_723a_dc84 |
| structs | 1.21 | 362.va_b_695ef4fdf9 |
| subversion | 2.13.2 | 1327.vc36e0019a_653 |
| swarm | 3.19 | 1258.vda_14ddda_0613 |
| test-results-analyzer | 0.3.5 | 353.vea_f73668d8b_5 |
| testng-plugin | 1.15 | 981.v1dc64d227855 |
| text-finder | 1.12 | 1.34 |
| thinBackup | 1.10 | 2.1.5 |
| throttle-concurrents | 2.0.2 | 625.vc8b_e469e9a_b_c |
| timestamper | 1.9 | 1.30 |
| token-macro | 2.12 | 477.vd4f0dc3cb_cf1 |
| translation | 1.16 | 1.16 |
| trilead-api | 1.0.1 | 2.284.v1974ea_324382 |
| uno-choice | 2.2.2 | 2.8.11 |
| validating-string-parameter | 2.5 | 345.v0f4330488a_f6 |
| variant | 1.3 | 70.va_d9f17f859e0 |
| versionnumber | 1.9 | 234.v315d3b_3cb_fb_5 |
| view-job-filters | 2.1.1 | 408.vb_ca_5e798e5b_5 |
| warnings-ng | 8.1.0 | 13.10305.v3231d55c5755 |
| woodstox-core-api | – | 7.3.0-16.v017f30ed004c (new, for jackson2-api, jackson3-api) |
| workflow-aggregator | 2.7 | 608.v67378e9d3db_1 |
| workflow-api | 2.39 | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 2.16.1 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 2.80 | 4395.v738456a_5da_6b_ |
| workflow-cps-global-lib | 2.16 | 636.v3d2c06d018a_6 |
| workflow-durable-task-step | 2.28 | 1479.v56e587f413a_7 |
| workflow-job | 2.39 | 1609.v8926e7962024 |
| workflow-multibranch | 2.21 | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 2.11 | 487.v41313a_96b_65d |
| workflow-step-api | 2.22 | 724.v538c2362b_dfb_ |
| workflow-support | 3.4 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.38 | 0.49 |
| xunit | 2.3.9 | 3.1.7 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
