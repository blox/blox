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

import java.util.Collections;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import lombok.Value;
import software.amazon.awssdk.services.cloudformation.CloudFormationClient;
import software.amazon.awssdk.services.cloudformation.model.DescribeStacksRequest;
import software.amazon.awssdk.services.cloudformation.model.DescribeStacksResponse;
import software.amazon.awssdk.services.cloudformation.model.Output;
import software.amazon.awssdk.services.cloudformation.model.Parameter;
import software.amazon.awssdk.services.cloudformation.model.Stack;

/**
 * Read-only CloudFormation wrapper to easily retrieve stack outputs from all stacks in an account.
 */
@RequiredArgsConstructor
public class CloudFormationStacks {
  private final CloudFormationClient cloudformation;
  private Map<String, CfnStack> stacks = null;

  public void refresh() {
    DescribeStacksResponse stacks =
        cloudformation.describeStacks(DescribeStacksRequest.builder().build());

    this.stacks =
        stacks
            .stacks()
            .stream()
            .map(CfnStack::from)
            .collect(Collectors.toMap(CfnStack::getName, Function.identity()));
  }

  public CfnStack get(String name) {
    if (stacks == null) refresh();

    if (!stacks.containsKey(name)) {
      throw new IndexOutOfBoundsException("No such stack: " + name);
    }

    return stacks.get(name);
  }

  @Value
  public static class CfnStack {
    private final String name;
    private final Map<String, String> parameters;
    private final Map<String, String> outputs;

    static CfnStack from(Stack s) {
      Map<String, String> parameters =
          s.parameters() == null
              ? Collections.emptyMap()
              : s.parameters()
                  .stream()
                  .collect(Collectors.toMap(Parameter::parameterKey, Parameter::parameterValue));

      Map<String, String> outputs =
          s.outputs() == null
              ? Collections.emptyMap()
              : s.outputs()
                  .stream()
                  .collect(Collectors.toMap(Output::outputKey, Output::outputValue));

      return new CfnStack(s.stackName(), parameters, outputs);
    }

    public String output(String key) {
      return outputs.get(key);
    }

    public String parameter(String key) {
      return parameters.get(key);
    }
  }
}
