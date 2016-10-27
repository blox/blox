@e2e @task @get-task
Feature: Get Task

  Scenario: Get task with a valid task ARN
    Given I start 1 task in the ECS cluster
    When I get task with the same ARN
    Then I get a task that matches the task started
