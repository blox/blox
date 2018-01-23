/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;
import com.amazonaws.protocol.StructuredPojo;
import com.amazonaws.protocol.ProtocolMarshaller;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/EnvironmentRevision" target="_top">AWS
 *      API Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class EnvironmentRevision implements Serializable, Cloneable, StructuredPojo {

    private TaskCounts counts;

    private String environmentRevisionId;

    private InstanceGroup instanceGroup;

    private String taskDefinition;

    /**
     * @param counts
     */

    public void setCounts(TaskCounts counts) {
        this.counts = counts;
    }

    /**
     * @return
     */

    public TaskCounts getCounts() {
        return this.counts;
    }

    /**
     * @param counts
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision counts(TaskCounts counts) {
        setCounts(counts);
        return this;
    }

    /**
     * @param environmentRevisionId
     */

    public void setEnvironmentRevisionId(String environmentRevisionId) {
        this.environmentRevisionId = environmentRevisionId;
    }

    /**
     * @return
     */

    public String getEnvironmentRevisionId() {
        return this.environmentRevisionId;
    }

    /**
     * @param environmentRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision environmentRevisionId(String environmentRevisionId) {
        setEnvironmentRevisionId(environmentRevisionId);
        return this;
    }

    /**
     * @param instanceGroup
     */

    public void setInstanceGroup(InstanceGroup instanceGroup) {
        this.instanceGroup = instanceGroup;
    }

    /**
     * @return
     */

    public InstanceGroup getInstanceGroup() {
        return this.instanceGroup;
    }

    /**
     * @param instanceGroup
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision instanceGroup(InstanceGroup instanceGroup) {
        setInstanceGroup(instanceGroup);
        return this;
    }

    /**
     * @param taskDefinition
     */

    public void setTaskDefinition(String taskDefinition) {
        this.taskDefinition = taskDefinition;
    }

    /**
     * @return
     */

    public String getTaskDefinition() {
        return this.taskDefinition;
    }

    /**
     * @param taskDefinition
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision taskDefinition(String taskDefinition) {
        setTaskDefinition(taskDefinition);
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
        if (getCounts() != null)
            sb.append("Counts: ").append(getCounts()).append(",");
        if (getEnvironmentRevisionId() != null)
            sb.append("EnvironmentRevisionId: ").append(getEnvironmentRevisionId()).append(",");
        if (getInstanceGroup() != null)
            sb.append("InstanceGroup: ").append(getInstanceGroup()).append(",");
        if (getTaskDefinition() != null)
            sb.append("TaskDefinition: ").append(getTaskDefinition());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof EnvironmentRevision == false)
            return false;
        EnvironmentRevision other = (EnvironmentRevision) obj;
        if (other.getCounts() == null ^ this.getCounts() == null)
            return false;
        if (other.getCounts() != null && other.getCounts().equals(this.getCounts()) == false)
            return false;
        if (other.getEnvironmentRevisionId() == null ^ this.getEnvironmentRevisionId() == null)
            return false;
        if (other.getEnvironmentRevisionId() != null && other.getEnvironmentRevisionId().equals(this.getEnvironmentRevisionId()) == false)
            return false;
        if (other.getInstanceGroup() == null ^ this.getInstanceGroup() == null)
            return false;
        if (other.getInstanceGroup() != null && other.getInstanceGroup().equals(this.getInstanceGroup()) == false)
            return false;
        if (other.getTaskDefinition() == null ^ this.getTaskDefinition() == null)
            return false;
        if (other.getTaskDefinition() != null && other.getTaskDefinition().equals(this.getTaskDefinition()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getCounts() == null) ? 0 : getCounts().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentRevisionId() == null) ? 0 : getEnvironmentRevisionId().hashCode());
        hashCode = prime * hashCode + ((getInstanceGroup() == null) ? 0 : getInstanceGroup().hashCode());
        hashCode = prime * hashCode + ((getTaskDefinition() == null) ? 0 : getTaskDefinition().hashCode());
        return hashCode;
    }

    @Override
    public EnvironmentRevision clone() {
        try {
            return (EnvironmentRevision) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

    @com.amazonaws.annotation.SdkInternalApi
    @Override
    public void marshall(ProtocolMarshaller protocolMarshaller) {
        com.amazonaws.blox.model.transform.EnvironmentRevisionMarshaller.getInstance().marshall(this, protocolMarshaller);
    }
}
