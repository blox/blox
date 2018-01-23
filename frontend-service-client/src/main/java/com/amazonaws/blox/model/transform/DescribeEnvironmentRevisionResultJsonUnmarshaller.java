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
 * DescribeEnvironmentRevisionResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DescribeEnvironmentRevisionResultJsonUnmarshaller implements Unmarshaller<DescribeEnvironmentRevisionResult, JsonUnmarshallerContext> {

    public DescribeEnvironmentRevisionResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        DescribeEnvironmentRevisionResult describeEnvironmentRevisionResult = new DescribeEnvironmentRevisionResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return describeEnvironmentRevisionResult;
        }

        while (true) {
            if (token == null)
                break;

            if (token == FIELD_NAME || token == START_OBJECT) {
                if (context.testExpression("environmentRevision", targetDepth)) {
                    context.nextToken();
                    describeEnvironmentRevisionResult.setEnvironmentRevision(EnvironmentRevisionJsonUnmarshaller.getInstance().unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return describeEnvironmentRevisionResult;
    }

    private static DescribeEnvironmentRevisionResultJsonUnmarshaller instance;

    public static DescribeEnvironmentRevisionResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new DescribeEnvironmentRevisionResultJsonUnmarshaller();
        return instance;
    }
}
