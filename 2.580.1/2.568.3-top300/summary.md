### 2.568.3 → 2.580.1, 300 most installed plugins

| Rehearsal | Image | Jenkins | Plugins | Failed | Causes | Partly loaded | Not downloaded |
|---|---|---|---|---|---|---|---|
| baseline: 2.568.3 | jenkins/jenkins:2.568.3-jdk21 | started | 308 | 75 | 38 | 0 | 9 |
| as-is: 2.580.1, plugins unchanged | jenkins/jenkins:2.580.1-jdk21 | started | 317 | 0 | 0 | 0 | 0 |
| planned: 2.580.1, plugins from plan | jenkins/jenkins:2.580.1-jdk21 | started | 317 | 0 | 0 | 0 | 0 |

- baseline: **junit 1442.vb_b_a_467df2767** missing-dependency: needs the checks-api plugin (415.vf022234a_931d or newer), which is not installed (takes down 31 more)
- baseline: **github 1.47.0** missing-dependency: needs the git plugin (5.8.1 or newer), which is not installed (takes down 2 more)
- baseline: **aws-java-sdk2-apigateway 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed (takes down 1 more)
- baseline: **aws-java-sdk2-ec2 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed (takes down 1 more)
- baseline: **code-coverage-api 4.99.0** missing-dependency: needs the checks-api plugin (2.0.2 or newer), which is not installed (takes down 1 more)
- baseline: **data-tables-api 2.3.8-1570.v1cb_1cd2a_0fb_c** missing-dependency: needs the jquery3-api plugin (3.7.1-682.vfa_cdce169929 or newer), which is not installed (takes down 1 more)
- baseline: **aws-java-sdk 1.12.780-480.v4a_0819121a_9e** missing-dependency: needs the aws-java-sdk-secretsmanager plugin (1.12.780-480.v4a_0819121a_9e or newer), which is not installed
- baseline: **aws-java-sdk2-cloudformation 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-cloudfront 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-codedeploy 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-ecr 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-elasticbeanstalk 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-elasticloadbalancingv2 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-iam 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-lambda 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-netty-nio-client 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-organizations 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-s3 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **aws-java-sdk2-sns 2.42.33-70.vd69c0763fa_60** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **bitbucket 322.v882c6a_b_82c52** missing-dependency: needs the git plugin (5.7.0 or newer), which is not installed
- baseline: **blueocean-autofavorite 1.2.5** missing-dependency: needs the git plugin (3.6.0 or newer), which is not installed
- baseline: **blueocean-personalization 1.27.26** missing-dependency: needs the blueocean-dashboard plugin (1.27.26 or newer), which is not installed
- baseline: **blueocean-pipeline-editor 1.27.26** missing-dependency: needs the blueocean-dashboard plugin (1.27.26 or newer), which is not installed
- baseline: **bootstrap4-api 4.6.0-6** missing-dependency: needs the jquery3-api plugin (3.6.0-2 or newer), which is not installed
- baseline: **build-pipeline-plugin 2.2.0** missing-dependency: needs the jquery3-api plugin (3.7.1-594.vb_3864f326cf0 or newer), which is not installed
- baseline: **cloudbees-bitbucket-branch-source 937.3.10** missing-dependency: needs the git plugin (5.10.1 or newer), which is not installed
- baseline: **coverage 3.3497.v3fee733ca_b_35** missing-dependency: needs the checks-api plugin (422.vc7e8b_51488c6 or newer), which is not installed
- baseline: **ec2 2064.v79de4fd2d03b_** missing-dependency: needs the aws-java-sdk2-core plugin (2.42.33-70.vd69c0763fa_60 or newer), which is not installed
- baseline: **echarts-api 6.1.0-1321.v86c4c5114494** missing-dependency: needs the jquery3-api plugin (3.7.1-687.v68d468e40b_30 or newer), which is not installed
- baseline: **git-parameter 462.463.v496a_59f698e5** missing-dependency: needs the jquery3-api plugin (3.7.1-3 or newer), which is not installed
- baseline: **github-oauth 685.v53b_070455063** missing-dependency: needs the git plugin (5.7.0 or newer), which is not installed
- baseline: **gitlab-branch-source 744.vb_d0403d08ec7** missing-dependency: needs the git plugin (5.8.1 or newer), which is not installed
- baseline: **gitlab-oauth 1.30** missing-dependency: needs the git plugin (5.8.1 or newer), which is not installed
- baseline: **gitlab-plugin 1.2158.v49db_f64b_fc28** missing-dependency: needs the git plugin (5.10.1 or newer), which is not installed
- baseline: **groovy-postbuild 335.vf10eb_c1b_5e14** missing-dependency: needs the badge plugin (3.605.v7a_3387b_97a_b_6 or newer), which is not installed
- baseline: **pipeline-github-lib 65.v203688e7727e** missing-dependency: needs the git plugin (5.7.0 or newer), which is not installed
- baseline: **uno-choice 2.8.11** missing-dependency: needs the jquery3-api plugin (3.7.1-594.vb_3864f326cf0 or newer), which is not installed
- baseline: **warnings-ng 13.10305.v3231d55c5755** missing-dependency: needs the checks-api plugin (422.vc7e8b_51488c6 or newer), which is not installed

