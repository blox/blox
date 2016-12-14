@e2e @instance @list-instances
Feature: List Instances

  Scenario Outline: List instances
    Given I have some instances registered with the ECS cluster
    When I list instances
    Then the list instances response contains all the registered instances

  Scenario: List instances with invalid status filter
    When I try to list instances with an invalid status filter
    Then I get a ListInstancesBadRequest instance exception
    And the instance exception message contains "Invalid status"

  Scenario: List instances with invalid cluster filter
    When I try to list instances with an invalid cluster filter
    Then I get a ListInstancesBadRequest instance exception
    And the instance exception message contains "Invalid cluster ARN or name"
