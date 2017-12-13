/**

*/
package com.amazonaws.blox.model;

import javax.annotation.Generated;

/**
 * 
 */
@Generated("com.amazonaws:aws-java-sdk-code-generator")
public class NotFoundException extends com.amazonaws.blox.model.BloxException {
    private static final long serialVersionUID = 1L;

    private Integer code;

    /**
     * Constructs a new NotFoundException with the specified error message.
     *
     * @param message
     *        Describes the error encountered.
     */
    public NotFoundException(String message) {
        super(message);
    }

    /**
     * @param code
     */

    @com.fasterxml.jackson.annotation.JsonProperty("code")
    public void setCode(Integer code) {
        this.code = code;
    }

    /**
     * @return
     */

    @com.fasterxml.jackson.annotation.JsonProperty("code")
    public Integer getCode() {
        return this.code;
    }

    /**
     * @param code
     * @return Returns a reference to this object so that method calls can be chained together.
     */

    public NotFoundException code(Integer code) {
        setCode(code);
        return this;
    }

}
