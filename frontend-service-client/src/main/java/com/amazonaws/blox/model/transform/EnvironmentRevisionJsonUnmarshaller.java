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
 * EnvironmentRevision JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class EnvironmentRevisionJsonUnmarshaller implements Unmarshaller<EnvironmentRevision, JsonUnmarshallerContext> {

    public EnvironmentRevision unmarshall(JsonUnmarshallerContext context) throws Exception {
        EnvironmentRevision environmentRevision = new EnvironmentRevision();

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
                if (context.testExpression("counts", targetDepth)) {
                    context.nextToken();
                    environmentRevision.setCounts(TaskCountsJsonUnmarshaller.getInstance().unmarshall(context));
                }
                if (context.testExpression("environmentRevisionId", targetDepth)) {
                    context.nextToken();
                    environmentRevision.setEnvironmentRevisionId(context.getUnmarshaller(String.class).unmarshall(context));
                }
                if (context.testExpression("instanceGroup", targetDepth)) {
                    context.nextToken();
                    environmentRevision.setInstanceGroup(InstanceGroupJsonUnmarshaller.getInstance().unmarshall(context));
                }
                if (context.testExpression("taskDefinition", targetDepth)) {
                    context.nextToken();
                    environmentRevision.setTaskDefinition(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return environmentRevision;
    }

    private static EnvironmentRevisionJsonUnmarshaller instance;

    public static EnvironmentRevisionJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new EnvironmentRevisionJsonUnmarshaller();
        return instance;
    }
}
