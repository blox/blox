/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

import com.amazonaws.auth.RequestSigner;
import com.amazonaws.opensdk.protect.auth.RequestSignerAware;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/ListEnvironments" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentsRequest extends com.amazonaws.opensdk.BaseRequest implements Serializable, Cloneable, RequestSignerAware {

    private String cluster;

    private String environmentNamePrefix;

    private String maxResults;

    private String nextToken;

    /**
     * @param cluster
     */

    public void setCluster(String cluster) {
        this.cluster = cluster;
    }

    /**
     * @return
     */

    public String getCluster() {
        return this.cluster;
    }

    /**
     * @param cluster
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsRequest cluster(String cluster) {
        setCluster(cluster);
        return this;
    }

    /**
     * @param environmentNamePrefix
     */

    public void setEnvironmentNamePrefix(String environmentNamePrefix) {
        this.environmentNamePrefix = environmentNamePrefix;
    }

    /**
     * @return
     */

    public String getEnvironmentNamePrefix() {
        return this.environmentNamePrefix;
    }

    /**
     * @param environmentNamePrefix
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsRequest environmentNamePrefix(String environmentNamePrefix) {
        setEnvironmentNamePrefix(environmentNamePrefix);
        return this;
    }

    /**
     * @param maxResults
     */

    public void setMaxResults(String maxResults) {
        this.maxResults = maxResults;
    }

    /**
     * @return
     */

    public String getMaxResults() {
        return this.maxResults;
    }

    /**
     * @param maxResults
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsRequest maxResults(String maxResults) {
        setMaxResults(maxResults);
        return this;
    }

    /**
     * @param nextToken
     */

    public void setNextToken(String nextToken) {
        this.nextToken = nextToken;
    }

    /**
     * @return
     */

    public String getNextToken() {
        return this.nextToken;
    }

    /**
     * @param nextToken
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsRequest nextToken(String nextToken) {
        setNextToken(nextToken);
        return this;
    }

    /**
     * Returns a string representation of this object; useful for testing and debugging.
     *
     * @return A string representation of this object.
     *
     * @see java.lang.Object#toString()
     */
    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("{");
        if (getCluster() != null)
            sb.append("Cluster: ").append(getCluster()).append(",");
        if (getEnvironmentNamePrefix() != null)
            sb.append("EnvironmentNamePrefix: ").append(getEnvironmentNamePrefix()).append(",");
        if (getMaxResults() != null)
            sb.append("MaxResults: ").append(getMaxResults()).append(",");
        if (getNextToken() != null)
            sb.append("NextToken: ").append(getNextToken());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof ListEnvironmentsRequest == false)
            return false;
        ListEnvironmentsRequest other = (ListEnvironmentsRequest) obj;
        if (other.getCluster() == null ^ this.getCluster() == null)
            return false;
        if (other.getCluster() != null && other.getCluster().equals(this.getCluster()) == false)
            return false;
        if (other.getEnvironmentNamePrefix() == null ^ this.getEnvironmentNamePrefix() == null)
            return false;
        if (other.getEnvironmentNamePrefix() != null && other.getEnvironmentNamePrefix().equals(this.getEnvironmentNamePrefix()) == false)
            return false;
        if (other.getMaxResults() == null ^ this.getMaxResults() == null)
            return false;
        if (other.getMaxResults() != null && other.getMaxResults().equals(this.getMaxResults()) == false)
            return false;
        if (other.getNextToken() == null ^ this.getNextToken() == null)
            return false;
        if (other.getNextToken() != null && other.getNextToken().equals(this.getNextToken()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getCluster() == null) ? 0 : getCluster().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentNamePrefix() == null) ? 0 : getEnvironmentNamePrefix().hashCode());
        hashCode = prime * hashCode + ((getMaxResults() == null) ? 0 : getMaxResults().hashCode());
        hashCode = prime * hashCode + ((getNextToken() == null) ? 0 : getNextToken().hashCode());
        return hashCode;
    }

    @Override
    public ListEnvironmentsRequest clone() {
        return (ListEnvironmentsRequest) super.clone();
    }

    @Override
    public Class<? extends RequestSigner> signerType() {
        return com.amazonaws.opensdk.protect.auth.IamRequestSigner.class;
    }

    /**
     * Set the configuration for this request.
     *
     * @param sdkRequestConfig
     *        Request configuration.
     * @return This object for method chaining.
     */
    public ListEnvironmentsRequest sdkRequestConfig(com.amazonaws.opensdk.SdkRequestConfig sdkRequestConfig) {
        super.sdkRequestConfig(sdkRequestConfig);
        return this;
    }

}
