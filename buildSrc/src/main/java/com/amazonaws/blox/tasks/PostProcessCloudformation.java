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
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;
import javafx.util.Pair;
import lombok.Getter;
import lombok.Setter;
import org.gradle.api.DefaultTask;
import org.gradle.api.file.FileCollection;
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
 * This allows us to use Cloudformation functions and references to derive the name of the Lambda
 * function to invoke, but brings its own set of limitations: - Swagger YAML definitions almost
 * always have at least one numeric field name (in the responses object), and Cloudformation
 * incorrectly reads this as an integer instead of a string key.
 *
 * <p>The solution then is to embed the Swagger definition into the SAM Template from YAML source,
 * but to produce the combined template in JSON format to avoid the non-string key issue.
 *
 * <p>In the same template, Cloudformation only supports paths relative to the template file or full
 * absolute paths for the CodeUri property of a Serverless::Function. This task will substitute in
 * the given lambdaZipFile as an absolute path into the CodeUri property.
 *
 * <p>TODO The lambda/swagger stuff is starting to feel like separate concerns. Using a Filter model
 * similar to the swagger task might make things a bit cleaner
 */
public class PostProcessCloudformation extends DefaultTask {

  /** Mapper used to read the source swagger/SAM YAML templates */
  private ObjectMapper inputMapper = new ObjectMapper(new YAMLFactory());
  /** Mapper used to produce combined JSON template */
  private ObjectMapper outputMapper = new ObjectMapper(new JsonFactory());
  /** The YAML file to read the SAM template from. */
  @Getter @Setter @InputFile private File templateFile;
  /** The JSON file to write the combined template into. */
  @Getter @Setter @OutputFile private File outputTemplateFile;

  /**
   * The lambda functions that should have their CodeUri property updated in the Cloudformation
   * template.
   */
  private Map<String, File> lambdaFunctions = new HashMap<>();

  /**
   * The API Gateway APIs that should have their DefinitionBody updated in the Cloudformation
   * template.
   */
  private Map<String, File> apis = new HashMap<>();

  /**
   * Register a lambda function that should have its CodeUri property updated in the Cloudformation
   * template
   *
   * @param ref the CloudFormation Ref of the function to update
   * @param zip the files containing the function's deployment package
   */
  public void lambdaFunction(String ref, FileCollection zip) {
    getInputs().files(zip);
    lambdaFunctions.put(ref, zip.getSingleFile());
  }

  /**
   * Register an API Gateway API that should have its DefinitionBody property updated in the
   * Cloudformation template
   *
   * @param ref the Cloudformation Ref of the API to update
   * @param swagger the file containing the swagger definition of the API
   */
  public void api(String ref, FileCollection swagger) {
    getInputs().files(swagger);
    apis.put(ref, swagger.getSingleFile());
  }

  @TaskAction
  public void embedSwagger() throws IOException {
    JsonNode template = inputMapper.readTree(templateFile);

    JsonNode resources = template.get("Resources");
    addApiDefinitionBody(resources);
    addLambdaCodeURI(resources);

    ObjectWriter writer = outputMapper.writer(new DefaultPrettyPrinter());
    outputTemplateFile.getParentFile().mkdirs();
    writer.writeValue(outputTemplateFile, template);
  }

  private void addLambdaCodeURI(JsonNode resources) {
    for (Entry<String, File> function : lambdaFunctions.entrySet()) {
      TextNode codeUri = new TextNode(function.getValue().getAbsolutePath());

      ObjectNode apiProperties = (ObjectNode) resources.get(function.getKey()).get("Properties");
      apiProperties.set("CodeUri", codeUri);
    }
  }

  private void addApiDefinitionBody(JsonNode resources) throws IOException {
    for (Entry<String, File> api : apis.entrySet()) {
      JsonNode swagger = inputMapper.readTree(api.getValue());

      ObjectNode properties = (ObjectNode) resources.get(api.getKey()).get("Properties");
      properties.set("DefinitionBody", swagger);
    }
  }
}
