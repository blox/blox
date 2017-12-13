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
 * ListEnvironmentRevisionsResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentRevisionsResultJsonUnmarshaller implements Unmarshaller<ListEnvironmentRevisionsResult, JsonUnmarshallerContext> {

    public ListEnvironmentRevisionsResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        ListEnvironmentRevisionsResult listEnvironmentRevisionsResult = new ListEnvironmentRevisionsResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return listEnvironmentRevisionsResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("revisionIds", targetDepth)) {
                    context.nextToken();
                    listEnvironmentRevisionsResult.setRevisionIds(new ListUnmarshaller<String>(context.getUnmarshaller(String.class)).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return listEnvironmentRevisionsResult;
    }

    private static ListEnvironmentRevisionsResultJsonUnmarshaller instance;

    public static ListEnvironmentRevisionsResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new ListEnvironmentRevisionsResultJsonUnmarshaller();
        return instance;
    }
}
