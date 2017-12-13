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
 * StartDeploymentResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class StartDeploymentResultJsonUnmarshaller implements Unmarshaller<StartDeploymentResult, JsonUnmarshallerContext> {

    public StartDeploymentResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        StartDeploymentResult startDeploymentResult = new StartDeploymentResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return startDeploymentResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("deploymentId", targetDepth)) {
                    context.nextToken();
                    startDeploymentResult.setDeploymentId(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return startDeploymentResult;
    }

    private static StartDeploymentResultJsonUnmarshaller instance;

    public static StartDeploymentResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new StartDeploymentResultJsonUnmarshaller();
        return instance;
    }
}
