package com.amazonaws.blox

import com.amazonaws.blox.tasks.GenerateSwaggerModel
import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory
import org.gradle.testfixtures.ProjectBuilder
import org.junit.Before
import org.junit.Test

import static org.junit.Assert.assertEquals

class GenerateSwaggerModelTest {
    private File swaggerFile = new File("build/tmp/swagger.yml")

    private classpath = (TestController.class.classLoader as URLClassLoader).URLs.collect { new File(it.toString()) }

    @Before
    void deleteSwaggerFile() {
        if (swaggerFile.exists()) {
            swaggerFile.delete()
        }
    }

    @Test
    void generatesSwaggerModelFromGivenClasses() throws Exception {
        GenerateSwaggerModel task = ProjectBuilder.builder().build().task("swagger", type: GenerateSwaggerModel)

        task.scanClasspath = this.classpath
        task.apiClasses = ["com.amazonaws.blox.TestController"]
        task.swaggerFile = this.swaggerFile

        task.execute()

        JsonNode swagger = readSwaggerFile()

        assertEquals("2.0", swagger.get("swagger").asText())
        assertEquals("test-summary", swagger
                .get("paths")
                .get("/test/{name}")
                .get("get")
                .get("summary")
                .asText())
    }

    private JsonNode readSwaggerFile() {
        new ObjectMapper(new YAMLFactory()).readTree(this.swaggerFile)
    }
}
