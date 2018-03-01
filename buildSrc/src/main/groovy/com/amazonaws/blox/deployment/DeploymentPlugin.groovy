package com.amazonaws.blox.deployment

import com.amazonaws.blox.tasks.PostProcessCloudformation
import groovy.json.JsonSlurper
import groovy.transform.PackageScope
import org.gradle.api.NamedDomainObjectContainer
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.tasks.Exec

/**
 * Gradle plugin to deploy Blox as a Serverless Application Model stack.
 */
class DeploymentPlugin implements Plugin<Project> {
    def deployment

    @Override
    void apply(Project project) {
        deployment = project.extensions.create("deployment", DeploymentPluginExtension, project)

        // Only create tasks after the Gradle configuration phase, so that we can make use of the
        // configured `deployment` DSL block
        project.afterEvaluate(this.&createTasks)
    }

    private def aws(... args) {
        return [deployment.aws.command, "--profile", deployment.aws.profile, "--region", deployment.aws.region, *args]
    }

    void createTasks(Project project) {
        deployment.outputs = new StackOutputs(deployment.getStackOutputsFile())

        def postProcessTask = project.task("postprocessCloudformationTemplate", type: PostProcessCloudformation) {
            description "Postprocess the Cloudformation template to insert Swagger/Lambda references."

            templateFile deployment.templateFile
            outputTemplateFile deployment.getProcessedTemplateFile()

            deployment.apis.each {
                api it.name, project.files(it.swaggerFile)
            }

            deployment.lambdaFunctions.each {
                lambdaFunction it.name, project.files(it.zipFile)
            }
        }
        project.assemble.dependsOn(postProcessTask)


        def packageTask = project.task("packageCloudformationResources", type: Exec) {
            description "Use the CloudFormation package command to upload the deployment bundle to S3."

            inputs.files postProcessTask
            deployment.lambdaFunctions.each { inputs.files(it.zipFile) }
            outputs.file deployment.getPackagedTemplateFile()

            commandLine aws("cloudformation", "package",
                    "--template-file", deployment.getProcessedTemplateFile(),
                    "--output-template-file", deployment.getPackagedTemplateFile(),
                    "--s3-bucket", deployment.s3Bucket)
        }

        project.task("deploy") {
            group "deployment"
            description "Deploy the CloudFormation stack defined by a template file."

            inputs.files packageTask
            outputs.file deployment.outputs.file

            doLast {
                def error = new ByteArrayOutputStream()

                def arguments = ["cloudformation", "deploy",
                                 "--template-file", deployment.getPackagedTemplateFile(),
                                 "--stack-name", deployment.stackName,
                                 "--capabilities", "CAPABILITY_IAM"]

                if (!deployment.parameters.isEmpty()) {
                    arguments << "--parameter-overrides"
                    arguments += deployment.parameters.collect { "${it.key}=${it.value}" }
                }

                def result = project.exec {
                    commandLine aws(*arguments)

                    errorOutput error
                    ignoreExitValue true
                }

                // HACK: The `deploy` command returns a nonzero status if the stack is
                // up to date.  We can remove this once
                // https://github.com/awslabs/serverless-application-model/issues/71 is
                // fixed.
                def errorString = error.toString()
                if (!(errorString.contains("No changes to deploy") || errorString.contains("The submitted information didn't contain changes"))) {
                    System.err.println(error.toString())
                    result.assertNormalExitValue()
                }

                // In order to make this task incremental, we store the stack outputs
                // from deploying the stack as a file. That way tasks that depend on
                // this one (such as downloadSdk) don't have to do a redeploy unless
                // there's actual code changes.
                def output = new ByteArrayOutputStream()
                project.exec {
                    commandLine aws("cloudformation", "describe-stacks",
                            "--stack-name", deployment.stackName,
                            "--query", "Stacks[0]",
                            "--output", "json")
                    standardOutput output
                }

                deployment.outputs.write(output.toString())
            }
        }

        project.task("createBucket", type: Exec) {
            group "deployment"
            description "Create the S3 bucket used to store CloudFormation/Lambda resources for deployment."

            commandLine aws("s3", "mb", "s3://${deployment.s3Bucket}")
        }

        project.task("deleteBucket", type: Exec) {
            group "deployment"
            description "Delete the S3 bucket used to store CloudFormation/Lambda resources for deployment."

            commandLine aws("s3", "rb", "--force", "s3://${deployment.s3Bucket}")
        }

        project.task("deleteStack", type: Exec) {
            group "deployment"
            description "Delete the CloudFormation stack for this project."

            commandLine aws("cloudformation", "delete-stack", "--stack-name", deployment.stackName)

            doLast {
                deployment.outputs.delete()
            }
        }

        project.task("describeStackEvents", type: Exec) {
            group "debug"
            description "Show a table of the events for the cloudformation stack for debugging"

            commandLine aws("cloudformation", "describe-stack-events",
                    "--stack-name", deployment.stackName,
                    "--query", "StackEvents[*].{Time:Timestamp,Type:ResourceType,Status:ResourceStatus,Reason:ResourceStatusReason}",
                    "--output", "table")
        }

    }
}

