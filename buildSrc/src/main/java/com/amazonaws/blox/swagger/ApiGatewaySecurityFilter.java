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
import io.swagger.models.auth.ApiKeyAuthDefinition;
import io.swagger.models.auth.In;
import java.util.Collections;
import lombok.Getter;
import lombok.Setter;
import org.gradle.api.tasks.Input;

@Getter
@Setter
/**
 * Apply a Swagger Security Definition to all methods of an API.
 *
 * <p>This defaults to applying AWS Signature Version 4 authentication.
 *
 * <p>TODO Move everything in com.amazonaws.blox.swagger to a separate project
 */
public class ApiGatewaySecurityFilter implements SwaggerFilter {
  private static final String SECURITY_SCHEME_NAME = "defaultAuthorization";

  /**
   * The authentication type for the generated swagger model, defaults to AWS Signature Version 4
   *
   * <p>See the documentation for more details:
   * http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-swagger-extensions-authtype.html
   */
  @Input private String authType = "awsSigv4";

  @Override
  public void apply(Swagger swagger) {
    ApiKeyAuthDefinition authorization = new ApiKeyAuthDefinition("Authorization", In.HEADER);
    authorization.setVendorExtension("x-amazon-apigateway-authtype", authType);
    swagger.securityDefinition(SECURITY_SCHEME_NAME, authorization);

    for (Path path : swagger.getPaths().values()) {
      for (Operation operation : path.getOperations()) {
        operation.addSecurity(SECURITY_SCHEME_NAME, Collections.emptyList());
      }
    }
  }
}
