@e2e @task @filter-tasks
Feature: Filter Tasks

  Scenario Outline: Filter tasks by running status
    Given I start <count> tasks in the ECS cluster
    When I filter tasks by running status
    Then the filter tasks response contains at least <count> tasks
    And all <count> tasks are present in the filter tasks response

    Examples:
      | count |
      |   1   |

  Scenario Outline: Filter tasks by stopped status
    Given I start <count> tasks in the ECS cluster
    And I stop the <count> tasks in the ECS cluster
    When I filter tasks by <stopped> status
    Then the filter tasks response contains at least <count> tasks
    And all <count> tasks are present in the filter tasks response

    Examples:
      | count | stopped |
      |   1   | stopped |
      |   1   | STOPPED |

  Scenario Outline: Filter tasks by cluster
    Given I start <count> tasks in the ECS cluster
    When I filter tasks by the ECS cluster name
    Then the filter tasks response contains at least <count> tasks
    And all <count> tasks are present in the filter tasks response

    Examples:
      | count |
      |   1   |
