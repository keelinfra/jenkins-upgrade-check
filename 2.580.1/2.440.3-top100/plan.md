# Jenkins upgrade plan: 2.440.3 → 2.580.1

Source: results/2.440.3-top100/plugins-2.440.3.txt · generated 2026-10-21 12:48 UTC by jenkins-upgrade-check dev

**2 core upgrade(s)**, 106 plugins: 100 to upgrade, 6 new dependencies; 16 with security warnings today.

Jenkins 2.440.3 runs on Java 11, 17; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

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

1. **On 2.440.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 17**, upgrade Jenkins to **2.541.3**, start it, then update plugins (106 change(s)).
3. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (12 change(s)).

No Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| credentials | 1380.va_435002fa_924 | 1514.vb_01eb_014e7d4 | [SECURITY-3373](https://www.jenkins.io/security/advisory/2024-10-02/#SECURITY-3373) fixed by the upgrade |
| credentials-binding | 681.vf91669a_32e45 | 728.v902a_273b_8947 | [SECURITY-3499](https://www.jenkins.io/security/advisory/2025-07-09/#SECURITY-3499) fixed by the upgrade<br>[SECURITY-3672](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3672) fixed by the upgrade<br>[SECURITY-3790](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3790) fixed by the upgrade |
| eddsa-api | 0.3.0-4.v84c6f0f4969e | 0.3.0.1-29.v67e9a_1c969b_b_ | [SECURITY-3404](https://www.jenkins.io/security/advisory/2025-03-19/#SECURITY-3404) fixed by the upgrade |
| email-ext | 1844.v3ea_a_b_842374a_ | 2038.v7b_8817a_499d9 | [SECURITY-3705](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3705) fixed by the upgrade |
| git-client | 5.0.3 | 6.6.1 | [SECURITY-3590](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3590) fixed by the upgrade<br>[SECURITY-3614](https://www.jenkins.io/security/advisory/2025-12-10/#SECURITY-3614) fixed by the upgrade<br>[SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| github | 1.40.0 | 1.47.0 | [SECURITY-3704](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3704) fixed by the upgrade |
| github-branch-source | 1807.v50351eb_7dd13 | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| gradle | 2.15 | 2.20.1253.vc116f0763a_eb_ | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| jakarta-mail-api | 2.1.3-1 | 2.1.5-1 | [SECURITY-3617](https://www.jenkins.io/security/advisory/2025-09-03/#SECURITY-3617) fixed by the upgrade |
| ldap | 725.v3cb_b_711b_1a_ef | 825.v2fca_37dd5b_cb_ | [SECURITY-3654](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3654) fixed by the upgrade<br>[SECURITY-3678](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3678) fixed by the upgrade |
| matrix-auth | 3.2.3 | 3.3 | [SECURITY-3676](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3676) fixed by the upgrade |
| pipeline-build-step | 540.vb_e8849e1a_b_d8 | 601.v6d4c6d1a_9dc7 | [SECURITY-3870](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3870) fixed by the upgrade |
| pipeline-groovy-lib | 745.vdf6077913de0 | 822.vc95f70dea_db_7 | [SECURITY-3727](https://www.jenkins.io/security/advisory/2026-05-27/#SECURITY-3727) fixed by the upgrade<br>[SECURITY-3815](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3815) fixed by the upgrade<br>[SECURITY-3796](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3796) fixed by the upgrade |
| script-security | 1369.v9b_98a_4e95b_2d | 1436.v1082791135e4 | [SECURITY-3662](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3662) fixed by the upgrade<br>[SECURITY-3792](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3792) fixed by the upgrade<br>[SECURITY-3793](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3793) fixed by the upgrade<br>[SECURITY-3986](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3986) fixed by the upgrade<br>[SECURITY-3897](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3897) fixed by the upgrade<br>[SECURITY-3923](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3923) fixed by the upgrade<br>[SECURITY-3925-1](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(1)) fixed by the upgrade<br>[SECURITY-3925-2](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3925%20(2)) fixed by the upgrade<br>[SECURITY-3929](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3929) fixed by the upgrade<br>[SECURITY-3931](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3931) fixed by the upgrade<br>[SECURITY-3932](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3932) fixed by the upgrade<br>[SECURITY-3977](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3977) fixed by the upgrade |
| workflow-cps | 4000.v5198556e9cea_ | 4395.v738456a_5da_6b_ | [SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-multibranch | 773.vc4fe1378f1d5 | 842.v3a_b_59b_57b_e6e | [SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| ant | 511.v0a_a_1a_334f41b_ | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 162.v0e6ec0fcfcf6 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.14-208.v438351942757 | 4.5.14-269.vfa_2321039a_83 |
| asm-api | 9.7.1-95.v9f552033802a_ | 9.11.1-226.vff833a_e4c58e |
| authentication-tokens | 1.119.v50285141b_7e1 | 1.144.v5ff4a_5ec5c33 |
| bootstrap5-api | 5.3.3-1 | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.30.1.78.1-248.ve27176eb_46cb_ | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1202.vc51412f3ed84 | 2.1303.v9f3b_95dc329d |
| build-timeout | 1.33 | 1.41 |
| caffeine-api | 3.2.0-161.v691ef352cee1 | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 2.2.3 | 422.vc7e8b_51488c6 |
| cloudbees-folder | 6.901.vb_4c7a_da_75da_3 | 6.1106.v3a_d9a_6d2465e |
| command-launcher | 115.vd8b_301cc15d0 | 141.v9578cf36b_311 |
| commons-lang3-api | 3.17.0-84.vb_b_938040b_078 | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.12.0-129.v99a_50df237f7 | 1.15.0-218.va_61573470393 |
| credentials | 1380.va_435002fa_924 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 681.vf91669a_32e45 | 728.v902a_273b_8947 |
| dark-theme | 439.vdef09f81f85e | 652.vea_da_dfea_e769 |
| display-url-api | 2.204.vf6fddd8a_8b_e9 | 2.217.va_6b_de84cc74b_ |
| docker-commons | 445.v6b_646c962a_94 | 477.v289085a_b_6896 |
| durable-task | 581.v299a_5609d767 | 686.v80ff80875b_82 |
| echarts-api | 5.5.1-4 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | 0.3.0-4.v84c6f0f4969e | 0.3.0.1-29.v67e9a_1c969b_b_ |
| email-ext | 1844.v3ea_a_b_842374a_ | 2038.v7b_8817a_499d9 |
| font-awesome-api | 6.6.0-2 | 7.3.1-1013.v0835a_879ec6d |
| git | 5.5.2 | 5.11.0 |
| git-client | 5.0.3 | 6.6.1 |
| github | 1.40.0 | 1.47.0 |
| github-api | 1.321-468.v6a_9f5f2d5a_7e | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 1807.v50351eb_7dd13 | 1986.v801e40614d00 |
| gradle | 2.15 | 2.20.1253.vc116f0763a_eb_ |
| gson-api | 2.11.0-41.v019fcf6125dc | 2.14.0-201.v8eefe5515533 |
| instance-identity | 201.vd2a_b_5a_468a_a_6 | 203.v15e81a_1b_7a_38 |
| ionicons-api | 74.v93d5eb_813d5f | 94.vcc3065403257 |
| jackson-annotations2-api | – | 2.22-19.v10a_a_582ea_26e (new, for jackson2-api, jackson3-api) |
| jackson2-api | 2.18.3-396.v93cc2d00b_b_7b_ | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | – | 3.2.4-103.vf5a_b_3f4ea_cec (new, for echarts-api, junit, pipeline-graph-view, theme-manager) |
| jakarta-activation-api | 2.1.3-1 | 2.1.4-1 |
| jakarta-mail-api | 2.1.3-1 | 2.1.5-1 |
| jakarta-xml-bind-api | – | 4.0.9-19.v2b_a_5b_44d9a_1c (new, for jackson2-api) |
| javadoc | 280.v050b_5c849f69 | 354.vee1a_660b_4990 |
| javax-activation-api | 1.2.0-7 | 1.2.0-8 |
| javax-mail-api | 1.6.2-10 | 1.6.2-11 |
| jaxb | 2.3.9-1 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 80.v8a_dee33ed6f0 | 83.v417146707a_3d |
| jjwt-api | 0.11.5-112.ve82dfb_224b_a_d | 0.13.0-141.vd58b_a_9592b_6c |
| joda-time-api | 2.13.0-85.vb_64d1c2921f1 | 2.15.0-208.v6b_084e7a_d9f0 |
| jquery3-api | 3.7.1-2 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.2.16-86.v42e010d9484b_ | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | 20240303-41.v94e11e6de726 | 20260814-226.v20f9685d642c |
| json-path-api | 2.9.0-58.v62e3e85b_a_655 | 3.0.0-218.vcd4dd1355de2 |
| jsoup | 1.10.3 | 1.23.2-106.v1161f282b_eb_f |
| junit | 1312.v1a_235a_b_94a_31 | 1442.vb_b_a_467df2767 |
| ldap | 725.v3cb_b_711b_1a_ef | 825.v2fca_37dd5b_cb_ |
| mailer | 488.v0c9639c1a_eb_3 | 534.v1b_36f5864073 |
| matrix-auth | 3.2.3 | 3.3 |
| matrix-project | 839.vff91cd7e3a_b_2 | 928.v1effc3f3a_236 |
| metrics | 4.2.21-458.vcf496cb_839e4 | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.14.0-143.v2b_362fc39576 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.14.0-143.v2b_362fc39576 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| okhttp-api | 4.11.0-172.vda_da_1feeb_c6e | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | – | 477.v7197b_eb_3d8b_8 (new, for javadoc) |
| pam-auth | 1.11 | 1.12 |
| pipeline-build-step | 540.vb_e8849e1a_b_d8 | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 61.v629f2cc41d83 | 65.v203688e7727e |
| pipeline-graph-analysis | 216.vfd8b_ece330ca_ | 254.v0f63a_a_447dca_ |
| pipeline-graph-view | 243.vc9e11fec486a_ | 1053.v9df19e6ceed2 |
| pipeline-groovy-lib | 745.vdf6077913de0 | 822.vc95f70dea_db_7 |
| pipeline-input-step | 508.v584c0e9a_2177 | 560.v56198a_642157 |
| pipeline-milestone-step | 119.vdfdc43fc3b_9a_ | 152.v6e22b_8cfc66c |
| pipeline-model-api | 2.2218.v56d0cda_37c72 | 2.2293.v6e7193cec599 |
| pipeline-model-definition | 2.2218.v56d0cda_37c72 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 2.2218.v56d0cda_37c72 | 2.2293.v6e7193cec599 |
| pipeline-rest-api | 2.34 | 2.41 |
| pipeline-stage-step | 312.v8cd10304c27a_ | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 2.2218.v56d0cda_37c72 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.34 | 2.41 |
| plain-credentials | 183.va_de8f1dd5a_2b_ | 199.v9f8e1f741799 |
| plugin-util-api | 4.1.0 | 7.1403.vd380b_c42e90c |
| prism-api | 1.29.0-18 | 1.30.0-741.v034eb_0b_0a_a_fa_ |
| resource-disposer | 0.24 | 0.25 |
| scm-api | 696.v778d637b_a_762 | 728.vc30dcf7a_0df5 |
| script-security | 1369.v9b_98a_4e95b_2d | 1436.v1082791135e4 |
| snakeyaml-api | 2.3-123.v13484c65210a_ | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | – | 3.2-18.v885ed7ff698f (new, for jackson3-api) |
| ssh-credentials | 349.vb_8b_6b_9709f5b_ | 372.va_250881b_08cd |
| ssh-slaves | 2.973.v0fa_8c0dea_f9f | 3.1097.v868116049892 |
| sshd | 3.330.vc866a_8389b_58 | 3.384.vc89b_5e138cf9 |
| structs | 338.v848422169819 | 362.va_b_695ef4fdf9 |
| theme-manager | 262.vc57ee4a_eda_5d | 346.v06cca_64c6a_37 |
| timestamper | 1.27 | 1.30 |
| token-macro | 400.v35420b_922dcb_ | 477.vd4f0dc3cb_cf1 |
| trilead-api | 2.147.vb_73cc728a_32e | 2.284.v1974ea_324382 |
| variant | 60.v7290fc0eb_b_cd | 70.va_d9f17f859e0 |
| woodstox-core-api | – | 7.3.0-16.v017f30ed004c (new, for jackson2-api, jackson3-api) |
| workflow-aggregator | 600.vb_57cdd26fdd7 | 608.v67378e9d3db_1 |
| workflow-api | 1336.vee415d95c521 | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 1058.vcb_fc1e3a_21a_9 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 4000.v5198556e9cea_ | 4395.v738456a_5da_6b_ |
| workflow-durable-task-step | 1378.v6a_3e903058a_3 | 1479.v56e587f413a_7 |
| workflow-job | 1400.v7fd111b_ec82f | 1609.v8926e7962024 |
| workflow-multibranch | 773.vc4fe1378f1d5 | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 427.v4ca_6512e7df1 | 487.v41313a_96b_65d |
| workflow-step-api | 678.v3ee58b_469476 | 724.v538c2362b_dfb_ |
| workflow-support | 936.v9fa_77211ca_e1 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.47 | 0.49 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
