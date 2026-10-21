# Jenkins upgrade plan: 2.568.3 → 2.580.1

Source: results/2.568.3-top100/plugins-2.568.3.txt · generated 2026-10-21 12:57 UTC by jenkins-upgrade-check dev

**1 core upgrade(s)**, 105 plugins: 0 to upgrade, 0 new dependencies.

Jenkins 2.568.3 runs on Java 21, 25; Jenkins 2.580.1 runs on Java 21, 25.

## Steps

1. **On 2.568.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (0 change(s)).

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| ant | 520.vd082ecfb_16a_9 | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 173.v680e3a_b_69ff3 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.14-269.vfa_2321039a_83 | 4.5.14-269.vfa_2321039a_83 |
| asm-api | 9.11.1-226.vff833a_e4c58e | 9.11.1-226.vff833a_e4c58e |
| authentication-tokens | 1.144.v5ff4a_5ec5c33 | 1.144.v5ff4a_5ec5c33 |
| bootstrap5-api | 5.3.8-1048.va_c299057e35c | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.30.1.86.2-310.v4655f6742752 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1303.v9f3b_95dc329d | 2.1303.v9f3b_95dc329d |
| build-timeout | 1.41 | 1.41 |
| caffeine-api | 3.3.1-214.v18eb_56b_34920 | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 422.vc7e8b_51488c6 | 422.vc7e8b_51488c6 |
| cloudbees-folder | 6.1106.v3a_d9a_6d2465e | 6.1106.v3a_d9a_6d2465e |
| command-launcher | 141.v9578cf36b_311 | 141.v9578cf36b_311 |
| commons-lang3-api | 3.21.0-122.v6a_ea_4fd34049 | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.15.0-218.va_61573470393 | 1.15.0-218.va_61573470393 |
| credentials | 1514.vb_01eb_014e7d4 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 728.v902a_273b_8947 | 728.v902a_273b_8947 |
| dark-theme | 652.vea_da_dfea_e769 | 652.vea_da_dfea_e769 |
| display-url-api | 2.217.va_6b_de84cc74b_ | 2.217.va_6b_de84cc74b_ |
| docker-commons | 477.v289085a_b_6896 | 477.v289085a_b_6896 |
| durable-task | 686.v80ff80875b_82 | 686.v80ff80875b_82 |
| echarts-api | 6.1.0-1321.v86c4c5114494 | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | 0.3.0.1-29.v67e9a_1c969b_b_ | 0.3.0.1-29.v67e9a_1c969b_b_ |
| email-ext | 2038.v7b_8817a_499d9 | 2038.v7b_8817a_499d9 |
| font-awesome-api | 7.3.1-1013.v0835a_879ec6d | 7.3.1-1013.v0835a_879ec6d |
| git | 5.11.0 | 5.11.0 |
| git-client | 6.6.1 | 6.6.1 |
| github | 1.47.0 | 1.47.0 |
| github-api | 1.330-492.v3941a_032db_2a_ | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 1986.v801e40614d00 | 1986.v801e40614d00 |
| gradle | 2.20.1253.vc116f0763a_eb_ | 2.20.1253.vc116f0763a_eb_ |
| gson-api | 2.14.0-201.v8eefe5515533 | 2.14.0-201.v8eefe5515533 |
| instance-identity | 203.v15e81a_1b_7a_38 | 203.v15e81a_1b_7a_38 |
| ionicons-api | 94.vcc3065403257 | 94.vcc3065403257 |
| jackson-annotations2-api | 2.22-19.v10a_a_582ea_26e | 2.22-19.v10a_a_582ea_26e |
| jackson2-api | 2.22.3-447.v6a_93ce7d2c82 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | 3.2.4-103.vf5a_b_3f4ea_cec | 3.2.4-103.vf5a_b_3f4ea_cec |
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
| jquery3-api | 3.7.1-687.v68d468e40b_30 | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.2.16-95.v3eecb_55fa_b_78 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | 20260814-226.v20f9685d642c | 20260814-226.v20f9685d642c |
| json-path-api | 3.0.0-218.vcd4dd1355de2 | 3.0.0-218.vcd4dd1355de2 |
| jsoup | 1.23.2-106.v1161f282b_eb_f | 1.23.2-106.v1161f282b_eb_f |
| junit | 1442.vb_b_a_467df2767 | 1442.vb_b_a_467df2767 |
| ldap | 825.v2fca_37dd5b_cb_ | 825.v2fca_37dd5b_cb_ |
| mailer | 534.v1b_36f5864073 | 534.v1b_36f5864073 |
| matrix-auth | 3.3 | 3.3 |
| matrix-project | 928.v1effc3f3a_236 | 928.v1effc3f3a_236 |
| metrics | 4.2.37-494.v06f9a_939d33a_ | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| okhttp-api | 5.5.0-242.v1482d48d3ea_3 | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | 477.v7197b_eb_3d8b_8 | 477.v7197b_eb_3d8b_8 |
| pam-auth | 1.12 | 1.12 |
| pipeline-build-step | 601.v6d4c6d1a_9dc7 | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 65.v203688e7727e | 65.v203688e7727e |
| pipeline-graph-analysis | 254.v0f63a_a_447dca_ | 254.v0f63a_a_447dca_ |
| pipeline-groovy-lib | 822.vc95f70dea_db_7 | 822.vc95f70dea_db_7 |
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
| plugin-util-api | 7.1403.vd380b_c42e90c | 7.1403.vd380b_c42e90c |
| prism-api | 1.30.0-741.v034eb_0b_0a_a_fa_ | 1.30.0-741.v034eb_0b_0a_a_fa_ |
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
| workflow-api | 1413.v2ff1a_5e720fa_ | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 1098.v808b_fd7f8cf4 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 4395.v738456a_5da_6b_ | 4395.v738456a_5da_6b_ |
| workflow-durable-task-step | 1479.v56e587f413a_7 | 1479.v56e587f413a_7 |
| workflow-job | 1609.v8926e7962024 | 1609.v8926e7962024 |
| workflow-multibranch | 842.v3a_b_59b_57b_e6e | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 487.v41313a_96b_65d | 487.v41313a_96b_65d |
| workflow-step-api | 724.v538c2362b_dfb_ | 724.v538c2362b_dfb_ |
| workflow-support | 1015.v785e5a_b_b_8b_22 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.49 | 0.49 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
