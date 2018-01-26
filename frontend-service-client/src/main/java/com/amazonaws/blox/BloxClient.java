/**

*/
package com.amazonaws.blox;

import java.net.*;
import java.util.*;

import javax.annotation.Generated;

import org.apache.commons.logging.*;

import com.amazonaws.*;
import com.amazonaws.opensdk.*;
import com.amazonaws.opensdk.model.*;
import com.amazonaws.opensdk.protect.model.transform.*;
import com.amazonaws.auth.*;
import com.amazonaws.handlers.*;
import com.amazonaws.http.*;
import com.amazonaws.internal.*;
import com.amazonaws.metrics.*;
import com.amazonaws.regions.*;
import com.amazonaws.transform.*;
import com.amazonaws.util.*;
import com.amazonaws.protocol.json.*;

import com.amazonaws.annotation.ThreadSafe;
import com.amazonaws.client.AwsSyncClientParams;

import com.amazonaws.client.ClientHandler;
import com.amazonaws.client.ClientHandlerParams;
import com.amazonaws.client.ClientExecutionParams;
import com.amazonaws.opensdk.protect.client.SdkClientHandler;
import com.amazonaws.SdkBaseException;

import com.amazonaws.blox.model.*;
import com.amazonaws.blox.model.transform.*;

/**
 * Client for accessing Blox. All service calls made using this client are blocking, and will not return until the
 * service call completes.
 * <p>
 * 
 */
@ThreadSafe
@Generated("com.amazonaws:aws-java-sdk-code-generator")
class BloxClient implements Blox {

    private final ClientHandler clientHandler;

    private static final com.amazonaws.opensdk.protect.protocol.ApiGatewayProtocolFactoryImpl protocolFactory = new com.amazonaws.opensdk.protect.protocol.ApiGatewayProtocolFactoryImpl(
            new JsonClientMetadata().withProtocolVersion("1.1").withSupportsCbor(false).withSupportsIon(false).withContentTypeOverride("application/json")
                    .withBaseServiceExceptionClass(com.amazonaws.blox.model.BloxException.class));

    /**
     * Constructs a new client to invoke service methods on Blox using the specified parameters.
     *
     * <p>
     * All service calls made using this new client object are blocking, and will not return until the service call
     * completes.
     *
     * @param clientParams
     *        Object providing client parameters.
     */
    BloxClient(AwsSyncClientParams clientParams) {
        this.clientHandler = new SdkClientHandler(new ClientHandlerParams().withClientParams(clientParams));
    }

    /**
     * @param createEnvironmentRequest
     * @return Result of the createEnvironment operation returned by the service.
     * @sample Blox.createEnvironment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/createEnvironment" target="_top">AWS
     *      API Documentation</a>
     */
    @Override
    public CreateEnvironmentResult createEnvironment(CreateEnvironmentRequest createEnvironmentRequest) {
        HttpResponseHandler<CreateEnvironmentResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata().withPayloadJson(true)
                .withHasStreamingSuccessResponse(false), new CreateEnvironmentResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<CreateEnvironmentRequest, CreateEnvironmentResult>()
                .withMarshaller(new CreateEnvironmentRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(createEnvironmentRequest));
    }

    /**
     * @param deleteEnvironmentRequest
     * @return Result of the deleteEnvironment operation returned by the service.
     * @sample Blox.deleteEnvironment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/deleteEnvironment" target="_top">AWS
     *      API Documentation</a>
     */
    @Override
    public DeleteEnvironmentResult deleteEnvironment(DeleteEnvironmentRequest deleteEnvironmentRequest) {
        HttpResponseHandler<DeleteEnvironmentResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata().withPayloadJson(true)
                .withHasStreamingSuccessResponse(false), new DeleteEnvironmentResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<DeleteEnvironmentRequest, DeleteEnvironmentResult>()
                .withMarshaller(new DeleteEnvironmentRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(deleteEnvironmentRequest));
    }

    /**
     * @param describeEnvironmentRequest
     * @return Result of the describeEnvironment operation returned by the service.
     * @sample Blox.describeEnvironment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/describeEnvironment"
     *      target="_top">AWS API Documentation</a>
     */
    @Override
    public DescribeEnvironmentResult describeEnvironment(DescribeEnvironmentRequest describeEnvironmentRequest) {
        HttpResponseHandler<DescribeEnvironmentResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata()
                .withPayloadJson(true).withHasStreamingSuccessResponse(false), new DescribeEnvironmentResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<DescribeEnvironmentRequest, DescribeEnvironmentResult>()
                .withMarshaller(new DescribeEnvironmentRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(describeEnvironmentRequest));
    }

    /**
     * @param describeEnvironmentDeploymentRequest
     * @return Result of the describeEnvironmentDeployment operation returned by the service.
     * @sample Blox.describeEnvironmentDeployment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/describeEnvironmentDeployment"
     *      target="_top">AWS API Documentation</a>
     */
    @Override
    public DescribeEnvironmentDeploymentResult describeEnvironmentDeployment(DescribeEnvironmentDeploymentRequest describeEnvironmentDeploymentRequest) {
        HttpResponseHandler<DescribeEnvironmentDeploymentResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata()
                .withPayloadJson(true).withHasStreamingSuccessResponse(false), new DescribeEnvironmentDeploymentResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<DescribeEnvironmentDeploymentRequest, DescribeEnvironmentDeploymentResult>()
                .withMarshaller(new DescribeEnvironmentDeploymentRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(describeEnvironmentDeploymentRequest));
    }

