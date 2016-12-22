@e2e @task @stream-tasks
Feature: Stream Tasks

  Scenario: Get task stream
    When I start streaming all task events
    And I start 1 task in the ECS cluster
    Then the stream tasks response contains at least 1 task
    And the stream tasks response contains the task started