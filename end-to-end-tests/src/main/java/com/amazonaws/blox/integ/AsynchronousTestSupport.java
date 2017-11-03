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
package com.amazonaws.blox.integ;

public class AsynchronousTestSupport {

  private AsynchronousTestSupport() {}

  /**
   * Re-run the given block until it no longer raises an AssertionError, or it times out.
   *
   * @param timeout timeout in ms
   * @param block block to run
   */
  public static void waitOrTimeout(long timeout, Runnable block) throws InterruptedException {
    long startTime = System.currentTimeMillis();
    long duration;

    while (true) {

      try {
        Thread.sleep(1000);

        block.run();
        return;
      } catch (AssertionError e) {
        duration = System.currentTimeMillis() - startTime;

        if (duration >= timeout) {
          throw e;
        }
      }
    }
  }
}
