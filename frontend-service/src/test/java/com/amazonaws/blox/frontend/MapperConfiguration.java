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
package com.amazonaws.blox.frontend;

import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

/**
 * Spring configuration for wiring together all the MapStruct mappers.
 *
 * <p>With componentModel = "spring", mapstruct will not automatically wire together dependent
 * mapper instances from {@link org.mapstruct.factory.Mappers#getMapper(Class)}. The only way to
 * wire the depdendent mappers is through a spring application context.
 *
 * <p>Unfortunately, we can't just wire all the unit test collaborators (i.e. Controllers, etc)
 * using Spring, as that requires wiring all of the dependencies of the controller. We should
 * investigate this again further.
 */
@Configuration
@ComponentScan("com.amazonaws.blox.frontend.mappers")
public class MapperConfiguration {}
