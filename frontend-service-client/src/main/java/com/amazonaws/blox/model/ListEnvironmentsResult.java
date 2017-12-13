/**

*/
package com.amazonaws.blox.model;

import java.io.Serializable;
import javax.annotation.Generated;

/**
 * 
 * @see <a href="http://docs.aws.amazon.com/goto/WebAPI/ecs-blox-v2017-07-11/ListEnvironments" target="_top">AWS API
 *      Documentation</a>
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class ListEnvironmentsResult extends com.amazonaws.opensdk.BaseResult implements Serializable, Cloneable {

    private java.util.List<String> environmentNames;

    private String nextToken;

    /**
     * @return
     */

    public java.util.List<String> getEnvironmentNames() {
        return environmentNames;
    }

    /**
     * @param environmentNames
     */

    public void setEnvironmentNames(java.util.Collection<String> environmentNames) {
        if (environmentNames == null) {
            this.environmentNames = null;
            return;
        }

        this.environmentNames = new java.util.ArrayList<String>(environmentNames);
    }

    /**
     * <p>
     * <b>NOTE:</b> This method appends the values to the existing list (if any). Use
     * {@link #setEnvironmentNames(java.util.Collection)} or {@link #withEnvironmentNames(java.util.Collection)} if you
     * want to override the existing values.
     * </p>
     * 
     * @param environmentNames
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsResult environmentNames(String... environmentNames) {
        if (this.environmentNames == null) {
            setEnvironmentNames(new java.util.ArrayList<String>(environmentNames.length));
        }
        for (String ele : environmentNames) {
            this.environmentNames.add(ele);
        }
        return this;
    }

    /**
     * @param environmentNames
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsResult environmentNames(java.util.Collection<String> environmentNames) {
        setEnvironmentNames(environmentNames);
        return this;
    }

    /**
     * @param nextToken
     */

    public void setNextToken(String nextToken) {
        this.nextToken = nextToken;
    }

    /**
     * @return
     */

    public String getNextToken() {
        return this.nextToken;
    }

    /**
     * @param nextToken
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public ListEnvironmentsResult nextToken(String nextToken) {
        setNextToken(nextToken);
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
        if (getEnvironmentNames() != null)
            sb.append("EnvironmentNames: ").append(getEnvironmentNames()).append(",");
        if (getNextToken() != null)
            sb.append("NextToken: ").append(getNextToken());
        sb.append("}");
        return sb.toString();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;

        if (obj instanceof ListEnvironmentsResult == false)
            return false;
        ListEnvironmentsResult other = (ListEnvironmentsResult) obj;
        if (other.getEnvironmentNames() == null ^ this.getEnvironmentNames() == null)
            return false;
        if (other.getEnvironmentNames() != null && other.getEnvironmentNames().equals(this.getEnvironmentNames()) == false)
            return false;
        if (other.getNextToken() == null ^ this.getNextToken() == null)
            return false;
        if (other.getNextToken() != null && other.getNextToken().equals(this.getNextToken()) == false)
            return false;
        return true;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int hashCode = 1;

        hashCode = prime * hashCode + ((getEnvironmentNames() == null) ? 0 : getEnvironmentNames().hashCode());
        hashCode = prime * hashCode + ((getNextToken() == null) ? 0 : getNextToken().hashCode());
        return hashCode;
    }

    @Override
    public ListEnvironmentsResult clone() {
        try {
            return (ListEnvironmentsResult) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new IllegalStateException("Got a CloneNotSupportedException from Object.clone() " + "even though we're Cloneable!", e);
        }
    }

}
