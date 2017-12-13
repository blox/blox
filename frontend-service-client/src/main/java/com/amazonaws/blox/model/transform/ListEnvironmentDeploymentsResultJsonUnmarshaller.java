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
 * ListEnvironmentDeploymentsResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentDeploymentsResultJsonUnmarshaller implements Unmarshaller<ListEnvironmentDeploymentsResult, JsonUnmarshallerContext> {

    public ListEnvironmentDeploymentsResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        ListEnvironmentDeploymentsResult listEnvironmentDeploymentsResult = new ListEnvironmentDeploymentsResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return listEnvironmentDeploymentsResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("deploymentIds", targetDepth)) {
                    context.nextToken();
                    listEnvironmentDeploymentsResult.setDeploymentIds(new ListUnmarshaller<String>(context.getUnmarshaller(String.class)).unmarshall(context));
                }
                if (context.testExpression("nextToken", targetDepth)) {
                    context.nextToken();
                    listEnvironmentDeploymentsResult.setNextToken(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return listEnvironmentDeploymentsResult;
    }

    private static ListEnvironmentDeploymentsResultJsonUnmarshaller instance;

    public static ListEnvironmentDeploymentsResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new ListEnvironmentDeploymentsResultJsonUnmarshaller();
        return instance;
    }
}
