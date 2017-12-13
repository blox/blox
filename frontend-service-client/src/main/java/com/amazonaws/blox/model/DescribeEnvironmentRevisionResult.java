/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/DescribeEnvironmentRevision"
 *      target="_top">AWS API Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class DescribeEnvironmentRevisionResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private EnvironmentRevision revision;

    /**
     * @param revision
     */

    public void setRevision(EnvironmentRevision revision) {
        this.revision = revision;
    }

    /**
     * @return
     */

    public EnvironmentRevision getRevision() {
        return this.revision;
    }

    /**
     * @param revision
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public DescribeEnvironmentRevisionResult revision(EnvironmentRevision revision) {
        setRevision(revision);
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
        if (getRevision() != null)
            sb.append("Revision: ").append(getRevision());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof DescribeEnvironmentRevisionResult == false)
            return false;
        DescribeEnvironmentRevisionResult other = (DescribeEnvironmentRevisionResult) obj;
        if (other.getRevision() == null ^ this.getRevision() == null)
            return false;
        if (other.getRevision() != null && other.getRevision().equals(this.getRevision()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getRevision() == null) ? 0 : getRevision().hashCode());
        return hashCode;
    }

    @Override
    public DescribeEnvironmentRevisionResult clone() {
        try {
            return (DescribeEnvironmentRevisionResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
