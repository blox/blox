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
package com.amazonaws.blox.tasks;

import com.amazonaws.blox.swagger.GenerationTimestampFilter;
import com.amazonaws.blox.swagger.SwaggerFilter;
import com.github.kongchen.swagger.docgen.GenerateException;
import com.github.kongchen.swagger.docgen.reader.SpringMvcApiReader;
import io.swagger.models.Swagger;
import io.swagger.util.Yaml;
import java.io.File;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import lombok.Getter;
import lombok.Setter;
import org.apache.maven.plugin.logging.Log;
import org.apache.maven.plugin.logging.SystemStreamLog;
import org.gradle.api.DefaultTask;
import org.gradle.api.tasks.Classpath;
import org.gradle.api.tasks.Input;
import org.gradle.api.tasks.InputFiles;
import org.gradle.api.tasks.Nested;
import org.gradle.api.tasks.OutputFile;
import org.gradle.api.tasks.TaskAction;

/**
 * Generate a YAML-formatted Swagger specification from the given classes.
 *
 * <p>This does basically the same as the swagger task provided by
 * https://github.com/gigaSproule/swagger-gradle-plugin, but gives us an extension point to wire in
 * our own filters for modifying the resultant Swagger document.
 */
@Getter
@Setter
public class GenerateSwaggerModel extends DefaultTask {
  private Log log = new SystemStreamLog();

  public static final SwaggerFilter DEFAULT_FILTER = new GenerationTimestampFilter();

  /** A list of class names to scan for Swagger annotations */
  @Input private List<String> apiClasses = new ArrayList<>();

  /** The file into which to write the swagger definition */
  @OutputFile private File swaggerFile;

  /**
   * The classpath to use when loading classes to scan for swagger annotations
   *
   * <p>Typically, we need the Runtime classpath specifically so that we can load the classes being
   * built in the parent project, as well as all its dependencies (in particular the Java
   * annotations used to declare RESTful controllers).
   */
  @Classpath @InputFiles private Iterable<File> scanClasspath;

  /**
   * An ordered list of SwaggerFilter instances to apply to the generated Swagger definition before
   * writing it out.
   */
  @Nested private List<SwaggerFilter> filters = new ArrayList<>();

  public GenerateSwaggerModel() {
    filters.add(DEFAULT_FILTER);
  }

  @TaskAction
  public void generateSpec() throws IOException, ClassNotFoundException, GenerateException {
    ClassLoader loader = projectClassLoader();

    Set<Class<?>> classes = new HashSet<>();

    for (String name : apiClasses) {
      classes.add(loader.loadClass(name));
    }

    SpringMvcApiReader reader = new SpringMvcApiReader(new Swagger(), log);
    Swagger swagger = reader.read(classes);

    for (SwaggerFilter filter : filters) {
      filter.apply(swagger);
    }

    Yaml.pretty().writeValue(swaggerFile, swagger);
  }

  private ClassLoader projectClassLoader() throws MalformedURLException {
    Set<URL> urls = new HashSet<>();
    for (File file : getScanClasspath()) {
      urls.add(file.toURI().toURL());
    }

    return new URLClassLoader(
        urls.toArray(new URL[] {}), GenerateSwaggerModel.class.getClassLoader());
  }
}
