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
package com.amazonaws.blox.swagger;

import io.swagger.models.Operation;
import io.swagger.models.Path;
import io.swagger.models.Swagger;
import java.util.HashMap;
import java.util.Map;
import lombok.AllArgsConstructor;
import org.gradle.api.tasks.Input;

/**
 * Add API Gateway extensions to all operations of a swagger spec.
 *
 * <p>This will add the necessary x-amazon-apigateway-integration section to all operations in a
 * Swagger spec to correctly proxy all requests to a single lambda function.
 *
 * <p>TODO Move everything in com.amazonaws.blox.swagger to a separate project
 */
@AllArgsConstructor
public class ApiGatewayExtensionsFilter implements SwaggerFilter {
  /**
   * The template for the Lambda function name to proxy to, in Cloudformation's Fn::Sub syntax.
   *
   * <p>A typical format for this is:
   *
   * <p>arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${MyLambdaFunction.Arn}/invocations
   */
  @Input private final String lambdaFunctionArnTemplate;

  public Map<String, Object> defaultExtensions() {
    Map<String, Object> extensions = new HashMap<>();

    extensions.put("passthroughBehavior", "when_no_match");
    extensions.put("httpMethod", "POST");
    extensions.put("type", "aws_proxy");

    extensions.put("uri", sub(lambdaFunctionArnTemplate));

    return extensions;
  }

  @Override
  public void apply(Swagger swagger) {
    Map<String, Object> extensions = defaultExtensions();

    for (Path path : swagger.getPaths().values()) {
      for (Operation operation : path.getOperations()) {
        operation.setVendorExtension("x-amazon-apigateway-integration", extensions);
      }
    }
  }

  /**
   * Create a Fn::Sub node with the given template as contents.
   *
   * <p>This is to support using the CloudFormation Fn::Sub intrinsic function in the swagger
   * definition. Using sub("foo${AWS::Region}") for example, would emit {"Fn::Sub":
   * "foo${AWS::Region}"} in the template.
   */
  Map<String, String> sub(String template) {
    Map<String, String> map = new HashMap<>();
    map.put("Fn::Sub", template);

    return map;
  }
}
