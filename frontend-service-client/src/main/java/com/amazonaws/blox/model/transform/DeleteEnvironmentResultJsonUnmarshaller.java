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
 * DeleteEnvironmentResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DeleteEnvironmentResultJsonUnmarshaller implements Unmarshaller<DeleteEnvironmentResult, JsonUnmarshallerContext> {

    public DeleteEnvironmentResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        DeleteEnvironmentResult deleteEnvironmentResult = new DeleteEnvironmentResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return deleteEnvironmentResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("environment", targetDepth)) {
                    context.nextToken();
                    deleteEnvironmentResult.setEnvironment(EnvironmentJsonUnmarshaller.getInstance().unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return deleteEnvironmentResult;
    }

    private static DeleteEnvironmentResultJsonUnmarshaller instance;

    public static DeleteEnvironmentResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new DeleteEnvironmentResultJsonUnmarshaller();
        return instance;
    }
}
