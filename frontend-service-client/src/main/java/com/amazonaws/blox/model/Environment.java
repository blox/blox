/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;
import com.amazonaws.protocol.StructuredPojo;
import com.amazonaws.protocol.ProtocolMarshaller;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/Environment" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class Environment implements Serializable, Cloneable, StructuredPojo {

    private String cluster;

    private String environmentName;

    private String targetRevisionId;

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

    public Environment cluster(String cluster) {
        setCluster(cluster);
        return this;
    }

    /**
     * @param environmentName
     */

    public void setEnvironmentName(String environmentName) {
        this.environmentName = environmentName;
    }

    /**
     * @return
     */

    public String getEnvironmentName() {
        return this.environmentName;
    }

    /**
     * @param environmentName
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment environmentName(String environmentName) {
        setEnvironmentName(environmentName);
        return this;
    }

    /**
     * @param targetRevisionId
     */

    public void setTargetRevisionId(String targetRevisionId) {
        this.targetRevisionId = targetRevisionId;
    }

    /**
     * @return
     */

    public String getTargetRevisionId() {
        return this.targetRevisionId;
    }

    /**
     * @param targetRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment targetRevisionId(String targetRevisionId) {
        setTargetRevisionId(targetRevisionId);
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
        if (getEnvironmentName() != null)
            sb.append("EnvironmentName: ").append(getEnvironmentName()).append(",");
        if (getTargetRevisionId() != null)
            sb.append("TargetRevisionId: ").append(getTargetRevisionId());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof Environment == false)
            return false;
        Environment other = (Environment) obj;
        if (other.getCluster() == null ^ this.getCluster() == null)
            return false;
        if (other.getCluster() != null && other.getCluster().equals(this.getCluster()) == false)
            return false;
        if (other.getEnvironmentName() == null ^ this.getEnvironmentName() == null)
            return false;
        if (other.getEnvironmentName() != null && other.getEnvironmentName().equals(this.getEnvironmentName()) == false)
            return false;
        if (other.getTargetRevisionId() == null ^ this.getTargetRevisionId() == null)
            return false;
        if (other.getTargetRevisionId() != null && other.getTargetRevisionId().equals(this.getTargetRevisionId()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getCluster() == null) ? 0 : getCluster().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentName() == null) ? 0 : getEnvironmentName().hashCode());
        hashCode = prime * hashCode + ((getTargetRevisionId() == null) ? 0 : getTargetRevisionId().hashCode());
        return hashCode;
    }

    @Override
    public Environment clone() {
        try {
            return (Environment) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

    @com.amazonaws.annotation.SdkInternalApi
    @Override
    public void marshall(ProtocolMarshaller protocolMarshaller) {
        com.amazonaws.blox.model.transform.EnvironmentMarshaller.getInstance().marshall(this, protocolMarshaller);
    }
}
