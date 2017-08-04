/*
 * Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"). You may
 * not use this file except in compliance with the License. A copy of the
 * License is located at
 *
 *     http://aws.amazon.com/apache2.0/
 *
 * or in the "LICENSE" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */
package com.amazonaws.blox.tasks;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.util.DefaultPrettyPrinter;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.fasterxml.jackson.databind.node.TextNode;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import java.io.File;
import java.io.IOException;
import lombok.Getter;
import lombok.Setter;
import org.gradle.api.DefaultTask;
import org.gradle.api.tasks.Input;
import org.gradle.api.tasks.InputFile;
import org.gradle.api.tasks.OutputFile;
import org.gradle.api.tasks.TaskAction;

/**
 * Substitute dynamic values for Serverless::Function.CodeURI and Serverless::Api.DefinitionBody
 * into a given Cloudformation/Serverless Application Model template.
 *
 * <p>Defining an API gateway API using Swagger using SAM has some limitations when using a separate
 * swagger definition file (as declared through the DefinitionUri property): - We must hard-code the
 * region and account ID of the arn of the lambda function it should invoke in the
 * `x-amazon-apigateway-integration/uri` property. - If we don't hardcode the Lambda function name,
 * we must provide it as a stage variable instead.
 *
 * <p>Instead, SAM supports specifying the Swagger definition inline using the DefinitionBody field.
 * This allows us to use CloudFormation functions and references to derive the name of the Lambda
 * function to invoke, but brings its own set of limitations: - Swagger YAML definitions almost
 * always have at least one numeric field name (in the responses object), and CloudFormation
 * incorrectly reads this as an integer instead of a string key.
 *
 * <p>The solution then is to embed the Swagger definition into the SAM Template from YAML source,
 * but to produce the combined template in JSON format to avoid the non-string key issue.
 *
 * <p>In the same template, cloudformation only supports paths relative to the template file or full
 * absolute paths for the CodeUri property of a Serverless::Function. This task will substitute in
 * the given lambdaZipFile as an absolute path into the CodeUri property.
 */
@Getter
@Setter
public class PostProcessCloudformation extends DefaultTask {
  /** Mapper used to read the source swagger/SAM YAML templates */
  private ObjectMapper inputMapper = new ObjectMapper(new YAMLFactory());
  /** Mapper used to produce combined JSON template */
  private ObjectMapper outputMapper = new ObjectMapper(new JsonFactory());

  /**
   * The name of the API GW resource in the SAM template into which the swagger template should be
   * embedded.
   */
  @Input private String apiName;

  /**
   * The name of the Lambda Function resource in the SAM template into which the Lambda zip path
   * should be embedded.
   */
  @Input private String handlerName;

  /** The YAML file to read the swagger definition from. */
  @InputFile private File swaggerFile;

  /** The YAML file to read the SAM template from. */
  @InputFile private File templateFile;

  /** The JSON file to write the combined template into. */
  @OutputFile private File outputTemplateFile;

  /** The Zip file to substitute into the CodeUri property of the lambda function */
  @InputFile private File lambdaZipFile;

  @TaskAction
  public void embedSwagger() throws IOException {
    JsonNode template = inputMapper.readTree(templateFile);

    JsonNode resources = template.get("Resources");
    embedSwaggerBody(resources);
    setCodeUri(resources);

    ObjectWriter writer = outputMapper.writer(new DefaultPrettyPrinter());
    outputTemplateFile.getParentFile().mkdirs();
    writer.writeValue(outputTemplateFile, template);
  }

  private void setCodeUri(JsonNode resources) {
    TextNode codeUri = new TextNode(lambdaZipFile.getAbsolutePath());

    ObjectNode apiProperties = (ObjectNode) resources.get(handlerName).get("Properties");
    apiProperties.set("CodeUri", codeUri);
  }

  private void embedSwaggerBody(JsonNode resources) throws IOException {
    JsonNode swagger = inputMapper.readTree(swaggerFile);

    ObjectNode properties = (ObjectNode) resources.get(apiName).get("Properties");
    properties.set("DefinitionBody", swagger);
  }
}
