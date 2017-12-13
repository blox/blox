/**

*/
package com.amazonaws.blox.model.transform;

import java.util.Map;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * EnvironmentRevisionMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class EnvironmentRevisionMarshaller {

    private static final MarshallingInfo<String> CLUSTER_BINDING = MarshallingInfo.builder(MarshallingType.STRING).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("cluster").build();
    private static final MarshallingInfo<StructuredPojo> COUNTS_BINDING = MarshallingInfo.builder(MarshallingType.STRUCTURED)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("counts").build();
    private static final MarshallingInfo<Map> DEPLOYMENTCONFIGURATION_BINDING = MarshallingInfo.builder(MarshallingType.MAP)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentConfiguration").build();
    private static final MarshallingInfo<String> DEPLOYMENTMETHOD_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentMethod").build();
    private static final MarshallingInfo<String> ENVIRONMENTNAME_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentName").build();
    private static final MarshallingInfo<String> ENVIRONMENTREVISIONID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentRevisionId").build();
    private static final MarshallingInfo<String> ENVIRONMENTTYPE_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentType").build();
    private static final MarshallingInfo<StructuredPojo> INSTANCEGROUP_BINDING = MarshallingInfo.builder(MarshallingType.STRUCTURED)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("instanceGroup").build();
    private static final MarshallingInfo<String> TASKDEFINITION_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("taskDefinition").build();

    private static final EnvironmentRevisionMarshaller instance = new EnvironmentRevisionMarshaller();

    public static EnvironmentRevisionMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(EnvironmentRevision environmentRevision, ProtocolMarshaller protocolMarshaller) {

        if (environmentRevision == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(environmentRevision.getCluster(), CLUSTER_BINDING);
            protocolMarshaller.marshall(environmentRevision.getCounts(), COUNTS_BINDING);
            protocolMarshaller.marshall(environmentRevision.getDeploymentConfiguration(), DEPLOYMENTCONFIGURATION_BINDING);
            protocolMarshaller.marshall(environmentRevision.getDeploymentMethod(), DEPLOYMENTMETHOD_BINDING);
            protocolMarshaller.marshall(environmentRevision.getEnvironmentName(), ENVIRONMENTNAME_BINDING);
            protocolMarshaller.marshall(environmentRevision.getEnvironmentRevisionId(), ENVIRONMENTREVISIONID_BINDING);
            protocolMarshaller.marshall(environmentRevision.getEnvironmentType(), ENVIRONMENTTYPE_BINDING);
            protocolMarshaller.marshall(environmentRevision.getInstanceGroup(), INSTANCEGROUP_BINDING);
            protocolMarshaller.marshall(environmentRevision.getTaskDefinition(), TASKDEFINITION_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
