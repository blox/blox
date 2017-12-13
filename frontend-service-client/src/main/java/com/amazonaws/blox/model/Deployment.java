/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;
import com.amazonaws.protocol.StructuredPojo;
import com.amazonaws.protocol.ProtocolMarshaller;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/Deployment" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class Deployment implements Serializable, Cloneable, StructuredPojo {

    private String deploymentId;

    private String environmentName;

    private String newTargetRevisionId;

    private String oldTargetRevisionId;

    private Integer timestamp;

    /**
     * @param deploymentId
     */

    public void setDeploymentId(String deploymentId) {
        this.deploymentId = deploymentId;
    }

    /**
     * @return
     */

    public String getDeploymentId() {
        return this.deploymentId;
    }

    /**
     * @param deploymentId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Deployment deploymentId(String deploymentId) {
        setDeploymentId(deploymentId);
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

    public Deployment environmentName(String environmentName) {
        setEnvironmentName(environmentName);
        return this;
    }

    /**
     * @param newTargetRevisionId
     */

    public void setNewTargetRevisionId(String newTargetRevisionId) {
        this.newTargetRevisionId = newTargetRevisionId;
    }

    /**
     * @return
     */

    public String getNewTargetRevisionId() {
        return this.newTargetRevisionId;
    }

    /**
     * @param newTargetRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Deployment newTargetRevisionId(String newTargetRevisionId) {
        setNewTargetRevisionId(newTargetRevisionId);
        return this;
    }

    /**
     * @param oldTargetRevisionId
     */

    public void setOldTargetRevisionId(String oldTargetRevisionId) {
        this.oldTargetRevisionId = oldTargetRevisionId;
    }

    /**
     * @return
     */

    public String getOldTargetRevisionId() {
        return this.oldTargetRevisionId;
    }

    /**
     * @param oldTargetRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Deployment oldTargetRevisionId(String oldTargetRevisionId) {
        setOldTargetRevisionId(oldTargetRevisionId);
        return this;
    }

    /**
     * @param timestamp
     */

    public void setTimestamp(Integer timestamp) {
        this.timestamp = timestamp;
    }

    /**
     * @return
     */

    public Integer getTimestamp() {
        return this.timestamp;
    }

    /**
     * @param timestamp
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Deployment timestamp(Integer timestamp) {
        setTimestamp(timestamp);
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
        if (getDeploymentId() != null)
            sb.append("DeploymentId: ").append(getDeploymentId()).append(",");
        if (getEnvironmentName() != null)
            sb.append("EnvironmentName: ").append(getEnvironmentName()).append(",");
        if (getNewTargetRevisionId() != null)
            sb.append("NewTargetRevisionId: ").append(getNewTargetRevisionId()).append(",");
        if (getOldTargetRevisionId() != null)
            sb.append("OldTargetRevisionId: ").append(getOldTargetRevisionId()).append(",");
        if (getTimestamp() != null)
            sb.append("Timestamp: ").append(getTimestamp());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof Deployment == false)
            return false;
        Deployment other = (Deployment) obj;
        if (other.getDeploymentId() == null ^ this.getDeploymentId() == null)
            return false;
        if (other.getDeploymentId() != null && other.getDeploymentId().equals(this.getDeploymentId()) == false)
            return false;
        if (other.getEnvironmentName() == null ^ this.getEnvironmentName() == null)
            return false;
        if (other.getEnvironmentName() != null && other.getEnvironmentName().equals(this.getEnvironmentName()) == false)
            return false;
        if (other.getNewTargetRevisionId() == null ^ this.getNewTargetRevisionId() == null)
            return false;
        if (other.getNewTargetRevisionId() != null && other.getNewTargetRevisionId().equals(this.getNewTargetRevisionId()) == false)
            return false;
        if (other.getOldTargetRevisionId() == null ^ this.getOldTargetRevisionId() == null)
            return false;
        if (other.getOldTargetRevisionId() != null && other.getOldTargetRevisionId().equals(this.getOldTargetRevisionId()) == false)
            return false;
        if (other.getTimestamp() == null ^ this.getTimestamp() == null)
            return false;
        if (other.getTimestamp() != null && other.getTimestamp().equals(this.getTimestamp()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getDeploymentId() == null) ? 0 : getDeploymentId().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentName() == null) ? 0 : getEnvironmentName().hashCode());
        hashCode = prime * hashCode + ((getNewTargetRevisionId() == null) ? 0 : getNewTargetRevisionId().hashCode());
        hashCode = prime * hashCode + ((getOldTargetRevisionId() == null) ? 0 : getOldTargetRevisionId().hashCode());
        hashCode = prime * hashCode + ((getTimestamp() == null) ? 0 : getTimestamp().hashCode());
        return hashCode;
    }

    @Override
    public Deployment clone() {
        try {
            return (Deployment) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

    @com.amazonaws.annotation.SdkInternalApi
    @Override
    public void marshall(ProtocolMarshaller protocolMarshaller) {
        com.amazonaws.blox.model.transform.DeploymentMarshaller.getInstance().marshall(this, protocolMarshaller);
    }
}
