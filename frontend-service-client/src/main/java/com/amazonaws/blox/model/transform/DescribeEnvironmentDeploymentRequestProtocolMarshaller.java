/**

*/
package com.amazonaws.blox.model.transform;

import javax.annotation.Generated;

import com.amazonaws.SdkClientException;
import com.amazonaws.Request;

import com.amazonaws.http.HttpMethodName;
import com.amazonaws.blox.model.*;
import com.amazonaws.transform.Marshaller;

import com.amazonaws.protocol.*;
import com.amazonaws.annotation.SdkInternalApi;

/**
 * DescribeEnvironmentDeploymentRequest Marshaller
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
@SdkInternalApi
public class DescribeEnvironmentDeploymentRequestProtocolMarshaller implements
        Marshaller<Request<DescribeEnvironmentDeploymentRequest>, DescribeEnvironmentDeploymentRequest> {

    private static final OperationInfo SDK_OPERATION_BINDING = OperationInfo.builder().protocol(Protocol.API_GATEWAY)
            .requestUri("/alpha/v1/{cluster}/environments/{environmentName}/deployments/{deploymentId}").httpMethodName(HttpMethodName.GET)
            .hasExplicitPayloadMember(false).hasPayloadMembers(false).serviceName("Blox").build();

    private final com.amazonaws.opensdk.protect.protocol.ApiGatewayProtocolFactoryImpl protocolFactory;

    public DescribeEnvironmentDeploymentRequestProtocolMarshaller(com.amazonaws.opensdk.protect.protocol.ApiGatewayProtocolFactoryImpl protocolFactory) {
        this.protocolFactory = protocolFactory;
    }

    public Request<DescribeEnvironmentDeploymentRequest> marshall(DescribeEnvironmentDeploymentRequest describeEnvironmentDeploymentRequest) {

        if (describeEnvironmentDeploymentRequest == null) {
            throw new SdkClientException("Invalid argument passed to marshall(...)");
        }

        try {
            final ProtocolRequestMarshaller<DescribeEnvironmentDeploymentRequest> protocolMarshaller = protocolFactory.createProtocolMarshaller(
                    SDK_OPERATION_BINDING, describeEnvironmentDeploymentRequest);

            protocolMarshaller.startMarshalling();
            DescribeEnvironmentDeploymentRequestMarshaller.getInstance().marshall(describeEnvironmentDeploymentRequest, protocolMarshaller);
            return protocolMarshaller.finishMarshalling();
        } catch (Exception e) {
            throw new SdkClientException("Unable to marshall request to JSON: " + e.getMessage(), e);
        }
    }

}
