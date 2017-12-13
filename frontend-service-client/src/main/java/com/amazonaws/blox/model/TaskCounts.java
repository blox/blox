/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;
import com.amazonaws.protocol.StructuredPojo;
import com.amazonaws.protocol.ProtocolMarshaller;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/TaskCounts" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class TaskCounts implements Serializable, Cloneable, StructuredPojo {

    private Integer desired;

    private Integer healthy;

    private Integer total;

    /**
     * @param desired
     */

    public void setDesired(Integer desired) {
        this.desired = desired;
    }

    /**
     * @return
     */

    public Integer getDesired() {
        return this.desired;
    }

    /**
     * @param desired
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public TaskCounts desired(Integer desired) {
        setDesired(desired);
        return this;
    }

    /**
     * @param healthy
     */

    public void setHealthy(Integer healthy) {
        this.healthy = healthy;
    }

    /**
     * @return
     */

    public Integer getHealthy() {
        return this.healthy;
    }

    /**
     * @param healthy
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public TaskCounts healthy(Integer healthy) {
        setHealthy(healthy);
        return this;
    }

    /**
     * @param total
     */

    public void setTotal(Integer total) {
        this.total = total;
    }

    /**
     * @return
     */

    public Integer getTotal() {
        return this.total;
    }

    /**
     * @param total
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public TaskCounts total(Integer total) {
        setTotal(total);
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
        if (getDesired() != null)
            sb.append("Desired: ").append(getDesired()).append(",");
        if (getHealthy() != null)
            sb.append("Healthy: ").append(getHealthy()).append(",");
        if (getTotal() != null)
            sb.append("Total: ").append(getTotal());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof TaskCounts == false)
            return false;
        TaskCounts other = (TaskCounts) obj;
        if (other.getDesired() == null ^ this.getDesired() == null)
            return false;
        if (other.getDesired() != null && other.getDesired().equals(this.getDesired()) == false)
            return false;
        if (other.getHealthy() == null ^ this.getHealthy() == null)
            return false;
        if (other.getHealthy() != null && other.getHealthy().equals(this.getHealthy()) == false)
            return false;
        if (other.getTotal() == null ^ this.getTotal() == null)
            return false;
        if (other.getTotal() != null && other.getTotal().equals(this.getTotal()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getDesired() == null) ? 0 : getDesired().hashCode());
        hashCode = prime * hashCode + ((getHealthy() == null) ? 0 : getHealthy().hashCode());
        hashCode = prime * hashCode + ((getTotal() == null) ? 0 : getTotal().hashCode());
        return hashCode;
    }

    @Override
    public TaskCounts clone() {
        try {
            return (TaskCounts) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

    @com.amazonaws.annotation.SdkInternalApi
    @Override
    public void marshall(ProtocolMarshaller protocolMarshaller) {
        com.amazonaws.blox.model.transform.TaskCountsMarshaller.getInstance().marshall(this, protocolMarshaller);
    }
}
