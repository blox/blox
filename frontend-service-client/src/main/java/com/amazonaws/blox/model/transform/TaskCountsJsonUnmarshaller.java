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
 * TaskCounts JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class TaskCountsJsonUnmarshaller implements Unmarshaller<TaskCounts, JsonUnmarshallerContext> {

    public TaskCounts unmarshall(JsonUnmarshallerContext context) throws Exception {
        TaskCounts taskCounts = new TaskCounts();

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
                if (context.testExpression("desired", targetDepth)) {
                    context.nextToken();
                    taskCounts.setDesired(context.getUnmarshaller(Integer.class).unmarshall(context));
                }
                if (context.testExpression("healthy", targetDepth)) {
                    context.nextToken();
                    taskCounts.setHealthy(context.getUnmarshaller(Integer.class).unmarshall(context));
                }
                if (context.testExpression("total", targetDepth)) {
                    context.nextToken();
                    taskCounts.setTotal(context.getUnmarshaller(Integer.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return taskCounts;
    }

    private static TaskCountsJsonUnmarshaller instance;

    public static TaskCountsJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new TaskCountsJsonUnmarshaller();
        return instance;
    }
}
