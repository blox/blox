/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/UpdateEnvironment" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class UpdateEnvironmentResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private String environmentRevisionId;

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

    public UpdateEnvironmentResult environmentRevisionId(String environmentRevisionId) {
        setEnvironmentRevisionId(environmentRevisionId);
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
        if (getEnvironmentRevisionId() != null)
            sb.append("EnvironmentRevisionId: ").append(getEnvironmentRevisionId());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof UpdateEnvironmentResult == false)
            return false;
        UpdateEnvironmentResult other = (UpdateEnvironmentResult) obj;
        if (other.getEnvironmentRevisionId() == null ^ this.getEnvironmentRevisionId() == null)
            return false;
        if (other.getEnvironmentRevisionId() != null && other.getEnvironmentRevisionId().equals(this.getEnvironmentRevisionId()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getEnvironmentRevisionId() == null) ? 0 : getEnvironmentRevisionId().hashCode());
        return hashCode;
    }

    @Override
    public UpdateEnvironmentResult clone() {
        try {
            return (UpdateEnvironmentResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
