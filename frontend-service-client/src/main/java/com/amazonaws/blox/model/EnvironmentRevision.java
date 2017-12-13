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

    private String cluster;

    private TaskCounts counts;

    private java.util.Map<String, String> deploymentConfiguration;

    private String deploymentMethod;

    private String environmentName;

    private String environmentRevisionId;

    private String environmentType;

    private InstanceGroup instanceGroup;

    private String taskDefinition;

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

    public EnvironmentRevision cluster(String cluster) {
        setCluster(cluster);
        return this;
    }

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
     * @return
     */

    public java.util.Map<String, String> getDeploymentConfiguration() {
        return deploymentConfiguration;
    }

    /**
     * @param deploymentConfiguration
     */

    public void setDeploymentConfiguration(java.util.Map<String, String> deploymentConfiguration) {
        this.deploymentConfiguration = deploymentConfiguration;
    }

    /**
     * @param deploymentConfiguration
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision deploymentConfiguration(java.util.Map<String, String> deploymentConfiguration) {
        setDeploymentConfiguration(deploymentConfiguration);
        return this;
    }

    public EnvironmentRevision addDeploymentConfigurationEntry(String key, String value) {
        if (null == this.deploymentConfiguration) {
            this.deploymentConfiguration = new java.util.HashMap<String, String>();
        }
        if (this.deploymentConfiguration.containsKey(key))
            throw new IllegalArgumentException("Duplicated keys (" + key.toString() + ") are provided.");
        this.deploymentConfiguration.put(key, value);
        return this;
    }

    /**
     * Removes all the entries added into DeploymentConfiguration.
     *
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision clearDeploymentConfigurationEntries() {
        this.deploymentConfiguration = null;
        return this;
    }

    /**
     * @param deploymentMethod
     */

    public void setDeploymentMethod(String deploymentMethod) {
        this.deploymentMethod = deploymentMethod;
    }

    /**
     * @return
     */

    public String getDeploymentMethod() {
        return this.deploymentMethod;
    }

    /**
     * @param deploymentMethod
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision deploymentMethod(String deploymentMethod) {
        setDeploymentMethod(deploymentMethod);
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

    public EnvironmentRevision environmentName(String environmentName) {
        setEnvironmentName(environmentName);
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
     * @param environmentType
     */

    public void setEnvironmentType(String environmentType) {
        this.environmentType = environmentType;
    }

    /**
     * @return
     */

    public String getEnvironmentType() {
        return this.environmentType;
    }

    /**
     * @param environmentType
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public EnvironmentRevision environmentType(String environmentType) {
        setEnvironmentType(environmentType);
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
        if (getCluster() != null)
            sb.append("Cluster: ").append(getCluster()).append(",");
        if (getCounts() != null)
            sb.append("Counts: ").append(getCounts()).append(",");
        if (getDeploymentConfiguration() != null)
            sb.append("DeploymentConfiguration: ").append(getDeploymentConfiguration()).append(",");
        if (getDeploymentMethod() != null)
            sb.append("DeploymentMethod: ").append(getDeploymentMethod()).append(",");
        if (getEnvironmentName() != null)
            sb.append("EnvironmentName: ").append(getEnvironmentName()).append(",");
        if (getEnvironmentRevisionId() != null)
            sb.append("EnvironmentRevisionId: ").append(getEnvironmentRevisionId()).append(",");
        if (getEnvironmentType() != null)
            sb.append("EnvironmentType: ").append(getEnvironmentType()).append(",");
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
        if (other.getCluster() == null ^ this.getCluster() == null)
            return false;
        if (other.getCluster() != null && other.getCluster().equals(this.getCluster()) == false)
            return false;
        if (other.getCounts() == null ^ this.getCounts() == null)
            return false;
        if (other.getCounts() != null && other.getCounts().equals(this.getCounts()) == false)
            return false;
        if (other.getDeploymentConfiguration() == null ^ this.getDeploymentConfiguration() == null)
            return false;
        if (other.getDeploymentConfiguration() != null && other.getDeploymentConfiguration().equals(this.getDeploymentConfiguration()) == false)
            return false;
        if (other.getDeploymentMethod() == null ^ this.getDeploymentMethod() == null)
            return false;
        if (other.getDeploymentMethod() != null && other.getDeploymentMethod().equals(this.getDeploymentMethod()) == false)
            return false;
        if (other.getEnvironmentName() == null ^ this.getEnvironmentName() == null)
            return false;
        if (other.getEnvironmentName() != null && other.getEnvironmentName().equals(this.getEnvironmentName()) == false)
            return false;
        if (other.getEnvironmentRevisionId() == null ^ this.getEnvironmentRevisionId() == null)
            return false;
        if (other.getEnvironmentRevisionId() != null && other.getEnvironmentRevisionId().equals(this.getEnvironmentRevisionId()) == false)
            return false;
        if (other.getEnvironmentType() == null ^ this.getEnvironmentType() == null)
            return false;
        if (other.getEnvironmentType() != null && other.getEnvironmentType().equals(this.getEnvironmentType()) == false)
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

        hashCode = prime * hashCode + ((getCluster() == null) ? 0 : getCluster().hashCode());
        hashCode = prime * hashCode + ((getCounts() == null) ? 0 : getCounts().hashCode());
        hashCode = prime * hashCode + ((getDeploymentConfiguration() == null) ? 0 : getDeploymentConfiguration().hashCode());
        hashCode = prime * hashCode + ((getDeploymentMethod() == null) ? 0 : getDeploymentMethod().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentName() == null) ? 0 : getEnvironmentName().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentRevisionId() == null) ? 0 : getEnvironmentRevisionId().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentType() == null) ? 0 : getEnvironmentType().hashCode());
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
