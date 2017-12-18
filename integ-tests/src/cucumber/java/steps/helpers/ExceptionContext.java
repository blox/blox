/*
 * Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"). You may
 * not use this file except in compliance with the License. A copy of the
 * License is located at
 *
 *     http://aws.amazon.com/apache2.0/
 *
 * or in the "LICENSE" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */
package steps.helpers;

/**
 * This class is used by Cucumber tests to keep track of any exceptions that tests need to check
 * for. It requires that exceptions are explicitly checked before storing another one. This ensures
 * that an exception isn't accidentally ignored in a test. If an attempt to store an exception is
 * made when the context already contains an exception that hasn't been checked an
 * IllegalStateException is thrown.
 */
public class ExceptionContext {
  private Exception exception;
  private boolean hasUnhandledException;

  /**
   * Gets the last exception that was thrown and marks it as handled, i.e. after calling this method
   * hasUnhandledException() will return false.
   */
  public Exception getException() {
    hasUnhandledException = false;
    return exception;
  }

  /**
   * Sets the exception or throws an IllegalStateException if the context already contains an
   * exception that hasn't been retrieved via getException().
   */
  public void setException(Exception exception) {
    if (hasUnhandledException) {
      hasUnhandledException = false;
      throw new IllegalStateException("Unhandled exception in Context", this.exception);
    }

    this.exception = exception;
    hasUnhandledException = true;
  }

  public boolean hasUnhandledException() {
    return hasUnhandledException;
  }
}
