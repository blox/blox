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
 * InstanceGroup JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class InstanceGroupJsonUnmarshaller implements Unmarshaller<InstanceGroup, JsonUnmarshallerContext> {

    public InstanceGroup unmarshall(JsonUnmarshallerContext context) throws Exception {
        InstanceGroup instanceGroup = new InstanceGroup();

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
                if (context.testExpression("query", targetDepth)) {
                    context.nextToken();
                    instanceGroup.setQuery(context.getUnmarshaller(String.class).unmarshall(context));
                }
            } else if (token == END_ARRAY || token == END_OBJECT) {
                if (context.getLastParsedParentElement() == null || context.getLastParsedParentElement().equals(currentParentElement)) {
                    if (context.getCurrentDepth() <= originalDepth)
                        break;
                }
            }
            token = context.nextToken();
        }

        return instanceGroup;
    }

    private static InstanceGroupJsonUnmarshaller instance;

    public static InstanceGroupJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new InstanceGroupJsonUnmarshaller();
        return instance;
    }
}
