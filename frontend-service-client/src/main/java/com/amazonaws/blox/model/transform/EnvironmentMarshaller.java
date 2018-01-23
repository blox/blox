/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * EnvironmentMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class EnvironmentMarshaller {

    private static final MarshallingInfo<String> ACTIVEENVIRONMENTREVISIONID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("activeEnvironmentRevisionId").build();
    private static final MarshallingInfo<String> CLUSTER_BINDING = MarshallingInfo.builder(MarshallingType.STRING).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("cluster").build();
    private static final MarshallingInfo<StructuredPojo> DEPLOYMENTCONFIGURATION_BINDING = MarshallingInfo.builder(MarshallingType.STRUCTURED)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentConfiguration").build();
    private static final MarshallingInfo<String> DEPLOYMENTMETHOD_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentMethod").build();
    private static final MarshallingInfo<String> ENVIRONMENTHEALTH_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentHealth").build();
    private static final MarshallingInfo<String> ENVIRONMENTNAME_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentName").build();
    private static final MarshallingInfo<String> ENVIRONMENTTYPE_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentType").build();
    private static final MarshallingInfo<String> LATESTENVIRONMENTREVISIONID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("latestEnvironmentRevisionId").build();
    private static final MarshallingInfo<String> ROLE_BINDING = MarshallingInfo.builder(MarshallingType.STRING).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("role").build();

    private static final EnvironmentMarshaller instance = new EnvironmentMarshaller();

    public static EnvironmentMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(Environment environment, ProtocolMarshaller protocolMarshaller) {

        if (environment == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(environment.getActiveEnvironmentRevisionId(), ACTIVEENVIRONMENTREVISIONID_BINDING);
            protocolMarshaller.marshall(environment.getCluster(), CLUSTER_BINDING);
            protocolMarshaller.marshall(environment.getDeploymentConfiguration(), DEPLOYMENTCONFIGURATION_BINDING);
            protocolMarshaller.marshall(environment.getDeploymentMethod(), DEPLOYMENTMETHOD_BINDING);
            protocolMarshaller.marshall(environment.getEnvironmentHealth(), ENVIRONMENTHEALTH_BINDING);
            protocolMarshaller.marshall(environment.getEnvironmentName(), ENVIRONMENTNAME_BINDING);
            protocolMarshaller.marshall(environment.getEnvironmentType(), ENVIRONMENTTYPE_BINDING);
            protocolMarshaller.marshall(environment.getLatestEnvironmentRevisionId(), LATESTENVIRONMENTREVISIONID_BINDING);
            protocolMarshaller.marshall(environment.getRole(), ROLE_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
