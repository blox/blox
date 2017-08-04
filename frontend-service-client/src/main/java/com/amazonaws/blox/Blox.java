/**

*/
package com.amazonaws.blox;

import javax.annotation.Generated;

import com.amazonaws.*;
import com.amazonaws.opensdk.*;
import com.amazonaws.opensdk.model.*;
import com.amazonaws.regions.*;

import com.amazonaws.blox.model.*;

/**
 * Interface for accessing Blox.
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public interface Blox {

    /**
     * @param describeEnvironmentRequest
     * @return Result of the describeEnvironment operation returned by the service.
     * @sample Blox.describeEnvironment
     * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/joufu8ief9-v2017-07-11/describeEnvironment"
     *      target="_top">AWS API Documentation</a>
     */
    DescribeEnvironmentResult describeEnvironment(DescribeEnvironmentRequest describeEnvironmentRequest);

    /**
     * @return Create new instance of builder with all defaults set.
     */
    public static BloxClientBuilder builder() {
        return new BloxClientBuilder();
    }

    /**
     * Shuts down this client object, releasing any resources that might be held open. This is an optional method, and
     * callers are not expected to call it, but can if they want to explicitly release any open resources. Once a client
     * has been shutdown, it should not be used to make any more requests.
     */
    void shutdown();

}