/**
 * DSL extension class for the "deployment" block
 */
class DeploymentPluginExtension {
    DeploymentPluginExtension(Project project) {
        this.project = project
        this.aws = extensions.create("aws", Aws)
        this.apis = project.container(Api)
        this.lambdaFunctions = project.container(LambdaFunction)

    }

    private Project project

    Aws aws

    /** Name of the cloudformation stack to deploy **/
    String stackName

    /** Name of the S3 bucket to which to upload deployment artifacts **/
    String s3Bucket

    /** File containing the CloudFormation/SAM template to deploy **/
    File templateFile

    /** Outputs of the deployed CloudFormation stack. Will not be available unless stack is already deployed. **/
    StackOutputs outputs

    /** CloudFormation template parameters to override **/
    Map<String, String> parameters = new HashMap<>()

    /** List of API Gateway APIs to deploy **/
    NamedDomainObjectContainer<Api> apis

    NamedDomainObjectContainer<Api> apis(Closure closure) { apis.configure(closure) }

    /** List of Lambda functions to deploy **/
    NamedDomainObjectContainer<LambdaFunction> lambdaFunctions

    NamedDomainObjectContainer<LambdaFunction> lambdaFunctions(Closure closure) { lambdaFunctions.configure(closure) }

    /** Internal plugin use only: name of the intermediate template file after postprocessing **/
    @PackageScope
    File getProcessedTemplateFile() {
        return project.file("${project.buildDir}/cloudformation/${templateFile.name}.processed.json")
    }

    /** Internal plugin use only: name of the intermediate template file after packaging **/
    @PackageScope
    File getPackagedTemplateFile() {
        return project.file("${project.buildDir}/cloudformation/${templateFile.name}.packaged.json")
    }

    /** Internal plugin use only: name of the file with cached outputs from the deployed stack **/
    @PackageScope
    File getStackOutputsFile() {
        return project.file("${project.buildDir}/cloudformation/${templateFile.name}.outputs.json")
    }
}

/**
 * DSL extension class for the "aws" block in "deployment"
 */
class Aws {
    String command = "aws"
    String profile = "default"
    String region = "us-west-2"
}

/**
 * DSL extension class values in the "apis" block in "deployment"
 */
class Api {
    Api(String name) { this.name = name }
    final String name
    Object swaggerFile
}

/**
 * DSL extension class values in the "lambdaFunctions" block in "deployment"
 */
class LambdaFunction {
    LambdaFunction(String name) { this.name = name }
    final String name
    Object zipFile
}

/**
 * Domain class for caching and parsing CloudFormation stack outputs lazily
 */
class StackOutputs {
    File file
    def outputs = null

    StackOutputs(File file) { this.file = file }

    String getAt(String property) {
        getOutputs()?.getAt(property)
    }

    def getOutputs() {
        if (outputs == null) {
            read()
        }

        return outputs
    }

    private void read() {
        if (file.exists()) {
            def stackInfo = new JsonSlurper().parse(file)
            if (stackInfo.Outputs != null) {
                outputs = stackInfo.Outputs.collectEntries { [(it.OutputKey): (it.OutputValue)] }
            }
        }
    }

    private void write(String contents) {
        file.write(contents)
    }

    void delete() {
        if (file.exists()) {
            file.delete()
        }
    }
}
