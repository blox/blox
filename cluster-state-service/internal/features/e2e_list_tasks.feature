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

  Scenario: List tasks with invalid status filter
    When I try to list tasks with an invalid status filter
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "Invalid status"

  Scenario: List tasks with invalid cluster filter
    When I try to list tasks with an invalid cluster filter
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "Invalid cluster ARN or name"
