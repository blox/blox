@e2e @task @stream-tasks
Feature: Stream Tasks

  Scenario: Get task stream
    When I start streaming all task events
    And I start 1 task in the ECS cluster
    Then the stream tasks response contains at least 1 task
    And the stream tasks response contains the task with desired status running

  Scenario: Get task stream with past entity version
    Given I start 1 task in the ECS cluster
    And I get task with the cluster name and task ARN
    And I stop the 1 task in the ECS cluster
    When I start streaming all task events with past entity version
    Then the stream tasks response contains at least 1 task
    And the stream tasks response contains the task with desired status running
    And the stream tasks response contains the task with desired status stopped
