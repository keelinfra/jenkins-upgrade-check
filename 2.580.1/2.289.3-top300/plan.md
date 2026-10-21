# Jenkins upgrade plan: 2.289.3 → 2.580.1

Source: results/2.289.3-top300/plugins-2.289.3.txt · generated 2026-10-21 12:43 UTC by jenkins-upgrade-check dev

**3 core upgrade(s)**, 351 plugins: 286 to upgrade, 49 new dependencies; 72 with security warnings today, 3 still with warnings after the upgrade, 3 deprecated, 2 not from the update center.

Jenkins 2.289.3 runs on Java 8, 11; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

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

1. **On 2.289.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 11**, upgrade Jenkins to **2.462.3**, start it, then update plugins (306 change(s)).
3. **Run the controller and every agent on Java 17**, upgrade Jenkins to **2.541.3**, start it, then update plugins (326 change(s)).
4. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (27 change(s)).

No Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| Parameterized-Remote-Trigger | 3.1.6.3 | 3.2.2 | [SECURITY-3768](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3768) **no fix yet**<br>[SECURITY-3800](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3800) **no fix yet** |
| ace-editor | 1.1 | 1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| active-directory | 2.27 | 2.971.v653d8b_6e5548 | [SECURITY-3059](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3059) fixed by the upgrade<br>[SECURITY-3659](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3659) fixed by the upgrade<br>[SECURITY-3651](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3651) fixed by the upgrade |
| allure-jenkins-plugin | 2.30.3 | 3.0.0 | [SECURITY-3645](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3645) fixed by the upgrade |
| analysis-model-api | 10.10.1 | 14.18.0-1028.v9c938eeda_ef2 | [SECURITY-3327](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3327) fixed by the upgrade |
| ansible | 1.1 | 635.v34b_48e979c86 | [SECURITY-3017](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3017) fixed by the upgrade |
| atlassian-bitbucket-server-integration | 3.3.1 | 6.0.2 | [SECURITY-3434](https://www.jenkins.io/security/advisory/2025-01-22/#SECURITY-3434) fixed by the upgrade<br>[SECURITY-3778](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3778) fixed by the upgrade |
| authorize-project | 1.4.0 | 534.v2f208c45e11c | [SECURITY-3010](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3010) fixed by the upgrade |
| blueocean | 1.26.0 | 1.27.26 | [SECURITY-3116](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3116) fixed by the upgrade |
| build-monitor-plugin | 1.13+build.202202112334 | 1131.va_2d8e68fc912 | [SECURITY-3280](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3280) fixed by the upgrade |
| build-pipeline-plugin | 1.5.8 | 2.2.0 | [SECURITY-879](https://jenkins.io/security/advisory/2019-08-07/#SECURITY-879) fixed by the upgrade |
| cloudbees-bitbucket-branch-source | 762.v969cfe087fc0 | 937.3.10 | [SECURITY-3300](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3300) fixed by the upgrade<br>[SECURITY-3363](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3363) fixed by the upgrade |
| cloudbees-folder | 6.722.v8165b_a_cf25e9 | 6.1106.v3a_d9a_6d2465e | [SECURITY-3105](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3105) fixed by the upgrade<br>[SECURITY-3106](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3106) fixed by the upgrade<br>[SECURITY-3109](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3109) fixed by the upgrade |
| config-file-provider | 3.8.1 | 1013.v73c323e52b_1f | [SECURITY-3090](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3090) fixed by the upgrade |
| credentials | 2.6.1.1 | 1514.vb_01eb_014e7d4 | [SECURITY-3373](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3373) fixed by the upgrade |
| credentials-binding | 1.27.1 | 728.v902a_273b_8947 | [SECURITY-3499](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3499) fixed by the upgrade<br>[SECURITY-3672](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3672) fixed by the upgrade<br>[SECURITY-3790](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3790) fixed by the upgrade |
| dependency-check-jenkins-plugin | 5.2.0 | 5.6.5 | [SECURITY-3344](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3344) fixed by the upgrade<br>[SECURITY-3992](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3992) fixed by the upgrade |
| docker-build-step | 2.8 | 2.14 | [SECURITY-3200](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3200) fixed by the upgrade |
| email-ext | 2.87 | 2038.v7b_8817a_499d9 | [SECURITY-2931](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2931) fixed by the upgrade<br>[SECURITY-2934](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2934) fixed by the upgrade<br>[SECURITY-2939](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2939) fixed by the upgrade<br>[SECURITY-3088-1](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(1)) fixed by the upgrade<br>[SECURITY-3088-2](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(2)) fixed by the upgrade<br>[SECURITY-3705](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3705) fixed by the upgrade |
| extended-choice-parameter | 346.vd87693c5a_86c | 388.ve7b_d0b_920e10 | [SECURITY-1350](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1350) **no fix yet**<br>[SECURITY-1351](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1351) **no fix yet**<br>[SECURITY-2232](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2232) **no fix yet**<br>[SECURITY-2617-extended-choice-parameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) **no fix yet**<br>[deprecated](https://github.com/jenkinsci/extended-choice-parameter-plugin?tab=readme-ov-file#end-of-life) |
| extensible-choice-parameter | 1.8.0 | 255.vfa_41f46828ef | [SECURITY-3583](https://www.jenkins.io/security/advisory/2025-10-29/#SECURITY-3583) fixed by the upgrade |
| external-monitor-job | 203.v683c09d993b_9 | 223.vb_fddcf42c9b_3 | [SECURITY-3133](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3133) fixed by the upgrade |
| file-parameters | 205.vf6ce13b_e5dee | 433.va_0b_80359d54d | [SECURITY-3123](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3123) fixed by the upgrade<br>[SECURITY-4093](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4093) fixed by the upgrade |
| git-client | 3.11.2 | 6.6.1 | [SECURITY-3590](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3590) fixed by the upgrade<br>[SECURITY-3614](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3614) fixed by the upgrade<br>[SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| git-parameter | 0.9.14 | 462.463.v496a_59f698e5 | [SECURITY-2617-git-parameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) fixed by the upgrade<br>[SECURITY-3419](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3419) fixed by the upgrade<br>[SECURITY-3745](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3745) fixed by the upgrade |
| git-server | 1.11 | 137.ve0060b_432302 | [SECURITY-3319](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3319) fixed by the upgrade<br>[SECURITY-3342](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3342) fixed by the upgrade |
| github | 1.34.3.1 | 1.47.0 | [SECURITY-3246](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3246) fixed by the upgrade<br>[SECURITY-3704](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3704) fixed by the upgrade |
| github-branch-source | 2.11.4 | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| github-pullrequest | 0.4.0 | 0.7.5 | [SECURITY-3776](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3776) fixed by the upgrade |
| gitlab-branch-source | 650.va_d1ce6d01959 | 744.vb_d0403d08ec7 | [SECURITY-2871](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-2871) fixed by the upgrade<br>[SECURITY-3040](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3040) fixed by the upgrade<br>[SECURITY-3251](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3251) fixed by the upgrade |
| gitlab-oauth | 1.16 | 1.30 | [SECURITY-2696](https://www.jenkins.io/security/advisory/2023-07-26/#SECURITY-2696) fixed by the upgrade |
| gitlab-plugin | 1.5.36 | 1.2158.v49db_f64b_fc28 | [SECURITY-3260](https://www.jenkins.io/security/advisory/2025-01-22/#SECURITY-3260) fixed by the upgrade<br>[SECURITY-4097](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4097) fixed by the upgrade<br>[SECURITY-3851](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3851) fixed by the upgrade |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| greenballs | 1.15.1 | 1.15.1 | [deprecated](https://github.com/jenkins-infra/update-center2/pull/735) |
| hashicorp-vault-plugin | 336.v182c0fbaaeb7 | 384.vda_86ec66c537 | [SECURITY-2521](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2521) fixed by the upgrade<br>[SECURITY-2593](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-2593) fixed by the upgrade<br>[SECURITY-3077](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3077) fixed by the upgrade<br>[SECURITY-3045](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3045) fixed by the upgrade |
| htmlpublisher | 1.31 | 429 | [SECURITY-3301](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3301) fixed by the upgrade<br>[SECURITY-3302](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3302) fixed by the upgrade<br>[SECURITY-3303](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3303) fixed by the upgrade<br>[SECURITY-3547](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3547) fixed by the upgrade<br>[SECURITY-3706](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3706) fixed by the upgrade |
| ivy | 2.2 | 609.v66779e68503d | [SECURITY-2924](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-2924) fixed by the upgrade<br>[SECURITY-3093](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3093) fixed by the upgrade |
| jakarta-mail-api | 2.0.0-5 | 2.1.5-1 | [SECURITY-3617](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3617) fixed by the upgrade |
| jenkins-multijob-plugin | 1.36 | 682.v4ed91b_52d621 | [SECURITY-3781](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3781) fixed by the upgrade<br>[SECURITY-3823-1](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3823%20(1)) fixed by the upgrade<br>[SECURITY-3823-2](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3823%20(2)) fixed by the upgrade |
| jira | 3.7.1 | 3.22 | [SECURITY-3225](https://www.jenkins.io/security/advisory/2023-11-29/#SECURITY-3225) fixed by the upgrade |
| job-import-plugin | 3.6 | 143.145.v48f9a_a_6ff384 | [SECURITY-3783](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3783) fixed by the upgrade |
| jobConfigHistory | 2.30 | 1380.v762185b_9a_793 | [SECURITY-2766](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-2766) fixed by the upgrade<br>[SECURITY-2765](https://www.jenkins.io/security/advisory/2022-08-23/#SECURITY-2765) fixed by the upgrade<br>[SECURITY-3233](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3233) fixed by the upgrade<br>[SECURITY-3235](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3235) fixed by the upgrade<br>[SECURITY-3742](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3742) fixed by the upgrade<br>[SECURITY-4098](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4098) fixed by the upgrade |
| junit | 1.53.0.1 | 1442.vb_b_a_467df2767 | [SECURITY-2888](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2888) fixed by the upgrade<br>[SECURITY-3032](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-3032) fixed by the upgrade |
| kubernetes | 1.30.11 | 4560.vc25b_83890c34 | [SECURITY-3075-kubernetes](https://www.jenkins.io/security/advisory/2023-04-12/#SECURITY-3075) fixed by the upgrade |
| kubernetes-credentials-provider | 0.20 | 1.315.v92008589c044 | [SECURITY-3022](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-3022) fixed by the upgrade |
| ldap | 2.9 | 825.v2fca_37dd5b_cb_ | [SECURITY-3046](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3046) fixed by the upgrade<br>[SECURITY-3654](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3654) fixed by the upgrade<br>[SECURITY-3678](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3678) fixed by the upgrade |
| matrix-auth | 3.1.5 | 3.3 | [SECURITY-3676](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3676) fixed by the upgrade |
| matrix-project | 772.v494f19991984 | 928.v1effc3f3a_236 | [SECURITY-3289](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3289) fixed by the upgrade |
| momentjs | 1.1.1 | 1.1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| nested-view | 1.24 | 241.v178f0b_a_cd76a_ | [SECURITY-2768](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2768) fixed by the upgrade |
| nodejs | 1.5.1 | 1.6.6 | [SECURITY-3196](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3196) fixed by the upgrade |
| oic-auth | 1.8 | 4.731.v4932b_506e767 | [SECURITY-2978](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-2978) fixed by the upgrade<br>[SECURITY-2979](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-2979) fixed by the upgrade<br>[SECURITY-3168](https://www.jenkins.io/security/advisory/2023-12-13/#SECURITY-3168) fixed by the upgrade<br>[SECURITY-3441-1](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3441%20(1)) fixed by the upgrade<br>[SECURITY-3441-2](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3441%20(2)) fixed by the upgrade<br>[SECURITY-3473](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3473) fixed by the upgrade<br>[SECURITY-3461](https://www.jenkins.io/security/advisory/2025-01-22/#SECURITY-3461) fixed by the upgrade |
| performance | 3.20 | 1022.v48fce959f0c3 | [SECURITY-2394](https://www.jenkins.io/security/advisory/2021-11-12/#SECURITY-2394) fixed by the upgrade<br>[SECURITY-4026](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4026) fixed by the upgrade |
| pipeline-build-step | 488.v8993df156e8d | 601.v6d4c6d1a_9dc7 | [SECURITY-3870](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3870) fixed by the upgrade |
| pipeline-groovy-lib | 612.v84da_9c54906d | 822.vc95f70dea_db_7 | [SECURITY-2824-2-pipeline-groovy-lib](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(2)) fixed by the upgrade<br>[SECURITY-3727](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3727) fixed by the upgrade<br>[SECURITY-3815](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3815) fixed by the upgrade<br>[SECURITY-3796](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3796) fixed by the upgrade |
| pipeline-maven | 3.10.0 | 1760.v9a_a_e6dcb_0444 | [SECURITY-3257](https://www.jenkins.io/security/advisory/2023-09-06/#SECURITY-3257) fixed by the upgrade |
| pipeline-model-definition | 1.9.3 | 2.2293.v6e7193cec599 | [SECURITY-3361](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3361) fixed by the upgrade |
| pipeline-utility-steps | 2.15.0 | 3.810.va_7672d206740 | [SECURITY-2196](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-2196) fixed by the upgrade |
| plain-credentials | 1.8 | 199.v9f8e1f741799 | [SECURITY-2495](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-2495) fixed by the upgrade |
| rebuild | 1.34 | 344.v805899551b_13 | [SECURITY-3033](https://www.jenkins.io/security/advisory/2023-07-12/#SECURITY-3033) fixed by the upgrade |
| robot | 3.3.0 | 6.3.0 | [SECURITY-3971](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3971) fixed by the upgrade |
| role-strategy | 483.v17281966f5c3 | 927.v9cf5527c4085 | [SECURITY-3053](https://www.jenkins.io/security/advisory/2023-03-21/#SECURITY-3053) fixed by the upgrade |
| saml | 2.0.9 | 4.623.v7875d61cd9f5 | [SECURITY-3613](https://www.jenkins.io/security/advisory/2025-10-29/#SECURITY-3613) fixed by the upgrade<br>[SECURITY-4065](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4065) fixed by the upgrade |
| script-security | 1138.v8e727069a_025 | 1436.v1082791135e4 | [SECURITY-2116](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2116) fixed by the upgrade<br>[SECURITY-2824-1-script-security](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-2564](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2564) fixed by the upgrade<br>[SECURITY-3016](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-3016) fixed by the upgrade<br>[SECURITY-3341](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3341) fixed by the upgrade<br>[SECURITY-3447](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3447) fixed by the upgrade<br>[SECURITY-3662](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3662) fixed by the upgrade<br>[SECURITY-3792](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3792) fixed by the upgrade<br>[SECURITY-3793](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3793) fixed by the upgrade<br>[SECURITY-3986](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3986) fixed by the upgrade<br>[SECURITY-3897](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3897) fixed by the upgrade<br>[SECURITY-3923](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3923) fixed by the upgrade<br>[SECURITY-3925-1](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(1)) fixed by the upgrade<br>[SECURITY-3925-2](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(2)) fixed by the upgrade<br>[SECURITY-3929](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3929) fixed by the upgrade<br>[SECURITY-3931](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3931) fixed by the upgrade<br>[SECURITY-3932](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3932) fixed by the upgrade<br>[SECURITY-3977](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3977) fixed by the upgrade |
| sonar | 2.14 | 2.19.0 | [SECURITY-3989](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3989) fixed by the upgrade |
| ssh | 2.6.1 | 159.v496ca_25e5e82 | [SECURITY-2093](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2093) **no fix yet**<br>[SECURITY-2315](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2315) **no fix yet** |
| structs | 324.va_f5d6774f3a_d | 362.va_b_695ef4fdf9 | [SECURITY-3371](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3371) fixed by the upgrade |
| thinBackup | 1.11 | 2.1.5 | [SECURITY-4099](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-4099) fixed by the upgrade |
| translation | 1.16 | 1.16 | [deprecated](https://github.com/jenkins-infra/update-center2/pull/719) |
| trilead-api | 1.0.13 | 2.284.v1974ea_324382 | [SECURITY-3333](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3333) fixed by the upgrade |
| warnings-ng | 9.11.1 | 13.10305.v3231d55c5755 | [SECURITY-3265](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3265) fixed by the upgrade<br>[SECURITY-3937](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3937) fixed by the upgrade |
| workflow-cps | 2660.2664.v4c114e93f4c1 | 4395.v738456a_5da_6b_ | [SECURITY-359](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-359) fixed by the upgrade<br>[SECURITY-2824-1-workflow-cps](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-3362](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3362) fixed by the upgrade<br>[SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-cps-global-lib | 2.21.3 | 636.v3d2c06d018a_6 | [SECURITY-2824-2-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(2)) fixed by the upgrade |
| workflow-job | 1145.v7f2433caa07f | 1609.v8926e7962024 | [SECURITY-3042](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3042) fixed by the upgrade |
| workflow-multibranch | 716.vc692a_e52371b_ | 842.v3a_b_59b_57b_e6e | [SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |
| workflow-support | 813.vb_d7c3d2984a_0 | 1015.v785e5a_b_b_8b_22 | [SECURITY-2881](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2881) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| Office-365-Connector | 4.17.0 | 5.3.1 |
| Parameterized-Remote-Trigger | 3.1.6.3 | 3.2.2 |
| ace-editor | 1.1 | 1.1 |
| active-directory | 2.27 | 2.971.v653d8b_6e5548 |
| allure-jenkins-plugin | 2.30.3 | 3.0.0 |
| amazon-ecr | 1.73.v741d474abe74 | 1.175.ve233e80e56fd |
| analysis-model-api | 10.10.1 | 14.18.0-1028.v9c938eeda_ef2 |
| ansible | 1.1 | 635.v34b_48e979c86 |
| ansicolor | 1.0.2 | 542.v03d235fee02d |
| ant | 481.v7b_09e538fcca | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 159.v25b_c67cd35fb_ | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.13-1.0 | 4.5.14-269.vfa_2321039a_83 |
| apache-httpcomponents-client-5-api | – | 5.6.4-207.vd5b_a_c15e4de5 (new, for Office-365-Connector, analysis-model-api, audit-trail, cloudbees-bitbucket-branch-source, docker-java-api, file-operations, http_request, saml, slack) |
| artifactory | 4.0.8 | 4.0.8 |
| asm-api | – | 9.11.1-226.vff833a_e4c58e (new, for analysis-model-api, cvs, jacoco, jnr-posix-api, json-path-api, oic-auth, scm-api) |
| atlassian-bitbucket-server-integration | 3.3.1 | 6.0.2 |
| audit-trail | 3.10 | 464.v78b_845ffc5d8 |
| authentication-tokens | 1.4 | 1.144.v5ff4a_5ec5c33 |
| authorize-project | 1.4.0 | 534.v2f208c45e11c |
| aws-credentials | 191.vcb_f183ce58b_9 | 265.v8422a_f384cd9 |
| aws-java-sdk | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-api-gateway | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-autoscaling | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-cloudformation | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-cloudfront | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-cloudwatch | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-codebuild | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-codedeploy | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-ec2 | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ecr | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ecs | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-efs | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-elasticbeanstalk | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-elasticloadbalancingv2 | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-iam | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-kinesis | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-lambda | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-logs | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-minimal | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-organizations | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk, pipeline-aws) |
| aws-java-sdk-secretsmanager | – | 1.12.780-480.v4a_0819121a_9e (new, for aws-java-sdk) |
| aws-java-sdk-sns | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-sqs | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ssm | 1.12.287-357.vf82d85a_6eefd | 1.12.780-480.v4a_0819121a_9e |
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
| azure-sdk | 118.v43f74dd9ca_dc | 299.v3996878fe83a_ |
| badge | 1.9.1 | 3.615.vb_ed967ca_b_e93 |
| basic-branch-build-strategies | 1.3.2 | 317.v85b_331d6cc42 |
| bitbucket | 214.v2fd4234d0554 | 322.v882c6a_b_82c52 |
| blueocean | 1.26.0 | 1.27.26 |
| blueocean-autofavorite | 1.2.5 | 1.2.5 |
| blueocean-bitbucket-pipeline | 1.26.0 | 1.27.26 |
| blueocean-commons | 1.26.0 | 1.27.26 |
| blueocean-config | 1.26.0 | 1.27.26 |
| blueocean-core-js | 1.26.0 | 1.27.26 |
| blueocean-dashboard | 1.26.0 | 1.27.26 |
| blueocean-display-url | 2.4.1 | 2.4.4 |
| blueocean-events | 1.26.0 | 1.27.26 |
| blueocean-git-pipeline | 1.26.0 | 1.27.26 |
| blueocean-github-pipeline | 1.26.0 | 1.27.26 |
| blueocean-i18n | 1.26.0 | 1.27.26 |
| blueocean-jira | 1.26.0 | 1.27.26 |
| blueocean-jwt | 1.26.0 | 1.27.26 |
| blueocean-personalization | 1.26.0 | 1.27.26 |
| blueocean-pipeline-api-impl | 1.26.0 | 1.27.26 |
| blueocean-pipeline-editor | 1.26.0 | 1.27.26 |
| blueocean-pipeline-scm-api | 1.26.0 | 1.27.26 |
| blueocean-rest | 1.26.0 | 1.27.26 |
| blueocean-rest-impl | 1.26.0 | 1.27.26 |
| blueocean-web | 1.26.0 | 1.27.26 |
| bootstrap4-api | 4.6.0-6 | 4.6.0-6 |
| bootstrap5-api | 5.1.3-6 | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.25 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1046.v0ca_37783ecc5 | 2.1303.v9f3b_95dc329d |
| build-blocker-plugin | 1.7.9 | 177.ve433d85b_3fd7 |
| build-monitor-plugin | 1.13+build.202202112334 | 1131.va_2d8e68fc912 |
| build-name-setter | 2.2.1 | 2.5.1 |
| build-pipeline-plugin | 1.5.8 | 2.2.0 |
| build-timeout | 1.20 | 1.41 |
| build-timestamp | 1.0.3 | 1.1.1 |
| build-token-root | 151.va_e52fe3215fc | 365.v717f8685a_09e |
| build-user-vars-plugin | 1.8 | 214.va_eed2ed849ca_ |
| build-with-parameters | 73.v3ed13cb_d6cca | 81.ve4a_9c2499d9a |
| built-on-column | 1.1 | 1.5 |
| caffeine-api | 2.9.3-65.v6a_47d0f4d1fe | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 1.7.4 | 422.vc7e8b_51488c6 |
| cloud-stats | 0.27 | 423.v63f017a_f582c |
| cloudbees-bitbucket-branch-source | 762.v969cfe087fc0 | 937.3.10 |
| cloudbees-disk-usage-simple | 0.10 | 275.v56e74c9dd0d3 |
| cloudbees-folder | 6.722.v8165b_a_cf25e9 | 6.1106.v3a_d9a_6d2465e |
| cobertura | 1.17 | 1.17 |
| code-coverage-api | 2.0.4 | 4.99.0 |
| command-launcher | 81.v9c2cb_cb_db_392 | 141.v9578cf36b_311 |
| commons-collections4-api | – | 4.5.0-8.va_d5448ef9011 (new, for cloudbees-bitbucket-branch-source, configurationslicing, cucumber-reports, jira) |
| commons-compress-api | – | 1.28.0-96.v8da_8cfa_57e2b_ (new, for docker-java-api, pipeline-utility-steps) |
| commons-httpclient3-api | – | 3.1-3 (new, for nexus-artifact-uploader) |
| commons-lang3-api | 3.12.0-36.vd97de6465d5b_ | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.10.0-27.vb_fa_3896786a_7 | 1.15.0-218.va_61573470393 |
| conditional-buildstep | 1.4.2 | 1.5.0 |
| config-file-provider | 3.8.1 | 1013.v73c323e52b_1f |
| configuration-as-code | 1512.vb_79d418d5fc8 | 2136.veb_3a_554550b_f |
| configurationslicing | 430.v966357576543 | 715.vcc64f669003e |
| copyartifact | 1.46.4 | 805.v9d1393350780 |
| coverage | – | 3.3497.v3fee733ca_b_35 (new, for code-coverage-api, pipeline-maven, prometheus) |
| credentials | 2.6.1.1 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 1.27.1 | 728.v902a_273b_8947 |
| cucumber-reports | 5.7.5 | 5.11.0 |
| cvs | 2.19.1 | 502.v02d5a_0c6b_9a_3 |
| dark-theme | 0.0.12 | 652.vea_da_dfea_e769 |
| dashboard-view | 2.447.vda_124dd35f11 | 2.569.v5ddb_c9b_a_7264 |
| data-tables-api | 1.11.4-4 | 2.3.8-1570.v1cb_1cd2a_0fb_c |
| dependency-check-jenkins-plugin | 5.2.0 | 5.6.5 |
| deploy | 1.16 | 329.vd9ec45431fb_e |
| description-setter | 1.10 | 264.v1957f215dcd5 |
| disk-usage | 0.28 | 573.v5269a_8d16a_db_ |
| display-url-api | 2.3.5 | 2.217.va_6b_de84cc74b_ |
| docker-build-step | 2.8 | 2.14 |
| docker-commons | 1.21 | 477.v289085a_b_6896 |
| docker-java-api | 3.2.13-37.vf3411c9828b9 | 3.7.1-136.v5d70f77a_c3d6 |
| docker-plugin | 1.2.9 | 1336.v3711dec712de |
| docker-workflow | 1.28 | 653.v2f2c08eff0ec |
| dtkit-api | 3.0.1 | 3.0.3 |
| durable-task | 507.v050055d0cb_dd | 686.v80ff80875b_82 |
| ec2 | 2.0.4 | 2064.v79de4fd2d03b_ |
| echarts-api | 5.3.2-1 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | – | 0.3.0.1-29.v67e9a_1c969b_b_ (new, for atlassian-bitbucket-server-integration, mina-sshd-api-common, sshd, subversion, trilead-api) |
| email-ext | 2.87 | 2038.v7b_8817a_499d9 |
| emailext-template | 1.4 | 249.v81d6ea_f33f2e |
| embeddable-build-status | 2.0.4 | 666.vcb_0de9eb_43d8 |
| emoji-symbols-api | – | 17.0-57.v8d44b_9a_b_d5ea_ (new, for badge) |
| envinject | 2.4.0 | 2.941.v351a_20c0a_3ca_ |
| envinject-api | 1.8 | 1.241.vdd714803b_403 |
| extended-choice-parameter | 346.vd87693c5a_86c | 388.ve7b_d0b_920e10 |
| extended-read-permission | 3.2 | 92.v1b_f0c20b_75db_ |
| extensible-choice-parameter | 1.8.0 | 255.vfa_41f46828ef |
| external-monitor-job | 203.v683c09d993b_9 | 223.vb_fddcf42c9b_3 |
| favorite | 2.3.3.1 | 2.267.vb_90d08408081 |
| file-operations | 1.11 | 439.vdc9d3b_d74b_d9 |
| file-parameters | 205.vf6ce13b_e5dee | 433.va_0b_80359d54d |
| folder-properties | 1.2.1 | 97.v0c61a_c452eb_b_ |
| font-awesome-api | 6.0.0-1 | 7.3.1-1013.v0835a_879ec6d |
| forensics-api | 1.13.0 | 4.1923.v445c56df9b_2e |
| generic-webhook-trigger | 1.86.4 | 2.4.3 |
| git | 4.11.5 | 5.11.0 |
| git-client | 3.11.2 | 6.6.1 |
| git-parameter | 0.9.14 | 462.463.v496a_59f698e5 |
| git-server | 1.11 | 137.ve0060b_432302 |
| github | 1.34.3.1 | 1.47.0 |
| github-api | 1.303-400.v35c2d8258028 | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 2.11.4 | 1986.v801e40614d00 |
| github-oauth | 0.39 | 685.v53b_070455063 |
| github-pullrequest | 0.4.0 | 0.7.5 |
| gitlab-api | 5.1.0-84.v491924123a_f7 | 6.3.0-126.vd066b_d896159 |
| gitlab-branch-source | 650.va_d1ce6d01959 | 744.vb_d0403d08ec7 |
| gitlab-oauth | 1.16 | 1.30 |
| gitlab-plugin | 1.5.36 | 1.2158.v49db_f64b_fc28 |
| google-oauth-plugin | 1.0.6 | 1.346.v712b_a_e0d1e60 |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ |
| greenballs | 1.15.1 | 1.15.1 |
| groovy | 2.4 | 537.v741a_5a_f1b_581 |
| groovy-postbuild | 2.5 | 335.vf10eb_c1b_5e14 |
| gson-api | – | 2.14.0-201.v8eefe5515533 (new, for Office-365-Connector, analysis-model-api, generic-webhook-trigger, git-client, google-oauth-plugin, performance, trilead-api, warnings-ng) |
| h2-api | 1.4.199 | 11.1.4.199-36.vb_ee07e965744 |
| handy-uri-templates-2-api | 2.1.8-22.v77d5b_75e6953 | 2.1.8-38.vcea_5d521d5f3 |
| hashicorp-vault-plugin | 336.v182c0fbaaeb7 | 384.vda_86ec66c537 |
| hidden-parameter | 0.0.5 | 546.vdeec4ce8b_26f |
| htmlpublisher | 1.31 | 429 |
| http_request | 1.16 | 1.659.v1b_9d9e942a_e6 |
| icon-shim | 3.0.0 | 3.0.0 |
| instance-identity | – | 203.v15e81a_1b_7a_38 (new, for atlassian-bitbucket-server-integration, build-monitor-plugin, docker-plugin, git-server, github, mailer, pubsub-light, sshd) |
| ionicons-api | – | 94.vcc3065403257 (new, for active-directory, badge, branch-api, cloudbees-disk-usage-simple, cloudbees-folder, code-coverage-api, config-file-provider, coverage, credentials, dashboard-view, docker-plugin, emailext-template, embeddable-build-status, favorite, git, github-branch-source, jobConfigHistory, junit, matrix-auth, matrix-project, metrics, naginator, nested-view, next-build-number, pipeline-graph-view, pipeline-model-definition, role-strategy, test-results-analyzer, theme-manager, thinBackup, workflow-cps, workflow-job, workflow-multibranch) |
| ivy | 2.2 | 609.v66779e68503d |
| jackson-annotations2-api | – | 2.22-19.v10a_a_582ea_26e (new, for jackson2-api, jackson3-api) |
| jackson2-api | 2.14.2-319.v37853346a_229 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | – | 3.2.4-103.vf5a_b_3f4ea_cec (new, for coverage, data-tables-api, echarts-api, file-operations, forensics-api, junit, pipeline-graph-view, pipeline-utility-steps, theme-manager) |
| jacoco | 3.3.2.1 | 3.3.7 |
| jakarta-activation-api | 2.0.0-2 | 2.1.4-1 |
| jakarta-mail-api | 2.0.0-5 | 2.1.5-1 |
| jakarta-xml-bind-api | – | 4.0.9-19.v2b_a_5b_44d9a_1c (new, for gitlab-oauth, gitlab-plugin, jackson2-api, jersey3-api, jobConfigHistory) |
| javadoc | 217.v905b_86277a_2a_ | 354.vee1a_660b_4990 |
| javax-activation-api | 1.2.0-3 | 1.2.0-8 |
| javax-mail-api | 1.6.2-6 | 1.6.2-11 |
| jaxb | 2.3.6-1 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 1.5 | 83.v417146707a_3d |
| jenkins-design-language | 1.26.0 | 1.27.26 |
| jenkins-multijob-plugin | 1.36 | 682.v4ed91b_52d621 |
| jersey2-api | 2.36-2 | 2.48-180.ve47b_264f849b_ |
| jersey3-api | – | 3.1.11-4.v77818819c2e1 (new, for gitlab-api, gitlab-plugin) |
| jira | 3.7.1 | 3.22 |
| jjwt-api | 0.11.5-77.v646c772fddb_0 | 0.13.0-141.vd58b_a_9592b_6c |
| jnr-posix-api | 3.0.45 | 3.2.4-223.v4f0cfd7fa_006 |
| job-dsl | 1.81.1 | 3732.v9a_c49a_61a_313 |
| job-import-plugin | 3.6 | 143.145.v48f9a_a_6ff384 |
| jobConfigHistory | 2.30 | 1380.v762185b_9a_793 |
| joda-time-api | – | 2.15.0-208.v6b_084e7a_d9f0 (new, for aws-java-sdk-minimal, cucumber-reports, gitlab-plugin, handy-uri-templates-2-api, jira, pipeline-model-api, saml) |
| jquery | 1.12.4-1 | 1.12.4-3 |
| jquery-detached | 1.2.1 | 1.2.1 |
| jquery-ui | 1.0.2 | 1.0.2 |
| jquery3-api | 3.6.0-3 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.1.55.2 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | – | 20260814-226.v20f9685d642c (new, for analysis-model-api, configuration-as-code, jackson2-api, jira, warnings-ng) |
| json-path-api | – | 3.0.0-218.vcd4dd1355de2 (new, for analysis-model-api, azure-sdk, build-timeout, generic-webhook-trigger, oic-auth, token-macro, warnings-ng) |
| jsoup | 1.10.3 | 1.23.2-106.v1161f282b_eb_f |
| junit | 1.53.0.1 | 1442.vb_b_a_467df2767 |
| kubernetes | 1.30.11 | 4560.vc25b_83890c34 |
| kubernetes-cli | 1.10.3 | 1.375.v8001982544ba |
| kubernetes-client-api | 5.12.2-193.v26a_6078f65a_9 | 7.3.1-256.v788a_0b_787114 |
| kubernetes-credentials | 0.9.0 | 207.v492f58828b_ed |
| kubernetes-credentials-provider | 0.20 | 1.315.v92008589c044 |
| ldap | 2.9 | 825.v2fca_37dd5b_cb_ |
| locale | 180.v207501dff9b_a_ | 641.v84f22d75fd8b_ |
| localization-support | 1.2 | 1.41.v186a_c1569458 |
| localization-zh-cn | 1.0.24 | 371.v23851f835d6b_ |
| lockable-resources | 2.18 | 1576.v217259b_f9976 |
| log-parser | 2.2 | 3.0.7 |
| m2release | 0.16.3 | 547.v78681a_45411a_ |
| mailer | 414.vcc4c33714601 | 534.v1b_36f5864073 |
| mapdb-api | 1.0.9.0 | 1.0.9-44.va_1e1310c9118 |
| mask-passwords | 3.3 | 220.v95819055b_265 |
| matrix-auth | 3.1.5 | 3.3 |
| matrix-project | 772.v494f19991984 | 928.v1effc3f3a_236 |
| maven-plugin | 3.16 | 3.28 |
| mercurial | 1260.vdfb_723cdcc81 | 1323.ve69d2a_db_8a_b_d |
| metrics | 4.0.2.8.1 | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.8.0-30.vf9df64641cb_d | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.8.0-30.vf9df64641cb_d | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-scp | 2.8.0-30.vf9df64641cb_d | 2.19.0-192.v2b_a_7b_2c1dc71 |
| momentjs | 1.1.1 | 1.1.1 |
| monitoring | 1.95.0 | 2.9.1 |
| msbuild | 1.30 | 1.45 |
| multiple-scms | 0.8 | 0.8 |
| naginator | 1.18.2 | 1.556.v14d723a_109a_c |
| nested-view | 1.24 | 241.v178f0b_a_cd76a_ |
| netty-api | – | 4.2.18.Final-34.v861b_21a_78c0f (new, for aws-java-sdk2-netty-nio-client, azure-sdk) |
| next-build-number | 1.8 | 66.v4b_4762172d53 |
| nexus-artifact-uploader | 2.13 | 2.16 |
| nimbus-jose-jwt-api | – | 10.10-12.vd4e6d651b_7f7 (new, for oic-auth) |
| node-iterator-api | 1.5.1 | 72.vc90e81737df1 |
| nodejs | 1.5.1 | 1.6.6 |
| nodelabelparameter | 1.10.3.1 | 851.vd94e5048d321 |
| oauth-credentials | 0.5 | 0.668.vb_431feb_d8a_4e |
| oic-auth | 1.8 | 4.731.v4932b_506e767 |
| okhttp-api | 4.9.3-108.v0feda04578cf | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | – | 477.v7197b_eb_3d8b_8 (new, for javadoc, maven-plugin, pipeline-maven-api) |
| pam-auth | 1.6.1 | 1.12 |
| parameter-separator | 1.3 | 334.v1fc0e534c71d |
| parameterized-scheduler | 1.1 | 379.v95c73f233a_df |
| parameterized-trigger | 2.44 | 893.va_383a_9b_a_4c11 |
| performance | 3.20 | 1022.v48fce959f0c3 |
| pipeline-aws | 1.43 | 767.v48d5c10b_1899 |
| pipeline-build-step | 488.v8993df156e8d | 601.v6d4c6d1a_9dc7 |
| pipeline-github | 2.8-141.7e985db66f30 | 2.8-191.cb_1a_07f10b_a_f |
| pipeline-github-lib | 36.v4c01db_ca_ed16 | 65.v203688e7727e |
| pipeline-graph-analysis | 1.11 | 254.v0f63a_a_447dca_ |
| pipeline-graph-view | 61.v5f35529d13b7 | 1053.v9df19e6ceed2 |
| pipeline-groovy-lib | 612.v84da_9c54906d | 822.vc95f70dea_db_7 |
| pipeline-input-step | 466.v6d0a_5df34f81 | 560.v56198a_642157 |
| pipeline-maven | 3.10.0 | 1760.v9a_a_e6dcb_0444 |
| pipeline-maven-api | – | 1760.v9a_a_e6dcb_0444 (new, for pipeline-maven) |
| pipeline-milestone-step | 111.v449306f708b_7 | 152.v6e22b_8cfc66c |
| pipeline-model-api | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-model-declarative-agent | 1.1.1 | 1.1.1 |
| pipeline-model-definition | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-multibranch-defaults | 2.1 | 2.1 |
| pipeline-rest-api | 2.28 | 2.41 |
| pipeline-stage-step | 293.v200037eefcd5 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.28 | 2.41 |
| pipeline-utility-steps | 2.15.0 | 3.810.va_7672d206740 |
| plain-credentials | 1.8 | 199.v9f8e1f741799 |
| plugin-usage-plugin | 2.2 | 418.v308c3c863d25 |
| plugin-util-api | 2.16.0 | 7.1403.vd380b_c42e90c |
| popper-api | 1.16.1-3 | 1.16.1-3 |
| popper2-api | 2.11.5-1 | 2.11.6-5 |
| postbuild-task | 1.9 | 78.v24529f1f5cdb_ |
| postbuildscript | 3.1.0-375.v3db_cd92485e1 | 3.4.1.821.v08352159216f |
| powershell | 1.8 | 185.v7a_026da_c54ee |
| prism-api | 1.26.0-2 | 1.30.0-741.v034eb_0b_0a_a_fa_ |
| prometheus | 2.0.10 | 860.v532442b_44e9a_ |
| promoted-builds | 881.vb_72da_69ffca_1 | 1029.v979cc7835e39 |
| publish-over | 0.22 | 242.vf56f6892cb_b_a_ |
| publish-over-ssh | 1.24 | 409.v73ca_23a_5b_48e |
| pubsub-light | 1.16 | 1.19 |
| python | 1.3 | 1.3 |
| rebuild | 1.34 | 344.v805899551b_13 |
| resource-disposer | 0.20 | 0.25 |
| robot | 3.3.0 | 6.3.0 |
| role-strategy | 483.v17281966f5c3 | 927.v9cf5527c4085 |
| run-condition | 1.5 | 276.v97298f3a_cd51 |
| s3 | 0.12.3445.vda_704535b_5a_d | 522.vc73d01440c29 |
| saferestart | 0.3 | 125.vdb_73765c38a_e |
| saml | 2.0.9 | 4.623.v7875d61cd9f5 |
| scm-api | 608.vfa_f971c5a_a_e9 | 728.vc30dcf7a_0df5 |
| script-security | 1138.v8e727069a_025 | 1436.v1082791135e4 |
| simple-theme-plugin | 103.va_161d09c38c7 | 230.v8b_fd91b_b_800c |
| slack | 631.v40deea_40323b | 795.v4b_9705b_e6d47 |
| snakeyaml-api | 1.31-84.ve43da_fb_49d0b | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | – | 3.2-18.v885ed7ff698f (new, for jackson3-api) |
| sonar | 2.14 | 2.19.0 |
| sse-gateway | 1.25 | 1.29 |
| ssh | 2.6.1 | 159.v496ca_25e5e82 |
| ssh-agent | 295.v9ca_a_1c7cc3a_a_ | 433.v9e73d67c3f74 |
| ssh-credentials | 277.v95c2fec1c047 | 372.va_250881b_08cd |
| ssh-slaves | 1.806.v2253cedd3295 | 3.1097.v868116049892 |
| ssh-steps | 2.0.0 | 2.0.92.vb_a_0583935f9b_ |
| sshd | 3.242.va_db_9da_b_26a_c3 | 3.384.vc89b_5e138cf9 |
| stashNotifier | 1.24 | 1.557.vb_d02a_d29f87f |
| structs | 324.va_f5d6774f3a_d | 362.va_b_695ef4fdf9 |
| subversion | 2.16.0 | 1327.vc36e0019a_653 |
| test-results-analyzer | 0.3.5 | 353.vea_f73668d8b_5 |
| theme-manager | 0.6 | 346.v06cca_64c6a_37 |
| thinBackup | 1.11 | 2.1.5 |
| throttle-concurrents | 2.10 | 625.vc8b_e469e9a_b_c |
| timestamper | 1.20 | 1.30 |
| token-macro | 321.vd7cc1f2a_52c8 | 477.vd4f0dc3cb_cf1 |
| translation | 1.16 | 1.16 |
| trilead-api | 1.0.13 | 2.284.v1974ea_324382 |
| uno-choice | 2.6.1 | 2.8.11 |
| validating-string-parameter | 2.8 | 345.v0f4330488a_f6 |
| variant | 1.4 | 70.va_d9f17f859e0 |
| view-job-filters | 2.3 | 408.vb_ca_5e798e5b_5 |
| warnings-ng | 9.11.1 | 13.10305.v3231d55c5755 |
| woodstox-core-api | – | 7.3.0-16.v017f30ed004c (new, for jackson2-api, jackson3-api) |
| workflow-aggregator | 578.vf9a_f99755f4a_ | 608.v67378e9d3db_1 |
| workflow-api | 1153.vb_912c0e47fb_a_ | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 2.24 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 2660.2664.v4c114e93f4c1 | 4395.v738456a_5da_6b_ |
| workflow-cps-global-lib | 2.21.3 | 636.v3d2c06d018a_6 |
| workflow-durable-task-step | 2.40 | 1479.v56e587f413a_7 |
| workflow-job | 1145.v7f2433caa07f | 1609.v8926e7962024 |
| workflow-multibranch | 716.vc692a_e52371b_ | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 408.v7d5b_135a_b_d49 | 487.v41313a_96b_65d |
| workflow-step-api | 639.v6eca_cd8c04a_a_ | 724.v538c2362b_dfb_ |
| workflow-support | 813.vb_d7c3d2984a_0 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.43 | 0.49 |
| xunit | 3.1.0 | 3.1.7 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
