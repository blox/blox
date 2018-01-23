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

    private String activeEnvironmentRevisionId;

    private String cluster;

    private DeploymentConfiguration deploymentConfiguration;

    private String deploymentMethod;

    private String environmentHealth;

    private String environmentName;

    private String environmentType;

    private String latestEnvironmentRevisionId;

    private String role;

    /**
     * @param activeEnvironmentRevisionId
     */

    public void setActiveEnvironmentRevisionId(String activeEnvironmentRevisionId) {
        this.activeEnvironmentRevisionId = activeEnvironmentRevisionId;
    }

    /**
     * @return
     */

    public String getActiveEnvironmentRevisionId() {
        return this.activeEnvironmentRevisionId;
    }

    /**
     * @param activeEnvironmentRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment activeEnvironmentRevisionId(String activeEnvironmentRevisionId) {
        setActiveEnvironmentRevisionId(activeEnvironmentRevisionId);
        return this;
    }

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
     * @param deploymentConfiguration
     */

    public void setDeploymentConfiguration(DeploymentConfiguration deploymentConfiguration) {
        this.deploymentConfiguration = deploymentConfiguration;
    }

    /**
     * @return
     */

    public DeploymentConfiguration getDeploymentConfiguration() {
        return this.deploymentConfiguration;
    }

    /**
     * @param deploymentConfiguration
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment deploymentConfiguration(DeploymentConfiguration deploymentConfiguration) {
        setDeploymentConfiguration(deploymentConfiguration);
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

    public Environment deploymentMethod(String deploymentMethod) {
        setDeploymentMethod(deploymentMethod);
        return this;
    }

    /**
     * @param environmentHealth
     */

    public void setEnvironmentHealth(String environmentHealth) {
        this.environmentHealth = environmentHealth;
    }

    /**
     * @return
     */

    public String getEnvironmentHealth() {
        return this.environmentHealth;
    }

    /**
     * @param environmentHealth
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment environmentHealth(String environmentHealth) {
        setEnvironmentHealth(environmentHealth);
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

    public Environment environmentType(String environmentType) {
        setEnvironmentType(environmentType);
        return this;
    }

    /**
     * @param latestEnvironmentRevisionId
     */

    public void setLatestEnvironmentRevisionId(String latestEnvironmentRevisionId) {
        this.latestEnvironmentRevisionId = latestEnvironmentRevisionId;
    }

    /**
     * @return
     */

    public String getLatestEnvironmentRevisionId() {
        return this.latestEnvironmentRevisionId;
    }

    /**
     * @param latestEnvironmentRevisionId
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment latestEnvironmentRevisionId(String latestEnvironmentRevisionId) {
        setLatestEnvironmentRevisionId(latestEnvironmentRevisionId);
        return this;
    }

    /**
     * @param role
     */

    public void setRole(String role) {
        this.role = role;
    }

    /**
     * @return
     */

    public String getRole() {
        return this.role;
    }

    /**
     * @param role
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public Environment role(String role) {
        setRole(role);
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
        if (getActiveEnvironmentRevisionId() != null)
            sb.append("ActiveEnvironmentRevisionId: ").append(getActiveEnvironmentRevisionId()).append(",");
        if (getCluster() != null)
            sb.append("Cluster: ").append(getCluster()).append(",");
        if (getDeploymentConfiguration() != null)
            sb.append("DeploymentConfiguration: ").append(getDeploymentConfiguration()).append(",");
        if (getDeploymentMethod() != null)
            sb.append("DeploymentMethod: ").append(getDeploymentMethod()).append(",");
        if (getEnvironmentHealth() != null)
            sb.append("EnvironmentHealth: ").append(getEnvironmentHealth()).append(",");
        if (getEnvironmentName() != null)
            sb.append("EnvironmentName: ").append(getEnvironmentName()).append(",");
        if (getEnvironmentType() != null)
            sb.append("EnvironmentType: ").append(getEnvironmentType()).append(",");
        if (getLatestEnvironmentRevisionId() != null)
            sb.append("LatestEnvironmentRevisionId: ").append(getLatestEnvironmentRevisionId()).append(",");
        if (getRole() != null)
            sb.append("Role: ").append(getRole());
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
        if (other.getActiveEnvironmentRevisionId() == null ^ this.getActiveEnvironmentRevisionId() == null)
            return false;
        if (other.getActiveEnvironmentRevisionId() != null && other.getActiveEnvironmentRevisionId().equals(this.getActiveEnvironmentRevisionId()) == false)
            return false;
        if (other.getCluster() == null ^ this.getCluster() == null)
            return false;
        if (other.getCluster() != null && other.getCluster().equals(this.getCluster()) == false)
            return false;
        if (other.getDeploymentConfiguration() == null ^ this.getDeploymentConfiguration() == null)
            return false;
        if (other.getDeploymentConfiguration() != null && other.getDeploymentConfiguration().equals(this.getDeploymentConfiguration()) == false)
            return false;
        if (other.getDeploymentMethod() == null ^ this.getDeploymentMethod() == null)
            return false;
        if (other.getDeploymentMethod() != null && other.getDeploymentMethod().equals(this.getDeploymentMethod()) == false)
            return false;
        if (other.getEnvironmentHealth() == null ^ this.getEnvironmentHealth() == null)
            return false;
        if (other.getEnvironmentHealth() != null && other.getEnvironmentHealth().equals(this.getEnvironmentHealth()) == false)
            return false;
        if (other.getEnvironmentName() == null ^ this.getEnvironmentName() == null)
            return false;
        if (other.getEnvironmentName() != null && other.getEnvironmentName().equals(this.getEnvironmentName()) == false)
            return false;
        if (other.getEnvironmentType() == null ^ this.getEnvironmentType() == null)
            return false;
        if (other.getEnvironmentType() != null && other.getEnvironmentType().equals(this.getEnvironmentType()) == false)
            return false;
        if (other.getLatestEnvironmentRevisionId() == null ^ this.getLatestEnvironmentRevisionId() == null)
            return false;
        if (other.getLatestEnvironmentRevisionId() != null && other.getLatestEnvironmentRevisionId().equals(this.getLatestEnvironmentRevisionId()) == false)
            return false;
        if (other.getRole() == null ^ this.getRole() == null)
            return false;
        if (other.getRole() != null && other.getRole().equals(this.getRole()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getActiveEnvironmentRevisionId() == null) ? 0 : getActiveEnvironmentRevisionId().hashCode());
        hashCode = prime * hashCode + ((getCluster() == null) ? 0 : getCluster().hashCode());
        hashCode = prime * hashCode + ((getDeploymentConfiguration() == null) ? 0 : getDeploymentConfiguration().hashCode());
        hashCode = prime * hashCode + ((getDeploymentMethod() == null) ? 0 : getDeploymentMethod().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentHealth() == null) ? 0 : getEnvironmentHealth().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentName() == null) ? 0 : getEnvironmentName().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentType() == null) ? 0 : getEnvironmentType().hashCode());
        hashCode = prime * hashCode + ((getLatestEnvironmentRevisionId() == null) ? 0 : getLatestEnvironmentRevisionId().hashCode());
        hashCode = prime * hashCode + ((getRole() == null) ? 0 : getRole().hashCode());
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
