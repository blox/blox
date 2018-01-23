/**

*/
package com.amazonaws.blox.model.transform;

import java.math.*;

import javax.annotation.Generated;

import com.amazonaws.blox.model.*;
import com.amazonaws.transform.SimpleTypeJsonUnmarshallers.*;
import com.amazonaws.transform.*;

import com.fasterxml.jackson.core.JsonToken;
import static com.fasterxml.jackson.core.JsonToken.*;

/**
 * Environment JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class EnvironmentJsonUnmarshaller implements Unmarshaller<Environment, JsonUnmarshallerContext> {

    public Environment unmarshall(JsonUnmarshallerContext context) throws Exception {
        Environment environment = new Environment();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return null;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("activeEnvironmentRevisionId", targetDepth)) {
                    context.nextToken();
                    environment.setActiveEnvironmentRevisionId(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("cluster", targetDepth)) {
                    context.nextToken();
                    environment.setCluster(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("deploymentConfiguration", targetDepth)) {
                    context.nextToken();
                    environment.setDeploymentConfiguration(DeploymentConfigurationJsonUnmarshaller.getInstance().unmarshall(context));
                }
                if (context.testExpression("deploymentMethod", targetDepth)) {
                    context.nextToken();
                    environment.setDeploymentMethod(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("environmentHealth", targetDepth)) {
                    context.nextToken();
                    environment.setEnvironmentHealth(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("environmentName", targetDepth)) {
                    context.nextToken();
                    environment.setEnvironmentName(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("environmentType", targetDepth)) {
                    context.nextToken();
                    environment.setEnvironmentType(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("latestEnvironmentRevisionId", targetDepth)) {
                    context.nextToken();
                    environment.setLatestEnvironmentRevisionId(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("role", targetDepth)) {
                    context.nextToken();
                    environment.setRole(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return environment;
    }

    private static EnvironmentJsonUnmarshaller instance;

    public static EnvironmentJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new EnvironmentJsonUnmarshaller();
        return instance;
    }
}
