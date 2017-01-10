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

  Scenario Outline: List tasks with cluster filter
    Given I start <count> tasks in the ECS cluster
    When I list tasks with cluster filter set to the ECS cluster name
    Then the list tasks response contains at least <count> tasks
    And all <count> tasks are present in the list tasks response

    Examples:
      | count |
      |   1   |

  Scenario Outline: List tasks with status filter stopped
    Given I start <count> tasks in the ECS cluster
    And I stop the <count> tasks in the ECS cluster
    When I list tasks with status filter set to <stopped>
    Then the list tasks response contains at least <count> tasks
    And all <count> tasks are present in the list tasks response

    Examples:
      | count | stopped |
      |   1   | stopped |
      |   1   | STOPPED |

  Scenario Outline: List tasks with status filter running
    Given I start <count> tasks in the ECS cluster
    When I list tasks with status filter set to running
    Then the list tasks response contains at least <count> tasks
    And all <count> tasks are present in the list tasks response

    Examples:
      | count |
      |   1   |

  Scenario Outline: List tasks with startedBy filter
    Given I start <count> tasks in the ECS cluster with startedBy set to e2eTester
    When I list tasks with startedBy filter set to e2eTester
    Then the list tasks response contains at least <count> tasks
    And all <count> tasks are present in the list tasks response

    Examples:
      | count |
      |   1   |

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

  Scenario: List tasks with redundant filters
    When I try to list tasks with redundant filters
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "At least one of the filters provided is specified multiple times"

  Scenario: List tasks with invalid filter combination
    When I try to list tasks with status, cluster and startedBy filters
    Then I get a ListTasksBadRequest task exception
    And the task exception message contains "The combination of filters provided are not supported"
