# Jenkins upgrade plan: 2.249.3 → 2.580.1

Source: results/2.249.3-top100/plugins-2.249.3.txt · generated 2026-10-21 13:34 UTC by jenkins-upgrade-check dev

**3 core upgrade(s)**, 134 plugins: 104 to upgrade, 27 new dependencies; 34 with security warnings today, 3 not from the update center.

Jenkins 2.249.3 runs on Java 8, 11; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

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

1. **On 2.249.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 11**, upgrade Jenkins to **2.462.3**, start it, then update plugins (122 change(s)).
3. **Run the controller and every agent on Java 17**, upgrade Jenkins to **2.541.3**, start it, then update plugins (128 change(s)).
4. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (15 change(s)).

No Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| ace-editor | 1.1 | 1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| cloudbees-bitbucket-branch-source | 726.vb0c1ea6c9336 | 937.3.10 | [SECURITY-2033](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2033) fixed by the upgrade<br>[SECURITY-2467](https://www.jenkins.io/security/advisory/2022-01-12/#SECURITY-2467) fixed by the upgrade<br>[SECURITY-3300](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3300) fixed by the upgrade<br>[SECURITY-3363](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3363) fixed by the upgrade |
| cloudbees-folder | 6.688.vfc7a_a_69059e0 | 6.1106.v3a_d9a_6d2465e | [SECURITY-3105](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3105) fixed by the upgrade<br>[SECURITY-3106](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3106) fixed by the upgrade<br>[SECURITY-3109](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3109) fixed by the upgrade |
| config-file-provider | 3.8.0 | 1013.v73c323e52b_1f | [SECURITY-3090](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3090) fixed by the upgrade |
| credentials | 2.6.1.1 | 1514.vb_01eb_014e7d4 | [SECURITY-3373](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3373) fixed by the upgrade |
| credentials-binding | 1.27.1 | 728.v902a_273b_8947 | [SECURITY-3499](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3499) fixed by the upgrade<br>[SECURITY-3672](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3672) fixed by the upgrade<br>[SECURITY-3790](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3790) fixed by the upgrade |
| email-ext | 2.84 | 2038.v7b_8817a_499d9 | [SECURITY-2931](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2931) fixed by the upgrade<br>[SECURITY-2934](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2934) fixed by the upgrade<br>[SECURITY-2939](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2939) fixed by the upgrade<br>[SECURITY-3088-1](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(1)) fixed by the upgrade<br>[SECURITY-3088-2](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(2)) fixed by the upgrade<br>[SECURITY-3705](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3705) fixed by the upgrade |
| favorite | 2.3.2 | 2.267.vb_90d08408081 | [SECURITY-2557](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2557) fixed by the upgrade |
| git | 4.6.0 | 5.11.0 | [SECURITY-2499](https://www.jenkins.io/security/advisory/2021-10-06/#SECURITY-2499) fixed by the upgrade<br>[SECURITY-2478-git](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2478) fixed by the upgrade<br>[SECURITY-284](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-284) fixed by the upgrade<br>[SECURITY-2796](https://www.jenkins.io/security/advisory/2022-08-23/#SECURITY-2796) fixed by the upgrade |
| git-client | 3.6.0 | 6.6.1 | [SECURITY-1468](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-1468) fixed by the upgrade<br>[SECURITY-3590](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3590) fixed by the upgrade<br>[SECURITY-3614](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3614) fixed by the upgrade<br>[SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| git-server | 1.9 | 137.ve0060b_432302 | [SECURITY-3319](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3319) fixed by the upgrade<br>[SECURITY-3342](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3342) fixed by the upgrade |
| github | 1.34.2 | 1.47.0 | [SECURITY-1849](https://www.jenkins.io/security/advisory/2022-07-27/#SECURITY-1849) fixed by the upgrade<br>[SECURITY-3246](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3246) fixed by the upgrade<br>[SECURITY-3704](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3704) fixed by the upgrade |
| github-branch-source | 2.9.9 | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| handlebars | 3.0.8 | 3.0.8 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| htmlpublisher | 1.25 | 429 | [SECURITY-3301](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3301) fixed by the upgrade<br>[SECURITY-3302](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3302) fixed by the upgrade<br>[SECURITY-3303](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3303) fixed by the upgrade<br>[SECURITY-3547](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3547) fixed by the upgrade<br>[SECURITY-3706](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3706) fixed by the upgrade |
| junit | 1.51 | 1442.vb_b_a_467df2767 | [SECURITY-2760](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2760) fixed by the upgrade<br>[SECURITY-2888](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2888) fixed by the upgrade<br>[SECURITY-3032](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-3032) fixed by the upgrade |
| ldap | 1.26 | 825.v2fca_37dd5b_cb_ | [SECURITY-3046](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3046) fixed by the upgrade<br>[SECURITY-3654](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3654) fixed by the upgrade<br>[SECURITY-3678](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3678) fixed by the upgrade |
| matrix-auth | 3.0.1 | 3.3 | [SECURITY-3676](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3676) fixed by the upgrade |
| matrix-project | 1.18.1 | 928.v1effc3f3a_236 | [SECURITY-3289](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3289) fixed by the upgrade |
| momentjs | 1.1.1 | 1.1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| nodejs | 1.4.3 | 1.6.6 | [SECURITY-3196](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3196) fixed by the upgrade |
| pipeline-build-step | 2.17 | 601.v6d4c6d1a_9dc7 | [SECURITY-3019](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-3019) fixed by the upgrade<br>[SECURITY-3870](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3870) fixed by the upgrade |
| pipeline-input-step | 427.va6441fa17010 | 560.v56198a_642157 | [SECURITY-2705](https://www.jenkins.io/security/advisory/2022-06-22/#SECURITY-2705) fixed by the upgrade<br>[SECURITY-2880](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2880) fixed by the upgrade |
| pipeline-model-definition | 1.9.3 | 2.2293.v6e7193cec599 | [SECURITY-3361](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3361) fixed by the upgrade |
| pipeline-stage-view | 2.23 | 2.41 | [SECURITY-2828](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2828) fixed by the upgrade |
| pipeline-utility-steps | 2.11.0 | 3.810.va_7672d206740 | [SECURITY-2948](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2948) fixed by the upgrade<br>[SECURITY-2949](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2949) fixed by the upgrade<br>[SECURITY-2196](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-2196) fixed by the upgrade |
| plain-credentials | 1.8 | 199.v9f8e1f741799 | [SECURITY-2495](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-2495) fixed by the upgrade |
| role-strategy | 483.v17281966f5c3 | 927.v9cf5527c4085 | [SECURITY-3053](https://www.jenkins.io/security/advisory/2023-03-21/#SECURITY-3053) fixed by the upgrade |
| script-security | 1138.v8e727069a_025 | 1436.v1082791135e4 | [SECURITY-2116](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2116) fixed by the upgrade<br>[SECURITY-2824-1-script-security](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-2564](https://www.jenkins.io/security/advisory/2022-11-15/#SECURITY-2564) fixed by the upgrade<br>[SECURITY-3016](https://www.jenkins.io/security/advisory/2023-01-24/#SECURITY-3016) fixed by the upgrade<br>[SECURITY-3341](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3341) fixed by the upgrade<br>[SECURITY-3447](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3447) fixed by the upgrade<br>[SECURITY-3662](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3662) fixed by the upgrade<br>[SECURITY-3792](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3792) fixed by the upgrade<br>[SECURITY-3793](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3793) fixed by the upgrade<br>[SECURITY-3986](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3986) fixed by the upgrade<br>[SECURITY-3897](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3897) fixed by the upgrade<br>[SECURITY-3923](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3923) fixed by the upgrade<br>[SECURITY-3925-1](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(1)) fixed by the upgrade<br>[SECURITY-3925-2](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(2)) fixed by the upgrade<br>[SECURITY-3929](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3929) fixed by the upgrade<br>[SECURITY-3931](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3931) fixed by the upgrade<br>[SECURITY-3932](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3932) fixed by the upgrade<br>[SECURITY-3977](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3977) fixed by the upgrade |
| structs | 308.v852b473a2b8c | 362.va_b_695ef4fdf9 | [SECURITY-3371](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3371) fixed by the upgrade |
| trilead-api | 1.0.13 | 2.284.v1974ea_324382 | [SECURITY-3333](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3333) fixed by the upgrade |
| workflow-cps | 2660.2664.v4c114e93f4c1 | 4395.v738456a_5da_6b_ | [SECURITY-359](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-359) fixed by the upgrade<br>[SECURITY-2824-1-workflow-cps](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(1)) fixed by the upgrade<br>[SECURITY-3362](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3362) fixed by the upgrade<br>[SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-cps-global-lib | 2.19 | 636.v3d2c06d018a_6 | [SECURITY-2422](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2422) fixed by the upgrade<br>[SECURITY-2441](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2441) fixed by the upgrade<br>[SECURITY-2463-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2463) fixed by the upgrade<br>[SECURITY-2586](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2586) fixed by the upgrade<br>[SECURITY-2613-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2613) fixed by the upgrade<br>[SECURITY-1951](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-1951) fixed by the upgrade<br>[SECURITY-2824-2-workflow-cps-global-lib](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2824%20(2)) fixed by the upgrade |
| workflow-job | 1145.v7f2433caa07f | 1609.v8926e7962024 | [SECURITY-3042](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3042) fixed by the upgrade |
| workflow-multibranch | 2.24 | 842.v3a_b_59b_57b_e6e | [SECURITY-2463-workflow-multibranch](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2463) fixed by the upgrade<br>[SECURITY-2613-workflow-multibranch](https://www.jenkins.io/security/advisory/2022-02-15/#SECURITY-2613) fixed by the upgrade<br>[SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |
| workflow-support | 3.8 | 1015.v785e5a_b_b_8b_22 | [SECURITY-2881](https://www.jenkins.io/security/advisory/2022-10-19/#SECURITY-2881) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| ace-editor | 1.1 | 1.1 |
| ant | 1.13 | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 2.1 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.13-1.0 | 4.5.14-269.vfa_2321039a_83 |
| apache-httpcomponents-client-5-api | – | 5.6.4-207.vd5b_a_c15e4de5 (new, for cloudbees-bitbucket-branch-source) |
| asm-api | – | 9.11.1-226.vff833a_e4c58e (new, for json-path-api, scm-api) |
| authentication-tokens | 1.4 | 1.144.v5ff4a_5ec5c33 |
| blueocean-commons | 1.24.8 | 1.27.26 |
| blueocean-rest | 1.24.8 | 1.27.26 |
| bootstrap4-api | 4.6.0-3 | 4.6.0-6 |
| bootstrap5-api | 5.0.1-2 | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.25 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.6.2 | 2.1303.v9f3b_95dc329d |
| build-timeout | 1.20 | 1.41 |
| caffeine-api | 2.9.2-29.v717aac953ff3 | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 1.7.0 | 422.vc7e8b_51488c6 |
| cloudbees-bitbucket-branch-source | 726.vb0c1ea6c9336 | 937.3.10 |
| cloudbees-folder | 6.688.vfc7a_a_69059e0 | 6.1106.v3a_d9a_6d2465e |
| command-launcher | 1.6 | 141.v9578cf36b_311 |
| commons-collections4-api | – | 4.5.0-8.va_d5448ef9011 (new, for cloudbees-bitbucket-branch-source) |
| commons-compress-api | – | 1.28.0-96.v8da_8cfa_57e2b_ (new, for pipeline-utility-steps) |
| commons-lang3-api | – | 3.21.0-122.v6a_ea_4fd34049 (new, for ant, authentication-tokens, blueocean-commons, bootstrap5-api, branch-api, checks-api, cloudbees-bitbucket-branch-source, commons-compress-api, commons-text-api, credentials, data-tables-api, echarts-api, email-ext, font-awesome-api, git, git-client, github, github-api, github-branch-source, gradle, htmlpublisher, ionicons-api, jquery3-api, maven-plugin, metrics, nodejs, pipeline-graph-analysis, pipeline-model-api, pipeline-model-definition, pipeline-model-extensions, pipeline-rest-api, pipeline-utility-steps, plugin-util-api, popper2-api, prism-api, role-strategy, scm-api, ssh-slaves, timestamper, token-macro, workflow-basic-steps) |
| commons-text-api | – | 1.15.0-218.va_61573470393 (new, for blueocean-commons, bootstrap5-api, checks-api, credentials, data-tables-api, echarts-api, email-ext, font-awesome-api, jquery3-api, pipeline-model-api, pipeline-utility-steps, plugin-util-api, popper2-api, prism-api, token-macro) |
| config-file-provider | 3.8.0 | 1013.v73c323e52b_1f |
| credentials | 2.6.1.1 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 1.27.1 | 728.v902a_273b_8947 |
| dark-theme | 0.0.10 | 652.vea_da_dfea_e769 |
| data-tables-api | 1.10.25-1 | 2.3.8-1570.v1cb_1cd2a_0fb_c |
| display-url-api | 2.3.5 | 2.217.va_6b_de84cc74b_ |
| docker-commons | 1.19 | 477.v289085a_b_6896 |
| docker-workflow | 1.28 | 653.v2f2c08eff0ec |
| durable-task | 1.37 | 686.v80ff80875b_82 |
| echarts-api | 5.1.2-2 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | – | 0.3.0.1-29.v67e9a_1c969b_b_ (new, for mina-sshd-api-common, sshd, trilead-api) |
| email-ext | 2.84 | 2038.v7b_8817a_499d9 |
| favorite | 2.3.2 | 2.267.vb_90d08408081 |
| font-awesome-api | 5.15.3-3 | 7.3.1-1013.v0835a_879ec6d |
| git | 4.6.0 | 5.11.0 |
| git-client | 3.6.0 | 6.6.1 |
| git-server | 1.9 | 137.ve0060b_432302 |
| github | 1.34.2 | 1.47.0 |
| github-api | 1.301-378.v9807bd746da5 | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 2.9.9 | 1986.v801e40614d00 |
| gradle | 1.40 | 2.20.1253.vc116f0763a_eb_ |
| gson-api | – | 2.14.0-201.v8eefe5515533 (new, for git-client, trilead-api) |
| handlebars | 3.0.8 | 3.0.8 |
| handy-uri-templates-2-api | 2.1.8-1.0 | 2.1.8-38.vcea_5d521d5f3 |
| htmlpublisher | 1.25 | 429 |
| instance-identity | – | 203.v15e81a_1b_7a_38 (new, for git-server, github, mailer, pubsub-light, sshd) |
| ionicons-api | – | 94.vcc3065403257 (new, for branch-api, cloudbees-folder, config-file-provider, credentials, favorite, git, github-branch-source, junit, matrix-auth, matrix-project, metrics, pipeline-model-definition, role-strategy, theme-manager, workflow-cps, workflow-job, workflow-multibranch) |
| jackson-annotations2-api | – | 2.22-19.v10a_a_582ea_26e (new, for jackson2-api, jackson3-api) |
| jackson2-api | 2.13.2.20220328-273.v11d70a_b_a_1a_52 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | – | 3.2.4-103.vf5a_b_3f4ea_cec (new, for data-tables-api, echarts-api, junit, pipeline-utility-steps, theme-manager) |
| jakarta-activation-api | – | 2.1.4-1 (new, for jackson2-api, jakarta-mail-api, jakarta-xml-bind-api) |
| jakarta-mail-api | – | 2.1.5-1 (new, for email-ext, mailer, workflow-basic-steps) |
| jakarta-xml-bind-api | – | 4.0.9-19.v2b_a_5b_44d9a_1c (new, for jackson2-api) |
| javadoc | 217.v905b_86277a_2a_ | 354.vee1a_660b_4990 |
| javax-activation-api | – | 1.2.0-8 (new, for jackson2-api, jaxb, jersey2-api) |
| jaxb | 2.3.0.1 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 1.5 | 83.v417146707a_3d |
| jenkins-design-language | 1.24.8 | 1.27.26 |
| jersey2-api | 2.35-5 | 2.48-180.ve47b_264f849b_ |
| jjwt-api | 0.11.2-9.c8b45b8bb173 | 0.13.0-141.vd58b_a_9592b_6c |
| joda-time-api | – | 2.15.0-208.v6b_084e7a_d9f0 (new, for handy-uri-templates-2-api, pipeline-model-api) |
| jquery3-api | 3.6.0-1 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.1.55.2 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | – | 20260814-226.v20f9685d642c (new, for jackson2-api) |
| json-path-api | – | 3.0.0-218.vcd4dd1355de2 (new, for build-timeout, token-macro) |
| jsoup | 1.10.3 | 1.23.2-106.v1161f282b_eb_f |
| junit | 1.51 | 1442.vb_b_a_467df2767 |
| ldap | 1.26 | 825.v2fca_37dd5b_cb_ |
| lockable-resources | 2.15 | 1576.v217259b_f9976 |
| mailer | 1.34.2 | 534.v1b_36f5864073 |
| mapdb-api | 1.0.9.0 | 1.0.9-44.va_1e1310c9118 |
| matrix-auth | 3.0.1 | 3.3 |
| matrix-project | 1.18.1 | 928.v1effc3f3a_236 |
| maven-plugin | 3.8.1 | 3.28 |
| metrics | 4.0.2.8.1 | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | – | 2.19.0-192.v2b_a_7b_2c1dc71 (new, for git-client, mina-sshd-api-core, sshd) |
| mina-sshd-api-core | – | 2.19.0-192.v2b_a_7b_2c1dc71 (new, for git-client, sshd) |
| momentjs | 1.1.1 | 1.1.1 |
| nodejs | 1.4.3 | 1.6.6 |
| okhttp-api | 4.9.3-105.vb96869f8ac3a | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | – | 477.v7197b_eb_3d8b_8 (new, for javadoc, maven-plugin) |
| pam-auth | 1.6.1 | 1.12 |
| pipeline-build-step | 2.17 | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 1.0 | 65.v203688e7727e |
| pipeline-graph-analysis | 1.11 | 254.v0f63a_a_447dca_ |
| pipeline-groovy-lib | – | 822.vc95f70dea_db_7 (new, for pipeline-github-lib, pipeline-model-definition, workflow-aggregator, workflow-cps-global-lib) |
| pipeline-input-step | 427.va6441fa17010 | 560.v56198a_642157 |
| pipeline-milestone-step | 1.3.2 | 152.v6e22b_8cfc66c |
| pipeline-model-api | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-model-definition | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-rest-api | 2.23 | 2.41 |
| pipeline-stage-step | 291.vf0a8a7aeeb50 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 1.9.3 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.23 | 2.41 |
| pipeline-utility-steps | 2.11.0 | 3.810.va_7672d206740 |
| plain-credentials | 1.8 | 199.v9f8e1f741799 |
| plugin-util-api | 2.3.0 | 7.1403.vd380b_c42e90c |
| popper-api | 1.16.1-2 | 1.16.1-3 |
| popper2-api | 2.5.4-2 | 2.11.6-5 |
| prism-api | – | 1.30.0-741.v034eb_0b_0a_a_fa_ (new, for junit) |
| pubsub-light | 1.13 | 1.19 |
| resource-disposer | 0.19 | 0.25 |
| role-strategy | 483.v17281966f5c3 | 927.v9cf5527c4085 |
| scm-api | 2.6.5 | 728.vc30dcf7a_0df5 |
| script-security | 1138.v8e727069a_025 | 1436.v1082791135e4 |
| snakeyaml-api | 1.30.1 | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | – | 3.2-18.v885ed7ff698f (new, for jackson3-api) |
| ssh-agent | 1.22.1 | 433.v9e73d67c3f74 |
| ssh-credentials | 1.18.1 | 372.va_250881b_08cd |
| ssh-slaves | 1.31.5 | 3.1097.v868116049892 |
| sshd | – | 3.384.vc89b_5e138cf9 (new, for git-server) |
| structs | 308.v852b473a2b8c | 362.va_b_695ef4fdf9 |
| theme-manager | 0.6 | 346.v06cca_64c6a_37 |
| timestamper | 1.16 | 1.30 |
| token-macro | 2.13 | 477.vd4f0dc3cb_cf1 |
| trilead-api | 1.0.13 | 2.284.v1974ea_324382 |
| variant | 1.4 | 70.va_d9f17f859e0 |
| woodstox-core-api | – | 7.3.0-16.v017f30ed004c (new, for jackson2-api, jackson3-api) |
| workflow-aggregator | 2.7 | 608.v67378e9d3db_1 |
| workflow-api | 1138.v619fd5201b_2f | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 2.22 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 2660.2664.v4c114e93f4c1 | 4395.v738456a_5da_6b_ |
| workflow-cps-global-lib | 2.19 | 636.v3d2c06d018a_6 |
| workflow-durable-task-step | 2.40 | 1479.v56e587f413a_7 |
| workflow-job | 1145.v7f2433caa07f | 1609.v8926e7962024 |
| workflow-multibranch | 2.24 | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 2.13 | 487.v41313a_96b_65d |
| workflow-step-api | 622.vb_8e7c15b_c95a_ | 724.v538c2362b_dfb_ |
| workflow-support | 3.8 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.42 | 0.49 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
