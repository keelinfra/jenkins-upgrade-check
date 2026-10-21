# Jenkins upgrade plan: 2.504.3 → 2.580.1

Source: results/2.504.3-top100/plugins-2.504.3.txt · generated 2026-10-21 15:41 UTC by jenkins-upgrade-check dev

**1 core upgrade(s)**, 105 plugins: 25 to upgrade, 0 new dependencies; 6 with security warnings today.

Jenkins 2.504.3 runs on Java 17, 21; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

- [core-2_528](https://www.jenkins.io/security/advisory/2025-09-17/): Multiple security vulnerabilities in Jenkins 2.527 and earlier, LTS 2.516.2 and earlier
- [core-2_541](https://www.jenkins.io/security/advisory/2025-12-10/): Multiple security vulnerabilities in Jenkins 2.540 and earlier, LTS 2.528.2 and earlier
- [core-2_551](https://www.jenkins.io/security/advisory/2026-02-18/): Multiple security vulnerabilities in Jenkins 2.550 and earlier, LTS 2.541.1 and earlier
- [core-2_555](https://www.jenkins.io/security/advisory/2026-03-18/): Multiple security vulnerabilities in Jenkins 2.554 and earlier, LTS 2.541.2 and earlier
- [core-2_568](https://www.jenkins.io/security/advisory/2026-06-10/): Multiple security vulnerabilities in Jenkins 2.567 and earlier, LTS 2.555.2 and earlier
- [core-2_576](https://www.jenkins.io/security/advisory/2026-08-05/): Multiple security vulnerabilities in Jenkins 2.575 and earlier, LTS 2.568.1 and earlier
- [core-2_580](https://www.jenkins.io/security/advisory/2026-09-02/): Multiple security vulnerabilities in Jenkins 2.579 and earlier, LTS 2.568.2 and earlier

## Steps

1. **On 2.504.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (25 change(s)).

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| git-client | 6.5.0 | 6.6.1 | [SECURITY-3723](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3723) fixed by the upgrade |
| github-branch-source | 1917.v9ee8a_39b_3d0d | 1986.v801e40614d00 | [SECURITY-3702](https://www.jenkins.io/security/advisory/2026-04-29/#SECURITY-3702) fixed by the upgrade<br>[SECURITY-3808](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3808) fixed by the upgrade |
| gradle | 2.18.1203.v2c96b_1243c72 | 2.20.1259.v6355c2b_e38c7 | [SECURITY-4028](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4028) fixed by the upgrade |
| pipeline-groovy-lib | 798.v5cc688825312 | 822.vc95f70dea_db_7 | [SECURITY-3815](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3815) fixed by the upgrade<br>[SECURITY-3796](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3796) fixed by the upgrade |
| workflow-cps | 4303.v8fa_2a_581f0a_6 | 4395.v738456a_5da_6b_ | [SECURITY-3677](https://www.jenkins.io/security/advisory/2026-06-24/#SECURITY-3677) fixed by the upgrade |
| workflow-multibranch | 841.vec5b_9e1806ec | 842.v3a_b_59b_57b_e6e | [SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| ant | 520.vd082ecfb_16a_9 | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 173.v680e3a_b_69ff3 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.14-269.vfa_2321039a_83 | 4.5.14-269.vfa_2321039a_83 |
| asm-api | 9.11.1-226.vff833a_e4c58e | 9.11.1-226.vff833a_e4c58e |
| authentication-tokens | 1.144.v5ff4a_5ec5c33 | 1.144.v5ff4a_5ec5c33 |
| bootstrap5-api | 5.3.8-890.v1c5cf4fa_178e | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.30.1.86.2-310.v4655f6742752 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1303.v9f3b_95dc329d | 2.1303.v9f3b_95dc329d |
| build-timeout | 1.38 | 1.41 |
| caffeine-api | 3.3.1-214.v18eb_56b_34920 | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 422.vc7e8b_51488c6 | 422.vc7e8b_51488c6 |
| cloudbees-folder | 6.1053.vd62fb_b_f7367b_ | 6.1106.v3a_d9a_6d2465e |
| command-launcher | 141.v9578cf36b_311 | 141.v9578cf36b_311 |
| commons-lang3-api | 3.21.0-122.v6a_ea_4fd34049 | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.15.0-218.va_61573470393 | 1.15.0-218.va_61573470393 |
| credentials | 1487.va_d001ede2c31 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 728.v902a_273b_8947 | 728.v902a_273b_8947 |
| dark-theme | 524.vd675b_22b_30cb_ | 652.vea_da_dfea_e769 |
| display-url-api | 2.217.va_6b_de84cc74b_ | 2.217.va_6b_de84cc74b_ |
| docker-commons | 457.v0f62a_94f11a_3 | 477.v289085a_b_6896 |
| durable-task | 686.v80ff80875b_82 | 686.v80ff80875b_82 |
| echarts-api | 6.0.0-1146.v5c8f3b_8f0573 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | 0.3.0.1-29.v67e9a_1c969b_b_ | 0.3.0.1-29.v67e9a_1c969b_b_ |
| email-ext | 1933.1935.v276319e3cc47 | 2038.v7b_8817a_499d9 |
| font-awesome-api | 7.0.1-859.v128d3a_efb_6e5 | 7.3.1-1013.v0835a_879ec6d |
| git | 5.11.0 | 5.11.0 |
| git-client | 6.5.0 | 6.6.1 |
| github | 1.47.0 | 1.47.0 |
| github-api | 1.330-492.v3941a_032db_2a_ | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 1917.v9ee8a_39b_3d0d | 1986.v801e40614d00 |
| gradle | 2.18.1203.v2c96b_1243c72 | 2.20.1259.v6355c2b_e38c7 |
| gson-api | 2.14.0-201.v8eefe5515533 | 2.14.0-201.v8eefe5515533 |
| instance-identity | 203.v15e81a_1b_7a_38 | 203.v15e81a_1b_7a_38 |
| ionicons-api | 94.vcc3065403257 | 94.vcc3065403257 |
| jackson-annotations2-api | 2.22-19.v10a_a_582ea_26e | 2.22-19.v10a_a_582ea_26e |
| jackson2-api | 2.21.2-436.v29efdb_7418ff | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | 3.2.0-83.vfb_04cee1ea_0f | 3.2.4-103.vf5a_b_3f4ea_cec |
| jakarta-activation-api | 2.1.4-1 | 2.1.4-1 |
| jakarta-mail-api | 2.1.5-1 | 2.1.5-1 |
| jakarta-xml-bind-api | 4.0.9-19.v2b_a_5b_44d9a_1c | 4.0.9-19.v2b_a_5b_44d9a_1c |
| javadoc | 354.vee1a_660b_4990 | 354.vee1a_660b_4990 |
| javax-activation-api | 1.2.0-8 | 1.2.0-8 |
| javax-mail-api | 1.6.2-11 | 1.6.2-11 |
| jaxb | 2.3.9-143.v5979df3304e6 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 83.v417146707a_3d | 83.v417146707a_3d |
| jjwt-api | 0.13.0-141.vd58b_a_9592b_6c | 0.13.0-141.vd58b_a_9592b_6c |
| joda-time-api | 2.15.0-208.v6b_084e7a_d9f0 | 2.15.0-208.v6b_084e7a_d9f0 |
| jquery3-api | 3.7.1-594.vb_3864f326cf0 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.2.16-95.v3eecb_55fa_b_78 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | 20260814-226.v20f9685d642c | 20260814-226.v20f9685d642c |
| json-path-api | 3.0.0-218.vcd4dd1355de2 | 3.0.0-218.vcd4dd1355de2 |
| jsoup | 1.23.2-106.v1161f282b_eb_f | 1.23.2-106.v1161f282b_eb_f |
| junit | 1369.v15da_00283f06 | 1442.vb_b_a_467df2767 |
| ldap | 825.v2fca_37dd5b_cb_ | 825.v2fca_37dd5b_cb_ |
| mailer | 534.v1b_36f5864073 | 534.v1b_36f5864073 |
| matrix-auth | 3.2.10 | 3.3 |
| matrix-project | 870.v9db_fcfc2f45b_ | 928.v1effc3f3a_236 |
| metrics | 4.2.37-494.v06f9a_939d33a_ | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| okhttp-api | 5.5.0-242.v1482d48d3ea_3 | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | 477.v7197b_eb_3d8b_8 | 477.v7197b_eb_3d8b_8 |
| pam-auth | 1.12 | 1.12 |
| pipeline-build-step | 601.v6d4c6d1a_9dc7 | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 65.v203688e7727e | 65.v203688e7727e |
| pipeline-graph-analysis | 254.v0f63a_a_447dca_ | 254.v0f63a_a_447dca_ |
| pipeline-groovy-lib | 798.v5cc688825312 | 822.vc95f70dea_db_7 |
| pipeline-input-step | 560.v56198a_642157 | 560.v56198a_642157 |
| pipeline-milestone-step | 152.v6e22b_8cfc66c | 152.v6e22b_8cfc66c |
| pipeline-model-api | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-model-definition | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-rest-api | 2.41 | 2.41 |
| pipeline-stage-step | 345.va_96187909426 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.41 | 2.41 |
| plain-credentials | 199.v9f8e1f741799 | 199.v9f8e1f741799 |
| plugin-util-api | 6.1167.v022176c7e0ca_ | 7.1403.vd380b_c42e90c |
| prism-api | 1.30.0-609.vf0a_df102d9a_f | 1.30.0-741.v034eb_0b_0a_a_fa_ |
| resource-disposer | 0.25 | 0.25 |
| scm-api | 728.vc30dcf7a_0df5 | 728.vc30dcf7a_0df5 |
| script-security | 1436.v1082791135e4 | 1436.v1082791135e4 |
| snakeyaml-api | 2.7-152.vc672728cf536 | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | 3.2-18.v885ed7ff698f | 3.2-18.v885ed7ff698f |
| ssh-credentials | 372.va_250881b_08cd | 372.va_250881b_08cd |
| ssh-slaves | 3.1097.v868116049892 | 3.1097.v868116049892 |
| sshd | 3.384.vc89b_5e138cf9 | 3.384.vc89b_5e138cf9 |
| structs | 362.va_b_695ef4fdf9 | 362.va_b_695ef4fdf9 |
| theme-manager | 346.v06cca_64c6a_37 | 346.v06cca_64c6a_37 |
| timestamper | 1.30 | 1.30 |
| token-macro | 477.vd4f0dc3cb_cf1 | 477.vd4f0dc3cb_cf1 |
| trilead-api | 2.284.v1974ea_324382 | 2.284.v1974ea_324382 |
| variant | 70.va_d9f17f859e0 | 70.va_d9f17f859e0 |
| woodstox-core-api | 7.3.0-16.v017f30ed004c | 7.3.0-16.v017f30ed004c |
| workflow-aggregator | 608.v67378e9d3db_1 | 608.v67378e9d3db_1 |
| workflow-api | 1384.vdc05a_48f535f | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 1098.v808b_fd7f8cf4 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 4303.v8fa_2a_581f0a_6 | 4395.v738456a_5da_6b_ |
| workflow-durable-task-step | 1479.v56e587f413a_7 | 1479.v56e587f413a_7 |
| workflow-job | 1571.vb_423c255d6d9 | 1609.v8926e7962024 |
| workflow-multibranch | 841.vec5b_9e1806ec | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 487.v41313a_96b_65d | 487.v41313a_96b_65d |
| workflow-step-api | 724.v538c2362b_dfb_ | 724.v538c2362b_dfb_ |
| workflow-support | 1015.v785e5a_b_b_8b_22 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.49 | 0.49 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
