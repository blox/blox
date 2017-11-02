# ![Logo](blox-logo.png)

# Blox: Open Source schedulers for Amazon ECS

[![Build Status](https://travis-ci.org/blox/blox.svg?branch=master)](https://travis-ci.org/blox/blox)

Blox provides open source schedulers optimized for running applications on Amazon ECS. Developers now have greater control over how their applications are deployed across clusters of resources, run and scale in production, and can take advantage of powerful placement capabilities of Amazon ECS.
Blox is being delivered as a managed service via the Amazon ECS Console, API and CLIs. Blox v1.0 provides daemon scheduling for Amazon ECS. We will continue to add additional schedulers as part of this project.
Blox schedulers are built using AWS primitives, and the Blox designs and code are open source. If you are interested in learning more or collaborating on the designs, please read the [design](docs/daemon_design.md).
If you are currently using Blox v0.3, please read the [FAQ](FAQ.md).

### Project structure
For an overview of the components of Blox, run:

```
./gradlew projects
```

### Testing
To run the full unit test suite, run:

```
./gradlew check
```

This will run the same tests that we run in the Travis CI build.

### Deploying
First, take a look at what Blox will put in your personal stack by running the
`showStackConfig` task:

```
$ ./gradlew showStackConfig

> Task :showStackConfig
Blox deployment stack configuration:

  Default resource name         (blox.name): blox-<username>-alpha-us-west-2 (default)
  API Gateway stage            (blox.stage): alpha (default)
  Stack prefix                (blox.prefix): <username>-alpha (default)
  AWS Region                  (blox.region): us-west-2 (default)
  AWS Credential Profile     (blox.profile): blox-<username>-alpha-us-west-2 (default)
  Cloudformation stack name (blox.cfnStack): blox-<username>-alpha-us-west-2 (default)
  Deployment S3 bucket name (blox.s3Bucket): blox-<username>-alpha-us-west-2 (default)

To customize these values, modify ~/.gradle/gradle.properties to override the property listed.

AWS CLI configuration for profile blox-<username>-alpha-us-west-2:

The config profile (blox-<username>-alpha-us-west-2) could not be found
```

If you wish to customize any of these values, you can do so by overriding the
property in parentheses using [any of the supported ways to override Gradle
properties](https://docs.gradle.org/current/userguide/build_environment.html#sec:gradle_properties_and_system_properties).
The easiest way is to override it for your user in `~/.gradle/gradle.properties`:

```
blox.profile=default
blox.region=us-east-1
```

Next, in order to deploy your personal stack:
- install the [official AWS CLI](https://aws.amazon.com/cli/)
- create an IAM user with the following permissions:

    ```json
    {
        "Version":"2012-10-17",
        "Statement":[{
            "Effect":"Allow",
            "Action":[
                "s3:*",
                "lambda:*",
                "apigateway:*",
                "cloudformation:*",
                "iam:*"
            ],
            "Resource":"*"
        }]
    }

    ```

  These permissions are pretty broad, so we recommend you use a separate, test account.

- configure the AWS Credential Profile shown in the `showStackOutput` task with
  the AWS credentials for the user you created above:

    ```
    aws configure --profile blox-<username>-alpha-us-west-2
    ```

- create an S3 bucket where all resources (code, cloudformation templates, etc) to be deployed will be stored:

    ```
    ./gradlew createBucket
    ```

- deploy the Blox stack:

    ```
    ./gradlew deploy
    ```

### End to end testing
Once you have a stack deployed, you can test it with:

```
./gradlew testEndToEnd
```


### Contact

* [Gitter](https://gitter.im/blox)
* [Planning/Roadmap](https://github.com/blox/blox/milestones)
* [Issues](https://github.com/blox/blox/issues)

### License
All projects under Blox are released under Apache 2.0 and contributions are accepted under individual Apache Contributor Agreements.
