/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/DescribeEnvironmentDeployment"
 *      target="_top">AWS API Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DescribeEnvironmentDeploymentResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private Deployment deployment;

    /**
     * @param deployment
     */

    public void setDeployment(Deployment deployment) {
        this.deployment = deployment;
    }

    /**
     * @return
     */

    public Deployment getDeployment() {
        return this.deployment;
    }

    /**
     * @param deployment
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public DescribeEnvironmentDeploymentResult deployment(Deployment deployment) {
        setDeployment(deployment);
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
        if (getDeployment() != null)
            sb.append("Deployment: ").append(getDeployment());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof DescribeEnvironmentDeploymentResult == false)
            return false;
        DescribeEnvironmentDeploymentResult other = (DescribeEnvironmentDeploymentResult) obj;
        if (other.getDeployment() == null ^ this.getDeployment() == null)
            return false;
        if (other.getDeployment() != null && other.getDeployment().equals(this.getDeployment()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getDeployment() == null) ? 0 : getDeployment().hashCode());
        return hashCode;
    }

    @Override
    public DescribeEnvironmentDeploymentResult clone() {
        try {
            return (DescribeEnvironmentDeploymentResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
