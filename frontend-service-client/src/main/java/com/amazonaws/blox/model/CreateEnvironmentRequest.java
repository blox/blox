/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

import com.amazonaws.auth.RequestSigner;
import com.amazonaws.opensdk.protect.auth.RequestSignerAware;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/CreateEnvironment" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class CreateEnvironmentRequest extends com.amazonaws.opensdk.BaseRequest implements Serializable, Cloneable, RequestSignerAware {

    private String cluster;

    private DeploymentConfiguration deploymentConfiguration;

    private String deploymentMethod;

    private String environmentName;

    private String environmentType;

    private InstanceGroup instanceGroup;

    private String role;

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

    public CreateEnvironmentRequest cluster(String cluster) {
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

    public CreateEnvironmentRequest deploymentConfiguration(DeploymentConfiguration deploymentConfiguration) {
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

    public CreateEnvironmentRequest deploymentMethod(String deploymentMethod) {
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

    public CreateEnvironmentRequest environmentName(String environmentName) {
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

    public CreateEnvironmentRequest environmentType(String environmentType) {
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

    public CreateEnvironmentRequest instanceGroup(InstanceGroup instanceGroup) {
        setInstanceGroup(instanceGroup);
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

    public CreateEnvironmentRequest role(String role) {
        setRole(role);
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

    public CreateEnvironmentRequest taskDefinition(String taskDefinition) {
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
        if (getDeploymentConfiguration() != null)
            sb.append("DeploymentConfiguration: ").append(getDeploymentConfiguration()).append(",");
        if (getDeploymentMethod() != null)
            sb.append("DeploymentMethod: ").append(getDeploymentMethod()).append(",");
        if (getEnvironmentName() != null)
            sb.append("EnvironmentName: ").append(getEnvironmentName()).append(",");
        if (getEnvironmentType() != null)
            sb.append("EnvironmentType: ").append(getEnvironmentType()).append(",");
        if (getInstanceGroup() != null)
            sb.append("InstanceGroup: ").append(getInstanceGroup()).append(",");
        if (getRole() != null)
            sb.append("Role: ").append(getRole()).append(",");
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

        if (obj instanceof CreateEnvironmentRequest == false)
            return false;
        CreateEnvironmentRequest other = (CreateEnvironmentRequest) obj;
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
        if (other.getEnvironmentName() == null ^ this.getEnvironmentName() == null)
            return false;
        if (other.getEnvironmentName() != null && other.getEnvironmentName().equals(this.getEnvironmentName()) == false)
            return false;
        if (other.getEnvironmentType() == null ^ this.getEnvironmentType() == null)
            return false;
        if (other.getEnvironmentType() != null && other.getEnvironmentType().equals(this.getEnvironmentType()) == false)
            return false;
        if (other.getInstanceGroup() == null ^ this.getInstanceGroup() == null)
            return false;
        if (other.getInstanceGroup() != null && other.getInstanceGroup().equals(this.getInstanceGroup()) == false)
            return false;
        if (other.getRole() == null ^ this.getRole() == null)
            return false;
        if (other.getRole() != null && other.getRole().equals(this.getRole()) == false)
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
        hashCode = prime * hashCode + ((getDeploymentConfiguration() == null) ? 0 : getDeploymentConfiguration().hashCode());
        hashCode = prime * hashCode + ((getDeploymentMethod() == null) ? 0 : getDeploymentMethod().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentName() == null) ? 0 : getEnvironmentName().hashCode());
        hashCode = prime * hashCode + ((getEnvironmentType() == null) ? 0 : getEnvironmentType().hashCode());
        hashCode = prime * hashCode + ((getInstanceGroup() == null) ? 0 : getInstanceGroup().hashCode());
        hashCode = prime * hashCode + ((getRole() == null) ? 0 : getRole().hashCode());
        hashCode = prime * hashCode + ((getTaskDefinition() == null) ? 0 : getTaskDefinition().hashCode());
        return hashCode;
    }

    @Override
    public CreateEnvironmentRequest clone() {
        return (CreateEnvironmentRequest) super.clone();
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
    public CreateEnvironmentRequest sdkRequestConfig(com.amazonaws.opensdk.SdkRequestConfig sdkRequestConfig) {
        super.sdkRequestConfig(sdkRequestConfig);
        return this;
    }

}
