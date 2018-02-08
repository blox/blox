# Setup CodeBuild for End to End Test

We use CodeBuild to run quick end-to-end test. Below are steps to set it up:

### Running `deploy` to setup blox application
To avoid too broad permissions required for CodeBuild project, we setup blox project firstly.

Make sure you have one AWS credential from the AWS account and used in **default** profile:

```
./gradlew -Pblox.region=us-west-2 -Pblox.profile=default -Pblox.name=blox-codebuild-us-west-2 createBucket
./gradlew -Pblox.region=us-west-2 -Pblox.profile=default -Pblox.name=blox-codebuild-us-west-2 deploy
```

### Create CodeBuild Project
Before we create the project, you must connect the AWS account to the GitHub account using AWS CodeBuild console as following:
 - Use the AWS CodeBuild console to create a build project
 - When you use the console to connect (or reconnect) with GitHub, on the GitHub Authorize application page, for Organization access, choose Request access next to blox
 - Choose Authorize application. You can close the AWS CodeBuild console.

Now you can create the CodeBuild project as following:

```
aws cloudformation create-stack --stack-name blox-code-build --template-body file://./end-to-end-tests/code-build/project_setup_cf_template.yml --capabilities CAPABILITY_IAM
```

Once the CloudFormation stack is finished, setup github webhook to this project:
```
aws codebuild create-webhook --project-name Blox-E2E-Test
```

Now you should see a webhook created in the blox repository.

