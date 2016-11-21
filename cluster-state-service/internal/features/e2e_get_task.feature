@e2e @task @get-task
Feature: Get Task

  Scenario: Get task with a valid task ARN
    Given I start 1 task in the ECS cluster
    When I get task with the cluster name and task ARN
    Then I get a task that matches the task started

  Scenario: Get non-existent task
    When I try to get task with a non-existent ARN
    Then I get a GetTaskNotFound task exception
