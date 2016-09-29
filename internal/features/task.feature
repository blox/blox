@task
Feature: Task APIs

  Scenario: Get task
    Given I put 1 tasks in the queue
    When I get task with the same arn
    Then I get a task that matches the task I pushed to the queue


  Scenario: List tasks
    Given I put 2 tasks in the queue
    When I list tasks
    Then I get a list of tasks that includes the tasks I pushed to the queue

  Scenario Outline: Filter tasks by status
    Given I put the following tasks in the queue: 1 pending, 3 running, and 2 stopped
    When I filter tasks by <status>
    Then I get <expected> number of tasks
    And each task matches the corresponding task pushed to the queue

    Examples:
      | status   | expected |
      | pending  |    1     |
      | running  |    3     |
      | stopped  |    2     |
