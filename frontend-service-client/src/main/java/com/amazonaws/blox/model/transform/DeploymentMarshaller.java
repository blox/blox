/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * DeploymentMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class DeploymentMarshaller {

    private static final MarshallingInfo<String> DEPLOYMENTID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("deploymentId").build();
    private static final MarshallingInfo<String> ENVIRONMENTNAME_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("environmentName").build();
    private static final MarshallingInfo<String> NEWTARGETREVISIONID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("newTargetRevisionId").build();
    private static final MarshallingInfo<String> OLDTARGETREVISIONID_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("oldTargetRevisionId").build();
    private static final MarshallingInfo<Integer> TIMESTAMP_BINDING = MarshallingInfo.builder(MarshallingType.INTEGER)
            .marshallLocation(MarshallLocation.PAYLOAD).marshallLocationName("timestamp").build();

    private static final DeploymentMarshaller instance = new DeploymentMarshaller();

    public static DeploymentMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(Deployment deployment, ProtocolMarshaller protocolMarshaller) {

        if (deployment == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(deployment.getDeploymentId(), DEPLOYMENTID_BINDING);
            protocolMarshaller.marshall(deployment.getEnvironmentName(), ENVIRONMENTNAME_BINDING);
            protocolMarshaller.marshall(deployment.getNewTargetRevisionId(), NEWTARGETREVISIONID_BINDING);
            protocolMarshaller.marshall(deployment.getOldTargetRevisionId(), OLDTARGETREVISIONID_BINDING);
            protocolMarshaller.marshall(deployment.getTimestamp(), TIMESTAMP_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
