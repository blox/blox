/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * CreateEnvironmentRequestMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class CreateEnvironmentRequestMarshaller {

    private static final MarshallingInfo<String> CLUSTER_BINDING = MarshallingInfo.builder(MarshallingType.STRING).marshallLocation(MarshallLocation.PATH)
            .marshallLocationName("cluster").build();
    private static final MarshallingInfo<StructuredPojo> DEPLOYMENTCONFIGURATION_BINDING = MarshallingInfo.builder(MarshallingType.STRUCTURED)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentConfiguration").build();
    private static final MarshallingInfo<String> DEPLOYMENTMETHOD_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentMethod").build();
    private static final MarshallingInfo<String> DEPLOYMENTTYPE_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentType").build();
    private static final MarshallingInfo<String> ENVIRONMENTNAME_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentName").build();
    private static final MarshallingInfo<StructuredPojo> INSTANCEGROUP_BINDING = MarshallingInfo.builder(MarshallingType.STRUCTURED)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("instanceGroup").build();
    private static final MarshallingInfo<String> TASKDEFINITION_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("taskDefinition").build();

    private static final CreateEnvironmentRequestMarshaller instance = new CreateEnvironmentRequestMarshaller();

    public static CreateEnvironmentRequestMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(CreateEnvironmentRequest createEnvironmentRequest, ProtocolMarshaller protocolMarshaller) {

        if (createEnvironmentRequest == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(createEnvironmentRequest.getCluster(), CLUSTER_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getDeploymentConfiguration(), DEPLOYMENTCONFIGURATION_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getDeploymentMethod(), DEPLOYMENTMETHOD_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getDeploymentType(), DEPLOYMENTTYPE_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getEnvironmentName(), ENVIRONMENTNAME_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getInstanceGroup(), INSTANCEGROUP_BINDING);
            protocolMarshaller.marshall(createEnvironmentRequest.getTaskDefinition(), TASKDEFINITION_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
