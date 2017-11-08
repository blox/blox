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
package steps.wrappers;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
public class DataServiceWrapper {

  private final DataService dataService;

  //TODO: implement data service methods used in tests
}
