# Jenkins upgrade plan: 2.346.3 → 2.580.1

Source: results/2.346.3-top100/plugins-2.346.3.txt · generated 2026-10-21 15:42 UTC by jenkins-upgrade-check dev

**2 core upgrade(s)**, 116 plugins: 101 to upgrade, 13 new dependencies; 24 with security warnings today, 2 not from the update center.

Jenkins 2.346.3 runs on Java 8, 11, 17; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

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

1. **On 2.346.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 17**, upgrade Jenkins to **2.541.3**, start it, then update plugins (114 change(s)).
3. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (14 change(s)).

No Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| ace-editor | 1.1 | 1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| cloudbees-folder | 6.740.ve4f4ffa_dea_54 | 6.1106.v3a_d9a_6d2465e | [SECURITY-3105](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3105) fixed by the upgrade<br>[SECURITY-3106](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3106) fixed by the upgrade<br>[SECURITY-3109](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3109) fixed by the upgrade |
| config-file-provider | 3.11.1 | 1013.v73c323e52b_1f | [SECURITY-3090](https://www.jenkins.io/security/advisory/2023-08-16/#SECURITY-3090) fixed by the upgrade |
| credentials | 1143.vb_e8b_b_ceee347 | 1514.vb_01eb_014e7d4 | [SECURITY-3373](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3373) fixed by the upgrade |
| credentials-binding | 523.525.vb_72269281873 | 728.v902a_273b_8947 | [SECURITY-3499](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3499) fixed by the upgrade<br>[SECURITY-3672](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3672) fixed by the upgrade<br>[SECURITY-3790](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3790) fixed by the upgrade |
| email-ext | 2.92 | 2038.v7b_8817a_499d9 | [SECURITY-2931](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2931) fixed by the upgrade<br>[SECURITY-2934](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2934) fixed by the upgrade<br>[SECURITY-2939](https://www.jenkins.io/security/advisory/2023-02-15/#SECURITY-2939) fixed by the upgrade<br>[SECURITY-3088-1](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(1)) fixed by the upgrade<br>[SECURITY-3088-2](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3088%20(2)) fixed by the upgrade<br>[SECURITY-3705](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3705) fixed by the upgrade |
| git-client | 3.13.1 | 6.6.1 | [SECURITY-3590](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3590) fixed by the upgrade<br>[SECURITY-3614](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3614) fixed by the upgrade<br>[SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| git-server | 1.11 | 137.ve0060b_432302 | [SECURITY-3319](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3319) fixed by the upgrade<br>[SECURITY-3342](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3342) fixed by the upgrade |
| github | 1.34.5 | 1.47.0 | [SECURITY-3246](https://www.jenkins.io/security/advisory/2023-10-25/#SECURITY-3246) fixed by the upgrade<br>[SECURITY-3704](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3704) fixed by the upgrade |
| github-branch-source | 1696.v3a_7603564d04 | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| gradle | 2.12.0.1 | 2.20.1259.v6355c2b_e38c7 | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| htmlpublisher | 1.31 | 429 | [SECURITY-3301](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3301) fixed by the upgrade<br>[SECURITY-3302](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3302) fixed by the upgrade<br>[SECURITY-3303](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3303) fixed by the upgrade<br>[SECURITY-3547](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3547) fixed by the upgrade<br>[SECURITY-3706](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3706) fixed by the upgrade |
| jakarta-mail-api | 2.0.1-2 | 2.1.5-1 | [SECURITY-3617](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3617) fixed by the upgrade |
| ldap | 659.v8ca_b_a_fe79fa_d | 825.v2fca_37dd5b_cb_ | [SECURITY-3046](https://www.jenkins.io/security/advisory/2023-05-16/#SECURITY-3046) fixed by the upgrade<br>[SECURITY-3654](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3654) fixed by the upgrade<br>[SECURITY-3678](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3678) fixed by the upgrade |
| matrix-auth | 3.1.5 | 3.3 | [SECURITY-3676](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3676) fixed by the upgrade |
| matrix-project | 772.v494f19991984 | 928.v1effc3f3a_236 | [SECURITY-3289](https://www.jenkins.io/security/advisory/2024-01-24/#SECURITY-3289) fixed by the upgrade |
| momentjs | 1.1.1 | 1.1.1 | not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself |
| pipeline-build-step | 488.v8993df156e8d | 601.v6d4c6d1a_9dc7 | [SECURITY-3870](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3870) fixed by the upgrade |
| pipeline-groovy-lib | 629.vb_5627b_ee2104 | 822.vc95f70dea_db_7 | [SECURITY-3727](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3727) fixed by the upgrade<br>[SECURITY-3815](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3815) fixed by the upgrade<br>[SECURITY-3796](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3796) fixed by the upgrade |
| pipeline-model-definition | 2.2121.vd87fb_6536d1e | 2.2293.v6e7193cec599 | [SECURITY-3361](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3361) fixed by the upgrade |
| plain-credentials | 143.v1b_df8b_d3b_e48 | 199.v9f8e1f741799 | [SECURITY-2495](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-2495) fixed by the upgrade |
| script-security | 1229.v4880b_b_e905a_6 | 1436.v1082791135e4 | [SECURITY-3341](https://www.jenkins.io/security/advisory/2024-05-02/#SECURITY-3341) fixed by the upgrade<br>[SECURITY-3447](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3447) fixed by the upgrade<br>[SECURITY-3662](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3662) fixed by the upgrade<br>[SECURITY-3792](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3792) fixed by the upgrade<br>[SECURITY-3793](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3793) fixed by the upgrade<br>[SECURITY-3986](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3986) fixed by the upgrade<br>[SECURITY-3897](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3897) fixed by the upgrade<br>[SECURITY-3923](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3923) fixed by the upgrade<br>[SECURITY-3925-1](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(1)) fixed by the upgrade<br>[SECURITY-3925-2](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(2)) fixed by the upgrade<br>[SECURITY-3929](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3929) fixed by the upgrade<br>[SECURITY-3931](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3931) fixed by the upgrade<br>[SECURITY-3932](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3932) fixed by the upgrade<br>[SECURITY-3977](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3977) fixed by the upgrade |
| structs | 324.va_f5d6774f3a_d | 362.va_b_695ef4fdf9 | [SECURITY-3371](https://www.jenkins.io/security/advisory/2024-06-26/#SECURITY-3371) fixed by the upgrade |
| trilead-api | 1.67.vc3938a_35172f | 2.284.v1974ea_324382 | [SECURITY-3333](https://www.jenkins.io/security/advisory/2024-03-06/#SECURITY-3333) fixed by the upgrade |
| workflow-cps | 3536.vb_8a_6628079d5 | 4395.v738456a_5da_6b_ | [SECURITY-3362](https://www.jenkins.io/security/advisory/2024-11-13/#SECURITY-3362) fixed by the upgrade<br>[SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-multibranch | 716.vc692a_e52371b_ | 842.v3a_b_59b_57b_e6e | [SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| ace-editor | 1.1 | 1.1 |
| ant | 481.v7b_09e538fcca | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 159.v25b_c67cd35fb_ | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.13-138.v4e7d9a_7b_a_e61 | 4.5.14-269.vfa_2321039a_83 |
| asm-api | – | 9.11.1-226.vff833a_e4c58e (new, for json-path-api, scm-api) |
| authentication-tokens | 1.4 | 1.144.v5ff4a_5ec5c33 |
| bootstrap5-api | 5.2.1-3 | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.27 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1071.v1a_188a_562481 | 2.1303.v9f3b_95dc329d |
| build-timeout | 1.24 | 1.41 |
| caffeine-api | 2.9.3-65.v6a_47d0f4d1fe | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 1.8.1 | 422.vc7e8b_51488c6 |
| cloudbees-folder | 6.740.ve4f4ffa_dea_54 | 6.1106.v3a_d9a_6d2465e |
| command-launcher | 90.v669d7ccb_7c31 | 141.v9578cf36b_311 |
| commons-lang3-api | 3.12.0-36.vd97de6465d5b_ | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.10.0-27.vb_fa_3896786a_7 | 1.15.0-218.va_61573470393 |
| config-file-provider | 3.11.1 | 1013.v73c323e52b_1f |
| credentials | 1143.vb_e8b_b_ceee347 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 523.525.vb_72269281873 | 728.v902a_273b_8947 |
| dark-theme | 262.v0202a_4c8fb_6a | 652.vea_da_dfea_e769 |
| display-url-api | 2.3.6 | 2.217.va_6b_de84cc74b_ |
| docker-commons | 1.21 | 477.v289085a_b_6896 |
| docker-workflow | 563.vd5d2e5c4007f | 653.v2f2c08eff0ec |
| durable-task | 507.v050055d0cb_dd | 686.v80ff80875b_82 |
| echarts-api | 5.4.0-1 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | – | 0.3.0.1-29.v67e9a_1c969b_b_ (new, for mina-sshd-api-common, trilead-api) |
| email-ext | 2.92 | 2038.v7b_8817a_499d9 |
| font-awesome-api | 6.2.1-1 | 7.3.1-1013.v0835a_879ec6d |
| git | 4.14.3 | 5.11.0 |
| git-client | 3.13.1 | 6.6.1 |
| git-server | 1.11 | 137.ve0060b_432302 |
| github | 1.34.5 | 1.47.0 |
| github-api | 1.303-400.v35c2d8258028 | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 1696.v3a_7603564d04 | 1986.v801e40614d00 |
| gradle | 2.12.0.1 | 2.20.1259.v6355c2b_e38c7 |
| gson-api | – | 2.14.0-201.v8eefe5515533 (new, for git-client, trilead-api) |
| htmlpublisher | 1.31 | 429 |
| instance-identity | – | 203.v15e81a_1b_7a_38 (new, for git-server, github, mailer, sshd) |
| ionicons-api | 31.v4757b_6987003 | 94.vcc3065403257 |
| jackson-annotations2-api | – | 2.22-19.v10a_a_582ea_26e (new, for jackson2-api, jackson3-api) |
| jackson2-api | 2.14.2-319.v37853346a_229 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | – | 3.2.4-103.vf5a_b_3f4ea_cec (new, for echarts-api, junit, pipeline-graph-view, theme-manager) |
| jakarta-activation-api | 2.0.1-2 | 2.1.4-1 |
| jakarta-mail-api | 2.0.1-2 | 2.1.5-1 |
| jakarta-xml-bind-api | – | 4.0.9-19.v2b_a_5b_44d9a_1c (new, for jackson2-api) |
| javadoc | 226.v71211feb_e7e9 | 354.vee1a_660b_4990 |
| javax-activation-api | 1.2.0-5 | 1.2.0-8 |
| javax-mail-api | 1.6.2-8 | 1.6.2-11 |
| jaxb | 2.3.7-1 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 63.v62d2fd4b_4793 | 83.v417146707a_3d |
| jjwt-api | 0.11.5-77.v646c772fddb_0 | 0.13.0-141.vd58b_a_9592b_6c |
| joda-time-api | – | 2.15.0-208.v6b_084e7a_d9f0 (new, for pipeline-model-api) |
| jquery3-api | 3.6.1-2 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.2.8-65.v052c39de79b_2 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | – | 20260814-226.v20f9685d642c (new, for jackson2-api) |
| json-path-api | – | 3.0.0-218.vcd4dd1355de2 (new, for build-timeout, token-macro) |
| jsoup | 1.10.3 | 1.23.2-106.v1161f282b_eb_f |
| junit | 1166.1168.vd6b_8042a_06de | 1442.vb_b_a_467df2767 |
| ldap | 659.v8ca_b_a_fe79fa_d | 825.v2fca_37dd5b_cb_ |
| lockable-resources | 2.18 | 1576.v217259b_f9976 |
| mailer | 435.438.v5b_81173f5b_a_1 | 534.v1b_36f5864073 |
| matrix-auth | 3.1.5 | 3.3 |
| matrix-project | 772.v494f19991984 | 928.v1effc3f3a_236 |
| maven-plugin | 3.22 | 3.28 |
| metrics | 4.2.13-420.vea_2f17932dd6 | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.9.2-50.va_0e1f42659a_a | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.9.2-50.va_0e1f42659a_a | 2.19.0-192.v2b_a_7b_2c1dc71 |
| momentjs | 1.1.1 | 1.1.1 |
| okhttp-api | 4.9.3-108.v0feda04578cf | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | – | 477.v7197b_eb_3d8b_8 (new, for javadoc, maven-plugin) |
| pam-auth | 1.10 | 1.12 |
| pipeline-build-step | 488.v8993df156e8d | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 38.v445716ea_edda_ | 65.v203688e7727e |
| pipeline-graph-analysis | 202.va_d268e64deb_3 | 254.v0f63a_a_447dca_ |
| pipeline-graph-view | 144.vf3924feb_7e35 | 1053.v9df19e6ceed2 |
| pipeline-groovy-lib | 629.vb_5627b_ee2104 | 822.vc95f70dea_db_7 |
| pipeline-input-step | 466.v6d0a_5df34f81 | 560.v56198a_642157 |
| pipeline-milestone-step | 111.v449306f708b_7 | 152.v6e22b_8cfc66c |
| pipeline-model-api | 2.2121.vd87fb_6536d1e | 2.2293.v6e7193cec599 |
| pipeline-model-definition | 2.2121.vd87fb_6536d1e | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 2.2121.vd87fb_6536d1e | 2.2293.v6e7193cec599 |
| pipeline-rest-api | 2.28 | 2.41 |
| pipeline-stage-step | 305.ve96d0205c1c6 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 2.2121.vd87fb_6536d1e | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.28 | 2.41 |
| plain-credentials | 143.v1b_df8b_d3b_e48 | 199.v9f8e1f741799 |
| plugin-util-api | 2.20.0 | 7.1403.vd380b_c42e90c |
| popper2-api | 2.11.6-2 | 2.11.6-5 |
| prism-api | 1.29.0-2 | 1.30.0-741.v034eb_0b_0a_a_fa_ |
| resource-disposer | 0.20 | 0.25 |
| role-strategy | 587.588.v850a_20a_30162 | 927.v9cf5527c4085 |
| scm-api | 631.v9143df5b_e4a_a | 728.vc30dcf7a_0df5 |
| script-security | 1229.v4880b_b_e905a_6 | 1436.v1082791135e4 |
| snakeyaml-api | 1.33-90.v80dcb_3814d35 | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | – | 3.2-18.v885ed7ff698f (new, for jackson3-api) |
| ssh-credentials | 305.v8f4381501156 | 372.va_250881b_08cd |
| ssh-slaves | 1.834.v622da_57f702c | 3.1097.v868116049892 |
| sshd | 3.242.va_db_9da_b_26a_c3 | 3.384.vc89b_5e138cf9 |
| structs | 324.va_f5d6774f3a_d | 362.va_b_695ef4fdf9 |
| theme-manager | 1.6 | 346.v06cca_64c6a_37 |
| timestamper | 1.21 | 1.30 |
| token-macro | 321.vd7cc1f2a_52c8 | 477.vd4f0dc3cb_cf1 |
| trilead-api | 1.67.vc3938a_35172f | 2.284.v1974ea_324382 |
| variant | 59.vf075fe829ccb | 70.va_d9f17f859e0 |
| woodstox-core-api | – | 7.3.0-16.v017f30ed004c (new, for jackson2-api, jackson3-api) |
| workflow-aggregator | 596.v8c21c963d92d | 608.v67378e9d3db_1 |
| workflow-api | 1200.v8005c684b_a_c6 | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 994.vd57e3ca_46d24 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 3536.vb_8a_6628079d5 | 4395.v738456a_5da_6b_ |
| workflow-durable-task-step | 1217.v38306d8fa_b_5c | 1479.v56e587f413a_7 |
| workflow-job | 1207.1209.v69351208a_5a_7 | 1609.v8926e7962024 |
| workflow-multibranch | 716.vc692a_e52371b_ | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 408.v7d5b_135a_b_d49 | 487.v41313a_96b_65d |
| workflow-step-api | 639.v6eca_cd8c04a_a_ | 724.v538c2362b_dfb_ |
| workflow-support | 839.v35e2736cfd5c | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.43 | 0.49 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
