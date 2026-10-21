# Jenkins upgrade plan: 2.541.3 → 2.580.1

Source: results/2.541.3-top300/plugins-2.541.3.txt · generated 2026-10-21 13:54 UTC by jenkins-upgrade-check dev

**1 core upgrade(s)**, 317 plugins: 25 to upgrade, 0 new dependencies; 7 with security warnings today, 3 still with warnings after the upgrade, 1 deprecated.

Jenkins 2.541.3 runs on Java 17, 21, 25; Jenkins 2.580.1 runs on Java 21, 25.

## Security warnings for the running core

- [core-2_568](https://www.jenkins.io/security/advisory/2026-06-10/): Multiple security vulnerabilities in Jenkins 2.567 and earlier, LTS 2.555.2 and earlier
- [core-2_576](https://www.jenkins.io/security/advisory/2026-08-05/): Multiple security vulnerabilities in Jenkins 2.575 and earlier, LTS 2.568.1 and earlier
- [core-2_580](https://www.jenkins.io/security/advisory/2026-09-02/): Multiple security vulnerabilities in Jenkins 2.579 and earlier, LTS 2.568.2 and earlier

## Steps

1. **On 2.541.3, update plugins** (0 change(s)). Back up `$JENKINS_HOME` first.
2. **Run the controller and every agent on Java 21**, upgrade Jenkins to **2.580.1**, start it, then update plugins (25 change(s)).

## Needs a decision

| Plugin | Now | After | Why |
|---|---|---|---|
| Parameterized-Remote-Trigger | 3.2.2 | 3.2.2 | [SECURITY-3768](https://www.jenkins.io/security/advisory/2026-08-05/#SECURITY-3768) **no fix yet**<br>[SECURITY-3800](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3800) **no fix yet** |
| coverage | 2.3139.v731d2e83ea_8a_ | 3.3497.v3fee733ca_b_35 | [SECURITY-4118](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-4118) fixed by the upgrade |
| extended-choice-parameter | 388.ve7b_d0b_920e10 | 388.ve7b_d0b_920e10 | [SECURITY-1350](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1350) **no fix yet**<br>[SECURITY-1351](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-1351) **no fix yet**<br>[SECURITY-2232](https://www.jenkins.io/security/advisory/2022-03-15/#SECURITY-2232) **no fix yet**<br>[SECURITY-2617-extended-choice-parameter](https://www.jenkins.io/security/advisory/2022-04-12/#SECURITY-2617) **no fix yet**<br>[deprecated](https://github.com/jenkinsci/extended-choice-parameter-plugin?tab=readme-ov-file#end-of-life) |
| pipeline-groovy-lib | 798.v5cc688825312 | 822.vc95f70dea_db_7 | [SECURITY-3815](https://www.jenkins.io/security/advisory/2026-09-02/#SECURITY-3815) fixed by the upgrade<br>[SECURITY-3796](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3796) fixed by the upgrade |
| ssh | 159.v496ca_25e5e82 | 159.v496ca_25e5e82 | [SECURITY-2093](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2093) **no fix yet**<br>[SECURITY-2315](https://www.jenkins.io/security/advisory/2022-05-17/#SECURITY-2315) **no fix yet** |
| warnings-ng | 13.10085.v94a_d6ca_82b_f6 | 13.10305.v3231d55c5755 | [SECURITY-3937](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3937) fixed by the upgrade |
| workflow-multibranch | 841.vec5b_9e1806ec | 842.v3a_b_59b_57b_e6e | [SECURITY-3729](https://www.jenkins.io/security/advisory/2026-09-16/#SECURITY-3729) fixed by the upgrade |

<details><summary>Every plugin</summary>

| Plugin | Now | After |
|---|---|---|
| Office-365-Connector | 5.3.1 | 5.3.1 |
| Parameterized-Remote-Trigger | 3.2.2 | 3.2.2 |
| active-directory | 2.971.v653d8b_6e5548 | 2.971.v653d8b_6e5548 |
| allure-jenkins-plugin | 3.0.0 | 3.0.0 |
| analysis-model-api | 13.21.0-945.va_f5a_0b_c1404b_ | 14.18.0-1028.v9c938eeda_ef2 |
| ansible | 635.v34b_48e979c86 | 635.v34b_48e979c86 |
| ansicolor | 542.v03d235fee02d | 542.v03d235fee02d |
| ant | 520.vd082ecfb_16a_9 | 520.vd082ecfb_16a_9 |
| antisamy-markup-formatter | 173.v680e3a_b_69ff3 | 173.v680e3a_b_69ff3 |
| apache-httpcomponents-client-4-api | 4.5.14-269.vfa_2321039a_83 | 4.5.14-269.vfa_2321039a_83 |
| apache-httpcomponents-client-5-api | 5.6.4-207.vd5b_a_c15e4de5 | 5.6.4-207.vd5b_a_c15e4de5 |
| artifactory | 4.0.8 | 4.0.8 |
| asm-api | 9.11.1-226.vff833a_e4c58e | 9.11.1-226.vff833a_e4c58e |
| audit-trail | 458.vf8493ccb_63b_1 | 464.v78b_845ffc5d8 |
| authentication-tokens | 1.144.v5ff4a_5ec5c33 | 1.144.v5ff4a_5ec5c33 |
| authorize-project | 534.v2f208c45e11c | 534.v2f208c45e11c |
| aws-credentials | 265.v8422a_f384cd9 | 265.v8422a_f384cd9 |
| aws-java-sdk | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-api-gateway | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-autoscaling | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-cloudformation | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-cloudfront | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-cloudwatch | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-codebuild | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-codedeploy | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ec2 | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ecr | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ecs | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-efs | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-elasticbeanstalk | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-elasticloadbalancingv2 | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-iam | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-kinesis | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-lambda | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-logs | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-minimal | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-organizations | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-secretsmanager | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-sns | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-sqs | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk-ssm | 1.12.780-480.v4a_0819121a_9e | 1.12.780-480.v4a_0819121a_9e |
| aws-java-sdk2-apigateway | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-cloudformation | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-cloudfront | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-codedeploy | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-core | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-ec2 | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-ecr | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-elasticbeanstalk | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-elasticloadbalancingv2 | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-iam | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-lambda | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-netty-nio-client | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-organizations | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-s3 | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| aws-java-sdk2-sns | 2.42.33-70.vd69c0763fa_60 | 2.42.33-70.vd69c0763fa_60 |
| azure-sdk | 299.v3996878fe83a_ | 299.v3996878fe83a_ |
| badge | 3.615.vb_ed967ca_b_e93 | 3.615.vb_ed967ca_b_e93 |
| basic-branch-build-strategies | 317.v85b_331d6cc42 | 317.v85b_331d6cc42 |
| bitbucket | 322.v882c6a_b_82c52 | 322.v882c6a_b_82c52 |
| blueocean | 1.27.26 | 1.27.26 |
| blueocean-autofavorite | 1.2.5 | 1.2.5 |
| blueocean-bitbucket-pipeline | 1.27.26 | 1.27.26 |
| blueocean-commons | 1.27.26 | 1.27.26 |
| blueocean-config | 1.27.26 | 1.27.26 |
| blueocean-core-js | 1.27.26 | 1.27.26 |
| blueocean-dashboard | 1.27.26 | 1.27.26 |
| blueocean-display-url | 2.4.4 | 2.4.4 |
| blueocean-events | 1.27.26 | 1.27.26 |
| blueocean-git-pipeline | 1.27.26 | 1.27.26 |
| blueocean-github-pipeline | 1.27.26 | 1.27.26 |
| blueocean-i18n | 1.27.26 | 1.27.26 |
| blueocean-jira | 1.27.26 | 1.27.26 |
| blueocean-jwt | 1.27.26 | 1.27.26 |
| blueocean-personalization | 1.27.26 | 1.27.26 |
| blueocean-pipeline-api-impl | 1.27.26 | 1.27.26 |
| blueocean-pipeline-editor | 1.27.26 | 1.27.26 |
| blueocean-pipeline-scm-api | 1.27.26 | 1.27.26 |
| blueocean-rest | 1.27.26 | 1.27.26 |
| blueocean-rest-impl | 1.27.26 | 1.27.26 |
| blueocean-web | 1.27.26 | 1.27.26 |
| bootstrap4-api | 4.6.0-6 | 4.6.0-6 |
| bootstrap5-api | 5.3.8-895.v4d0d8e47fea_d | 5.3.8-1048.va_c299057e35c |
| bouncycastle-api | 2.30.1.86.2-310.v4655f6742752 | 2.30.1.86.2-310.v4655f6742752 |
| branch-api | 2.1303.v9f3b_95dc329d | 2.1303.v9f3b_95dc329d |
| build-blocker-plugin | 177.ve433d85b_3fd7 | 177.ve433d85b_3fd7 |
| build-monitor-plugin | 1131.va_2d8e68fc912 | 1131.va_2d8e68fc912 |
| build-name-setter | 2.5.1 | 2.5.1 |
| build-pipeline-plugin | 2.2.0 | 2.2.0 |
| build-timeout | 1.41 | 1.41 |
| build-timestamp | 1.1.1 | 1.1.1 |
| build-user-vars-plugin | 214.va_eed2ed849ca_ | 214.va_eed2ed849ca_ |
| build-with-parameters | 81.ve4a_9c2499d9a | 81.ve4a_9c2499d9a |
| built-on-column | 1.5 | 1.5 |
| caffeine-api | 3.3.1-214.v18eb_56b_34920 | 3.3.1-214.v18eb_56b_34920 |
| checks-api | 422.vc7e8b_51488c6 | 422.vc7e8b_51488c6 |
| cloud-stats | 423.v63f017a_f582c | 423.v63f017a_f582c |
| cloudbees-bitbucket-branch-source | 937.3.10 | 937.3.10 |
| cloudbees-disk-usage-simple | 274.v5e87a_f798201 | 275.v56e74c9dd0d3 |
| cloudbees-folder | 6.1106.v3a_d9a_6d2465e | 6.1106.v3a_d9a_6d2465e |
| cobertura | 1.17 | 1.17 |
| code-coverage-api | 4.99.0 | 4.99.0 |
| command-launcher | 141.v9578cf36b_311 | 141.v9578cf36b_311 |
| commons-collections4-api | 4.5.0-8.va_d5448ef9011 | 4.5.0-8.va_d5448ef9011 |
| commons-compress-api | 1.28.0-96.v8da_8cfa_57e2b_ | 1.28.0-96.v8da_8cfa_57e2b_ |
| commons-httpclient3-api | 3.1-3 | 3.1-3 |
| commons-lang3-api | 3.21.0-122.v6a_ea_4fd34049 | 3.21.0-122.v6a_ea_4fd34049 |
| commons-text-api | 1.15.0-218.va_61573470393 | 1.15.0-218.va_61573470393 |
| conditional-buildstep | 1.5.0 | 1.5.0 |
| config-file-provider | 1006.vc7366c201f57 | 1013.v73c323e52b_1f |
| configuration-as-code | 2136.veb_3a_554550b_f | 2136.veb_3a_554550b_f |
| copyartifact | 805.v9d1393350780 | 805.v9d1393350780 |
| coverage | 2.3139.v731d2e83ea_8a_ | 3.3497.v3fee733ca_b_35 |
| credentials | 1511.v2e3cb_0008ef0 | 1514.vb_01eb_014e7d4 |
| credentials-binding | 728.v902a_273b_8947 | 728.v902a_273b_8947 |
| cvs | 502.v02d5a_0c6b_9a_3 | 502.v02d5a_0c6b_9a_3 |
| dark-theme | 652.vea_da_dfea_e769 | 652.vea_da_dfea_e769 |
| dashboard-view | 2.569.v5ddb_c9b_a_7264 | 2.569.v5ddb_c9b_a_7264 |
| data-tables-api | 2.3.7-1534.v539d4edf109d | 2.3.8-1570.v1cb_1cd2a_0fb_c |
| display-url-api | 2.217.va_6b_de84cc74b_ | 2.217.va_6b_de84cc74b_ |
| docker-build-step | 2.14 | 2.14 |
| docker-commons | 477.v289085a_b_6896 | 477.v289085a_b_6896 |
| docker-java-api | 3.7.1-136.v5d70f77a_c3d6 | 3.7.1-136.v5d70f77a_c3d6 |
| docker-plugin | 1335.vfc7e88a_14461 | 1336.v3711dec712de |
| docker-workflow | 653.v2f2c08eff0ec | 653.v2f2c08eff0ec |
| dtkit-api | 3.0.3 | 3.0.3 |
| durable-task | 686.v80ff80875b_82 | 686.v80ff80875b_82 |
| ec2 | 2064.v79de4fd2d03b_ | 2064.v79de4fd2d03b_ |
| echarts-api | 6.0.0-1247.vf3e35a_c1813f | 6.1.0-1321.v86c4c5114494 |
| eddsa-api | 0.3.0.1-29.v67e9a_1c969b_b_ | 0.3.0.1-29.v67e9a_1c969b_b_ |
| email-ext | 2038.v7b_8817a_499d9 | 2038.v7b_8817a_499d9 |
| emailext-template | 249.v81d6ea_f33f2e | 249.v81d6ea_f33f2e |
| embeddable-build-status | 666.vcb_0de9eb_43d8 | 666.vcb_0de9eb_43d8 |
| emoji-symbols-api | 17.0-57.v8d44b_9a_b_d5ea_ | 17.0-57.v8d44b_9a_b_d5ea_ |
| envinject | 2.941.v351a_20c0a_3ca_ | 2.941.v351a_20c0a_3ca_ |
| envinject-api | 1.241.vdd714803b_403 | 1.241.vdd714803b_403 |
| extended-choice-parameter | 388.ve7b_d0b_920e10 | 388.ve7b_d0b_920e10 |
| extended-read-permission | 92.v1b_f0c20b_75db_ | 92.v1b_f0c20b_75db_ |
| extensible-choice-parameter | 255.vfa_41f46828ef | 255.vfa_41f46828ef |
| external-monitor-job | 223.vb_fddcf42c9b_3 | 223.vb_fddcf42c9b_3 |
| favorite | 2.263.v941d21defef7 | 2.267.vb_90d08408081 |
| file-parameters | 433.va_0b_80359d54d | 433.va_0b_80359d54d |
| flatpickr-api | 4.6.13-32.v60a_51029c136 | 4.6.13-32.v60a_51029c136 |
| font-awesome-api | 7.2.0-965.ve3840b_696418 | 7.3.1-1013.v0835a_879ec6d |
| forensics-api | 3.1832.va_1179842528b_ | 4.1923.v445c56df9b_2e |
| generic-webhook-trigger | 2.4.3 | 2.4.3 |
| git | 5.11.0 | 5.11.0 |
| git-client | 6.6.1 | 6.6.1 |
| git-parameter | 462.463.v496a_59f698e5 | 462.463.v496a_59f698e5 |
| git-server | 137.ve0060b_432302 | 137.ve0060b_432302 |
| github | 1.47.0 | 1.47.0 |
| github-api | 1.330-492.v3941a_032db_2a_ | 1.330-492.v3941a_032db_2a_ |
| github-branch-source | 1986.v801e40614d00 | 1986.v801e40614d00 |
| github-oauth | 685.v53b_070455063 | 685.v53b_070455063 |
| github-pullrequest | 0.7.5 | 0.7.5 |
| gitlab-api | 6.3.0-126.vd066b_d896159 | 6.3.0-126.vd066b_d896159 |
| gitlab-branch-source | 744.vb_d0403d08ec7 | 744.vb_d0403d08ec7 |
| gitlab-oauth | 1.30 | 1.30 |
| gitlab-plugin | 1.2158.v49db_f64b_fc28 | 1.2158.v49db_f64b_fc28 |
| google-oauth-plugin | 1.346.v712b_a_e0d1e60 | 1.346.v712b_a_e0d1e60 |
| gradle | 2.20.1253.vc116f0763a_eb_ | 2.20.1253.vc116f0763a_eb_ |
| groovy | 537.v741a_5a_f1b_581 | 537.v741a_5a_f1b_581 |
| groovy-postbuild | 335.vf10eb_c1b_5e14 | 335.vf10eb_c1b_5e14 |
| gson-api | 2.14.0-201.v8eefe5515533 | 2.14.0-201.v8eefe5515533 |
| h2-api | 11.1.4.199-36.vb_ee07e965744 | 11.1.4.199-36.vb_ee07e965744 |
| handy-uri-templates-2-api | 2.1.8-38.vcea_5d521d5f3 | 2.1.8-38.vcea_5d521d5f3 |
| hashicorp-vault-plugin | 384.vda_86ec66c537 | 384.vda_86ec66c537 |
| htmlpublisher | 429 | 429 |
| http_request | 1.659.v1b_9d9e942a_e6 | 1.659.v1b_9d9e942a_e6 |
| icon-shim | 3.0.0 | 3.0.0 |
| instance-identity | 203.v15e81a_1b_7a_38 | 203.v15e81a_1b_7a_38 |
| ionicons-api | 94.vcc3065403257 | 94.vcc3065403257 |
| ivy | 609.v66779e68503d | 609.v66779e68503d |
| jackson-annotations2-api | 2.22-19.v10a_a_582ea_26e | 2.22-19.v10a_a_582ea_26e |
| jackson2-api | 2.22.3-447.v6a_93ce7d2c82 | 2.22.3-447.v6a_93ce7d2c82 |
| jackson3-api | 3.2.4-103.vf5a_b_3f4ea_cec | 3.2.4-103.vf5a_b_3f4ea_cec |
| jacoco | 3.3.7 | 3.3.7 |
| jakarta-activation-api | 2.1.4-1 | 2.1.4-1 |
| jakarta-mail-api | 2.1.5-1 | 2.1.5-1 |
| jakarta-xml-bind-api | 4.0.9-19.v2b_a_5b_44d9a_1c | 4.0.9-19.v2b_a_5b_44d9a_1c |
| javadoc | 354.vee1a_660b_4990 | 354.vee1a_660b_4990 |
| javax-activation-api | 1.2.0-8 | 1.2.0-8 |
| javax-mail-api | 1.6.2-11 | 1.6.2-11 |
| jaxb | 2.3.9-143.v5979df3304e6 | 2.3.9-143.v5979df3304e6 |
| jdk-tool | 83.v417146707a_3d | 83.v417146707a_3d |
| jenkins-design-language | 1.27.26 | 1.27.26 |
| jenkins-multijob-plugin | 682.v4ed91b_52d621 | 682.v4ed91b_52d621 |
| jersey2-api | 2.48-180.ve47b_264f849b_ | 2.48-180.ve47b_264f849b_ |
| jersey3-api | 3.1.11-4.v77818819c2e1 | 3.1.11-4.v77818819c2e1 |
| jira | 3.21 | 3.22 |
| jjwt-api | 0.13.0-141.vd58b_a_9592b_6c | 0.13.0-141.vd58b_a_9592b_6c |
| jnr-posix-api | 3.2.4-223.v4f0cfd7fa_006 | 3.2.4-223.v4f0cfd7fa_006 |
| job-dsl | 3732.v9a_c49a_61a_313 | 3732.v9a_c49a_61a_313 |
| job-import-plugin | 143.145.v48f9a_a_6ff384 | 143.145.v48f9a_a_6ff384 |
| jobConfigHistory | 1380.v762185b_9a_793 | 1380.v762185b_9a_793 |
| joda-time-api | 2.15.0-208.v6b_084e7a_d9f0 | 2.15.0-208.v6b_084e7a_d9f0 |
| jquery | 1.12.4-3 | 1.12.4-3 |
| jquery-detached | 1.2.1 | 1.2.1 |
| jquery3-api | 3.7.1-619.vdb_10e002501a_ | 3.7.1-687.v68d468e40b_30 |
| jsch | 0.2.16-95.v3eecb_55fa_b_78 | 0.2.16-95.v3eecb_55fa_b_78 |
| json-api | 20260814-226.v20f9685d642c | 20260814-226.v20f9685d642c |
| json-path-api | 3.0.0-218.vcd4dd1355de2 | 3.0.0-218.vcd4dd1355de2 |
| jsoup | 1.23.2-106.v1161f282b_eb_f | 1.23.2-106.v1161f282b_eb_f |
| junit | 1442.vb_b_a_467df2767 | 1442.vb_b_a_467df2767 |
| kubernetes | 4560.vc25b_83890c34 | 4560.vc25b_83890c34 |
| kubernetes-cli | 1.375.v8001982544ba | 1.375.v8001982544ba |
| kubernetes-client-api | 7.3.1-256.v788a_0b_787114 | 7.3.1-256.v788a_0b_787114 |
| kubernetes-credentials | 207.v492f58828b_ed | 207.v492f58828b_ed |
| ldap | 825.v2fca_37dd5b_cb_ | 825.v2fca_37dd5b_cb_ |
| locale | 641.v84f22d75fd8b_ | 641.v84f22d75fd8b_ |
| localization-support | 1.41.v186a_c1569458 | 1.41.v186a_c1569458 |
| localization-zh-cn | 371.v23851f835d6b_ | 371.v23851f835d6b_ |
| lockable-resources | 1576.v217259b_f9976 | 1576.v217259b_f9976 |
| log-parser | 3.0.7 | 3.0.7 |
| mailer | 534.v1b_36f5864073 | 534.v1b_36f5864073 |
| mapdb-api | 1.0.9-44.va_1e1310c9118 | 1.0.9-44.va_1e1310c9118 |
| mask-passwords | 220.v95819055b_265 | 220.v95819055b_265 |
| matrix-auth | 3.3 | 3.3 |
| matrix-project | 905.vcc6831e8760a_ | 928.v1effc3f3a_236 |
| maven-plugin | 3.28 | 3.28 |
| mercurial | 1323.ve69d2a_db_8a_b_d | 1323.ve69d2a_db_8a_b_d |
| metrics | 4.2.37-494.v06f9a_939d33a_ | 4.2.37-494.v06f9a_939d33a_ |
| mina-sshd-api-common | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-core | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| mina-sshd-api-scp | 2.19.0-192.v2b_a_7b_2c1dc71 | 2.19.0-192.v2b_a_7b_2c1dc71 |
| monitoring | 2.9.1 | 2.9.1 |
| msbuild | 1.44 | 1.45 |
| multiple-scms | 0.8 | 0.8 |
| naginator | 1.556.v14d723a_109a_c | 1.556.v14d723a_109a_c |
| netty-api | 4.2.18.Final-34.v861b_21a_78c0f | 4.2.18.Final-34.v861b_21a_78c0f |
| node-iterator-api | 72.vc90e81737df1 | 72.vc90e81737df1 |
| nodejs | 1.6.6 | 1.6.6 |
| nodelabelparameter | 851.vd94e5048d321 | 851.vd94e5048d321 |
| oauth-credentials | 0.668.vb_431feb_d8a_4e | 0.668.vb_431feb_d8a_4e |
| okhttp-api | 5.5.0-242.v1482d48d3ea_3 | 5.5.0-242.v1482d48d3ea_3 |
| oss-symbols-api | 477.v7197b_eb_3d8b_8 | 477.v7197b_eb_3d8b_8 |
| pam-auth | 1.12 | 1.12 |
| parameterized-scheduler | 379.v95c73f233a_df | 379.v95c73f233a_df |
| parameterized-trigger | 893.va_383a_9b_a_4c11 | 893.va_383a_9b_a_4c11 |
| performance | 1022.v48fce959f0c3 | 1022.v48fce959f0c3 |
| pipeline-aws | 767.v48d5c10b_1899 | 767.v48d5c10b_1899 |
| pipeline-build-step | 601.v6d4c6d1a_9dc7 | 601.v6d4c6d1a_9dc7 |
| pipeline-github-lib | 65.v203688e7727e | 65.v203688e7727e |
| pipeline-graph-analysis | 254.v0f63a_a_447dca_ | 254.v0f63a_a_447dca_ |
| pipeline-graph-view | 1013.v9f83fd83c063 | 1053.v9df19e6ceed2 |
| pipeline-groovy-lib | 798.v5cc688825312 | 822.vc95f70dea_db_7 |
| pipeline-input-step | 560.v56198a_642157 | 560.v56198a_642157 |
| pipeline-maven | 1760.v9a_a_e6dcb_0444 | 1760.v9a_a_e6dcb_0444 |
| pipeline-maven-api | 1760.v9a_a_e6dcb_0444 | 1760.v9a_a_e6dcb_0444 |
| pipeline-milestone-step | 152.v6e22b_8cfc66c | 152.v6e22b_8cfc66c |
| pipeline-model-api | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-model-declarative-agent | 1.1.1 | 1.1.1 |
| pipeline-model-definition | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-model-extensions | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-rest-api | 2.41 | 2.41 |
| pipeline-stage-step | 345.va_96187909426 | 345.va_96187909426 |
| pipeline-stage-tags-metadata | 2.2293.v6e7193cec599 | 2.2293.v6e7193cec599 |
| pipeline-stage-view | 2.41 | 2.41 |
| pipeline-utility-steps | 3.810.va_7672d206740 | 3.810.va_7672d206740 |
| plain-credentials | 199.v9f8e1f741799 | 199.v9f8e1f741799 |
| plugin-util-api | 6.1192.v30fe6e2837ff | 7.1403.vd380b_c42e90c |
| popper-api | 1.16.1-3 | 1.16.1-3 |
| popper2-api | 2.11.6-5 | 2.11.6-5 |
| postbuild-task | 78.v24529f1f5cdb_ | 78.v24529f1f5cdb_ |
| postbuildscript | 3.4.1.821.v08352159216f | 3.4.1.821.v08352159216f |
| powershell | 185.v7a_026da_c54ee | 185.v7a_026da_c54ee |
| prism-api | 1.30.0-717.vb_f8360844b_53 | 1.30.0-741.v034eb_0b_0a_a_fa_ |
| prometheus | 860.v532442b_44e9a_ | 860.v532442b_44e9a_ |
| publish-over | 242.vf56f6892cb_b_a_ | 242.vf56f6892cb_b_a_ |
| publish-over-ssh | 409.v73ca_23a_5b_48e | 409.v73ca_23a_5b_48e |
| pubsub-light | 1.19 | 1.19 |
| rebuild | 344.v805899551b_13 | 344.v805899551b_13 |
| resource-disposer | 0.25 | 0.25 |
| role-strategy | 867.vd09254229f9b_ | 927.v9cf5527c4085 |
| run-condition | 276.v97298f3a_cd51 | 276.v97298f3a_cd51 |
| saferestart | 125.vdb_73765c38a_e | 125.vdb_73765c38a_e |
| saml | 4.623.v7875d61cd9f5 | 4.623.v7875d61cd9f5 |
| scm-api | 728.vc30dcf7a_0df5 | 728.vc30dcf7a_0df5 |
| script-security | 1436.v1082791135e4 | 1436.v1082791135e4 |
| simple-theme-plugin | 230.v8b_fd91b_b_800c | 230.v8b_fd91b_b_800c |
| slack | 795.v4b_9705b_e6d47 | 795.v4b_9705b_e6d47 |
| snakeyaml-api | 2.7-152.vc672728cf536 | 2.7-152.vc672728cf536 |
| snakeyaml-engine-api | 3.2-18.v885ed7ff698f | 3.2-18.v885ed7ff698f |
| sonar | 2.19.0 | 2.19.0 |
| sse-gateway | 1.29 | 1.29 |
| ssh | 159.v496ca_25e5e82 | 159.v496ca_25e5e82 |
| ssh-agent | 433.v9e73d67c3f74 | 433.v9e73d67c3f74 |
| ssh-credentials | 372.va_250881b_08cd | 372.va_250881b_08cd |
| ssh-slaves | 3.1097.v868116049892 | 3.1097.v868116049892 |
| ssh-steps | 2.0.92.vb_a_0583935f9b_ | 2.0.92.vb_a_0583935f9b_ |
| sshd | 3.384.vc89b_5e138cf9 | 3.384.vc89b_5e138cf9 |
| structs | 362.va_b_695ef4fdf9 | 362.va_b_695ef4fdf9 |
| subversion | 1327.vc36e0019a_653 | 1327.vc36e0019a_653 |
| theme-manager | 346.v06cca_64c6a_37 | 346.v06cca_64c6a_37 |
| thinBackup | 2.1.5 | 2.1.5 |
| throttle-concurrents | 625.vc8b_e469e9a_b_c | 625.vc8b_e469e9a_b_c |
| timestamper | 1.30 | 1.30 |
| token-macro | 477.vd4f0dc3cb_cf1 | 477.vd4f0dc3cb_cf1 |
| trilead-api | 2.284.v1974ea_324382 | 2.284.v1974ea_324382 |
| uno-choice | 2.8.11 | 2.8.11 |
| variant | 70.va_d9f17f859e0 | 70.va_d9f17f859e0 |
| view-job-filters | 408.vb_ca_5e798e5b_5 | 408.vb_ca_5e798e5b_5 |
| warnings-ng | 13.10085.v94a_d6ca_82b_f6 | 13.10305.v3231d55c5755 |
| woodstox-core-api | 7.3.0-16.v017f30ed004c | 7.3.0-16.v017f30ed004c |
| workflow-aggregator | 608.v67378e9d3db_1 | 608.v67378e9d3db_1 |
| workflow-api | 1413.v2ff1a_5e720fa_ | 1413.v2ff1a_5e720fa_ |
| workflow-basic-steps | 1098.v808b_fd7f8cf4 | 1098.v808b_fd7f8cf4 |
| workflow-cps | 4395.v738456a_5da_6b_ | 4395.v738456a_5da_6b_ |
| workflow-cps-global-lib | 636.v3d2c06d018a_6 | 636.v3d2c06d018a_6 |
| workflow-durable-task-step | 1479.v56e587f413a_7 | 1479.v56e587f413a_7 |
| workflow-job | 1571.1585.v21fce07cd04d | 1609.v8926e7962024 |
| workflow-multibranch | 841.vec5b_9e1806ec | 842.v3a_b_59b_57b_e6e |
| workflow-scm-step | 487.v41313a_96b_65d | 487.v41313a_96b_65d |
| workflow-step-api | 724.v538c2362b_dfb_ | 724.v538c2362b_dfb_ |
| workflow-support | 1015.v785e5a_b_b_8b_22 | 1015.v785e5a_b_b_8b_22 |
| ws-cleanup | 0.49 | 0.49 |
| xunit | 3.1.7 | 3.1.7 |

</details>

---
This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/
