/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.blox.model.*;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * ListEnvironmentsRequestMarshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class ListEnvironmentsRequestMarshaller {

    private static final MarshallingInfo<String> CLUSTER_BINDING = MarshallingInfo.builder(MarshallingType.STRING).marshallLocation(MarshallLocation.PATH)
            .marshallLocationName("cluster").build();
    private static final MarshallingInfo<String> ENVIRONMENTNAMEPREFIX_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.QUERY_PARAM).marshallLocationName("environmentNamePrefix").build();
    private static final MarshallingInfo<String> MAXRESULTS_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.QUERY_PARAM).marshallLocationName("maxResults").build();
    private static final MarshallingInfo<String> NEXTTOKEN_BINDING = MarshallingInfo.builder(MarshallingType.STRING)
            .marshallLocation(MarshallLocation.QUERY_PARAM).marshallLocationName("nextToken").build();

    private static final ListEnvironmentsRequestMarshaller instance = new ListEnvironmentsRequestMarshaller();

    public static ListEnvironmentsRequestMarshaller getInstance() {
        return instance;
    }

    /**
     * Marshall the given parameter object.
     */
    public void marshall(ListEnvironmentsRequest listEnvironmentsRequest, ProtocolMarshaller protocolMarshaller) {

        if (listEnvironmentsRequest == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            protocolMarshaller.marshall(listEnvironmentsRequest.getCluster(), CLUSTER_BINDING);
            protocolMarshaller.marshall(listEnvironmentsRequest.getEnvironmentNamePrefix(), ENVIRONMENTNAMEPREFIX_BINDING);
            protocolMarshaller.marshall(listEnvironmentsRequest.getMaxResults(), MAXRESULTS_BINDING);
            protocolMarshaller.marshall(listEnvironmentsRequest.getNextToken(), NEXTTOKEN_BINDING);
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
