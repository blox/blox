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
