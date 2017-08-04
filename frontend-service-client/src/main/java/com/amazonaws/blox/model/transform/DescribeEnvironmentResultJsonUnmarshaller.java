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
 * DescribeEnvironmentResult JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DescribeEnvironmentResultJsonUnmarshaller implements Unmarshaller<DescribeEnvironmentResult, JsonUnmarshallerContext> {

    public DescribeEnvironmentResult unmarshall(JsonUnmarshallerContext context) throws Exception {
        DescribeEnvironmentResult describeEnvironmentResult = new DescribeEnvironmentResult();

        int originalDepth = context.getCurrentDepth();
        String currentParentElement = context.getCurrentParentElement();
        int targetDepth = originalDepth + 1;

        JsonToken token = context.getCurrentToken();
        if (token == null)
            token = context.nextToken();
        if (token == VALUE_NULL) {
            return describeEnvironmentResult;
        }

        while (true) {
            if (token == null)
                break;

            describeEnvironmentResult.setEnvironment(EnvironmentJsonUnmarshaller.getInstance().unmarshall(context));
            token = context.nextToken();
        }

        return describeEnvironmentResult;
    }

    private static DescribeEnvironmentResultJsonUnmarshaller instance;

    public static DescribeEnvironmentResultJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new DescribeEnvironmentResultJsonUnmarshaller();
        return instance;
    }
}
