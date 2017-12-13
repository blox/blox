/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/ListEnvironmentDeployments"
 *      target="_top">AWS API Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentDeploymentsResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private java.util.List<String> deploymentIds;

    private String nextToken;

    /**
     * @return
     */

    public java.util.List<String> getDeploymentIds() {
        return deploymentIds;
    }

    /**
     * @param deploymentIds
     */

    public void setDeploymentIds(java.util.Collection<String> deploymentIds) {
        if (deploymentIds == null) {
            this.deploymentIds = null;
            return;
        }

        this.deploymentIds = new java.util.ArrayList<String>(deploymentIds);
    }

    /**
     * <p>
     * <b>NOTE:</b> This method appends the values to the existing list (if any). Use
     * {@link #setDeploymentIds(java.util.Collection)} or {@link #withDeploymentIds(java.util.Collection)} if you want
     * to override the existing values.
     * </p>
     * 
     * @param deploymentIds
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentDeploymentsResult deploymentIds(String... deploymentIds) {
        if (this.deploymentIds == null) {
            setDeploymentIds(new java.util.ArrayList<String>(deploymentIds.length));
        }
        for (String ele : deploymentIds) {
            this.deploymentIds.add(ele);
        }
        return this;
    }

    /**
     * @param deploymentIds
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentDeploymentsResult deploymentIds(java.util.Collection<String> deploymentIds) {
        setDeploymentIds(deploymentIds);
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

    public ListEnvironmentDeploymentsResult nextToken(String nextToken) {
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
        if (getDeploymentIds() != null)
            sb.append("DeploymentIds: ").append(getDeploymentIds()).append(",");
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

        if (obj instanceof ListEnvironmentDeploymentsResult == false)
            return false;
        ListEnvironmentDeploymentsResult other = (ListEnvironmentDeploymentsResult) obj;
        if (other.getDeploymentIds() == null ^ this.getDeploymentIds() == null)
            return false;
        if (other.getDeploymentIds() != null && other.getDeploymentIds().equals(this.getDeploymentIds()) == false)
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

        hashCode = prime * hashCode + ((getDeploymentIds() == null) ? 0 : getDeploymentIds().hashCode());
        hashCode = prime * hashCode + ((getNextToken() == null) ? 0 : getNextToken().hashCode());
        return hashCode;
    }

    @Override
    public ListEnvironmentDeploymentsResult clone() {
        try {
            return (ListEnvironmentDeploymentsResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