    /**
     * @param describeEnvironmentRevisionRequest
     * @return Result of the describeEnvironmentRevision operation returned by the service.
     * @sample Blox.describeEnvironmentRevision
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/describeEnvironmentRevision"
     *      target="_top">AWS API Documentation</a>
     */
    @Override
    public DescribeEnvironmentRevisionResult describeEnvironmentRevision(DescribeEnvironmentRevisionRequest describeEnvironmentRevisionRequest) {
        HttpResponseHandler<DescribeEnvironmentRevisionResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata()
                .withPayloadJson(true).withHasStreamingSuccessResponse(false), new DescribeEnvironmentRevisionResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<DescribeEnvironmentRevisionRequest, DescribeEnvironmentRevisionResult>()
                .withMarshaller(new DescribeEnvironmentRevisionRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(describeEnvironmentRevisionRequest));
    }

    /**
     * @param listEnvironmentDeploymentsRequest
     * @return Result of the listEnvironmentDeployments operation returned by the service.
     * @sample Blox.listEnvironmentDeployments
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/listEnvironmentDeployments"
     *      target="_top">AWS API Documentation</a>
     */
    @Override
    public ListEnvironmentDeploymentsResult listEnvironmentDeployments(ListEnvironmentDeploymentsRequest listEnvironmentDeploymentsRequest) {
        HttpResponseHandler<ListEnvironmentDeploymentsResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata()
                .withPayloadJson(true).withHasStreamingSuccessResponse(false), new ListEnvironmentDeploymentsResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<ListEnvironmentDeploymentsRequest, ListEnvironmentDeploymentsResult>()
                .withMarshaller(new ListEnvironmentDeploymentsRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(listEnvironmentDeploymentsRequest));
    }

    /**
     * @param listEnvironmentRevisionsRequest
     * @return Result of the listEnvironmentRevisions operation returned by the service.
     * @sample Blox.listEnvironmentRevisions
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/listEnvironmentRevisions"
     *      target="_top">AWS API Documentation</a>
     */
    @Override
    public ListEnvironmentRevisionsResult listEnvironmentRevisions(ListEnvironmentRevisionsRequest listEnvironmentRevisionsRequest) {
        HttpResponseHandler<ListEnvironmentRevisionsResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata()
                .withPayloadJson(true).withHasStreamingSuccessResponse(false), new ListEnvironmentRevisionsResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<ListEnvironmentRevisionsRequest, ListEnvironmentRevisionsResult>()
                .withMarshaller(new ListEnvironmentRevisionsRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(listEnvironmentRevisionsRequest));
    }

    /**
     * @param listEnvironmentsRequest
     * @return Result of the listEnvironments operation returned by the service.
     * @sample Blox.listEnvironments
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/listEnvironments" target="_top">AWS
     *      API Documentation</a>
     */
    @Override
    public ListEnvironmentsResult listEnvironments(ListEnvironmentsRequest listEnvironmentsRequest) {
        HttpResponseHandler<ListEnvironmentsResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata().withPayloadJson(true)
                .withHasStreamingSuccessResponse(false), new ListEnvironmentsResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<ListEnvironmentsRequest, ListEnvironmentsResult>()
                .withMarshaller(new ListEnvironmentsRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(listEnvironmentsRequest));
    }

    /**
     * @param startDeploymentRequest
     * @return Result of the startDeployment operation returned by the service.
     * @sample Blox.startDeployment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/startDeployment" target="_top">AWS
     *      API Documentation</a>
     */
    @Override
    public StartDeploymentResult startDeployment(StartDeploymentRequest startDeploymentRequest) {
        HttpResponseHandler<StartDeploymentResult> responseHandler = protocolFactory.createResponseHandler(new JsonOperationMetadata().withPayloadJson(true)
                .withHasStreamingSuccessResponse(false), new StartDeploymentResultJsonUnmarshaller());

        HttpResponseHandler<SdkBaseException> errorResponseHandler = createErrorResponseHandler();

        return clientHandler.execute(new ClientExecutionParams<StartDeploymentRequest, StartDeploymentResult>()
                .withMarshaller(new StartDeploymentRequestProtocolMarshaller(protocolFactory)).withResponseHandler(responseHandler)
                .withErrorResponseHandler(errorResponseHandler).withInput(startDeploymentRequest));
    }

    /**
     * Create the error response handler for the operation.
     * 
     * @param errorShapeMetadata
     *        Error metadata for the given operation
     * @return Configured error response handler to pass to HTTP layer
     */
    private HttpResponseHandler<SdkBaseException> createErrorResponseHandler(JsonErrorShapeMetadata... errorShapeMetadata) {
        return protocolFactory.createErrorResponseHandler(new JsonErrorResponseMetadata().withErrorShapes(Arrays.asList(errorShapeMetadata)));
    }

    @Override
    public void shutdown() {
        clientHandler.shutdown();
    }

}
