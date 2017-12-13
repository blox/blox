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
 * DescribeEnvironmentDeploymentResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DescribeEnvironmentDeploymentResultJsonUnmarshaller implements Unmarshaller<DescribeEnvironmentDeploymentResult, JsonUnmarshallerContext> {

    public DescribeEnvironmentDeploymentResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        DescribeEnvironmentDeploymentResult describeEnvironmentDeploymentResult = new DescribeEnvironmentDeploymentResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return describeEnvironmentDeploymentResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("deployment", targetDepth)) {
                    context.nextToken();
                    describeEnvironmentDeploymentResult.setDeployment(DeploymentJsonUnmarshaller.getInstance().unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return describeEnvironmentDeploymentResult;
    }

    private static DescribeEnvironmentDeploymentResultJsonUnmarshaller instance;

    public static DescribeEnvironmentDeploymentResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new DescribeEnvironmentDeploymentResultJsonUnmarshaller();
        return instance;
    }
}
