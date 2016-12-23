@e2e @task @list-tasks
Feature: List Tasks

  Scenario Outline: List tasks
    Given I start <count> tasks in the ECS cluster
    When I list tasks
    Then the list tasks response contains at least <count> tasks
    And all <count> tasks are present in the list tasks response

    Examples:
      | count |
      |   1   |
      |   3   |

  Scenario: List tasks with status and cluster filter returns tasks
    Given I start 1 task in the ECS cluster
    When I list tasks with filters set to running status and cluster name
    Then the list tasks response contains at least 1 task
    And all tasks in the list tasks response belong to the cluster and have status set to running

  Scenario: List tasks with status and cluster filter returns no tasks
    Given I start 1 task in the ECS cluster
    When I list tasks with filters set to running status and a different cluster name
    Then the list tasks response contains 0 tasks

  Scenario: List tasks with invalid status filter
    When I try to list tasks with an invalid status filter
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "Invalid status"

  Scenario: List tasks with invalid cluster filter
    When I try to list tasks with an invalid cluster filter
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "Invalid cluster ARN or name"
