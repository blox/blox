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
 * UpdateEnvironmentResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class UpdateEnvironmentResultJsonUnmarshaller implements Unmarshaller<UpdateEnvironmentResult, JsonUnmarshallerContext> {

    public UpdateEnvironmentResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        UpdateEnvironmentResult updateEnvironmentResult = new UpdateEnvironmentResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return updateEnvironmentResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("environmentRevisionId", targetDepth)) {
                    context.nextToken();
                    updateEnvironmentResult.setEnvironmentRevisionId(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return updateEnvironmentResult;
    }

    private static UpdateEnvironmentResultJsonUnmarshaller instance;

    public static UpdateEnvironmentResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new UpdateEnvironmentResultJsonUnmarshaller();
        return instance;
    }
}
