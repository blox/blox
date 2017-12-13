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
 * ListEnvironmentsResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentsResultJsonUnmarshaller implements Unmarshaller<ListEnvironmentsResult, JsonUnmarshallerContext> {

    public ListEnvironmentsResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        ListEnvironmentsResult listEnvironmentsResult = new ListEnvironmentsResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return listEnvironmentsResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("environmentNames", targetDepth)) {
                    context.nextToken();
                    listEnvironmentsResult.setEnvironmentNames(new ListUnmarshaller<String>(context.getUnmarshaller(String.class)).unmarshall(context));
                }
                if (context.testExpression("nextToken", targetDepth)) {
                    context.nextToken();
                    listEnvironmentsResult.setNextToken(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return listEnvironmentsResult;
    }

    private static ListEnvironmentsResultJsonUnmarshaller instance;

    public static ListEnvironmentsResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new ListEnvironmentsResultJsonUnmarshaller();
        return instance;
    }
}
