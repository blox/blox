/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/ListEnvironmentRevisions"
 *      target="_top">AWS API Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentRevisionsResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private java.util.List<String> revisionIds;

    /**
     * @return
     */

    public java.util.List<String> getRevisionIds() {
        return revisionIds;
    }

    /**
     * @param revisionIds
     */

    public void setRevisionIds(java.util.Collection<String> revisionIds) {
        if (revisionIds == null) {
            this.revisionIds = null;
            return;
        }

        this.revisionIds = new java.util.ArrayList<String>(revisionIds);
    }

    /**
     * <p>
     * <b>NOTE:</b> This method appends the values to the existing list (if any). Use
     * {@link #setRevisionIds(java.util.Collection)} or {@link #withRevisionIds(java.util.Collection)} if you want to
     * override the existing values.
     * </p>
     * 
     * @param revisionIds
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentRevisionsResult revisionIds(String... revisionIds) {
        if (this.revisionIds == null) {
            setRevisionIds(new java.util.ArrayList<String>(revisionIds.length));
        }
        for (String ele : revisionIds) {
            this.revisionIds.add(ele);
        }
        return this;
    }

    /**
     * @param revisionIds
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentRevisionsResult revisionIds(java.util.Collection<String> revisionIds) {
        setRevisionIds(revisionIds);
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
        if (getRevisionIds() != null)
            sb.append("RevisionIds: ").append(getRevisionIds());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof ListEnvironmentRevisionsResult == false)
            return false;
        ListEnvironmentRevisionsResult other = (ListEnvironmentRevisionsResult) obj;
        if (other.getRevisionIds() == null ^ this.getRevisionIds() == null)
            return false;
        if (other.getRevisionIds() != null && other.getRevisionIds().equals(this.getRevisionIds()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getRevisionIds() == null) ? 0 : getRevisionIds().hashCode());
        return hashCode;
    }

    @Override
    public ListEnvironmentRevisionsResult clone() {
        try {
            return (ListEnvironmentRevisionsResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
