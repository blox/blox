/**

*/
package com.amazonaws.blox.model.transform;

import java.math.*;

import javax.annotation.Generated;

import com.amazonaws.blox.model.*;
import com.amazonaws.transform.SimpleTypeJsonUnmarshallers.*;
import com.amazonaws.transform.*;

import static com.fasterxml.jackson.core.JsonToken.*;

/**
 * DeploymentConfiguration JSON Unmarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DeploymentConfigurationJsonUnmarshaller implements Unmarshaller<DeploymentConfiguration, JsonUnmarshallerContext> {

    public DeploymentConfiguration unmarshall(JsonUnmarshallerContext context) throws Exception {
        DeploymentConfiguration deploymentConfiguration = new DeploymentConfiguration();

        return deploymentConfiguration;
    }

    private static DeploymentConfigurationJsonUnmarshaller instance;

    public static DeploymentConfigurationJsonUnmarshaller getInstance() {
        if (instance == null)
            instance = new DeploymentConfigurationJsonUnmarshaller();
        return instance;
    }
}
