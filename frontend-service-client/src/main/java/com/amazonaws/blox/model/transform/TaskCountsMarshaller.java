/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * TaskCountsMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class TaskCountsMarshaller {

    private static final MarshallingInfo<Integer> DESIRED_BINDING = MarshallingInfo.builder(MarshallingType.INTEGER).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("desired").build();
    private static final MarshallingInfo<Integer> HEALTHY_BINDING = MarshallingInfo.builder(MarshallingType.INTEGER).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("healthy").build();
    private static final MarshallingInfo<Integer> TOTAL_BINDING = MarshallingInfo.builder(MarshallingType.INTEGER).marshallLocation(MarshallLocation.PAYLOAD)
            .marshallLocationName("total").build();

    private static final TaskCountsMarshaller instance = new TaskCountsMarshaller();

    public static TaskCountsMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(TaskCounts taskCounts, ProtocolMarshaller protocolMarshaller) {

        if (taskCounts == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(taskCounts.getDesired(), DESIRED_BINDING);
            protocolMarshaller.marshall(taskCounts.getHealthy(), HEALTHY_BINDING);
            protocolMarshaller.marshall(taskCounts.getTotal(), TOTAL_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
