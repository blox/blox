@e2e @instance @list-instances
Feature: List Instances

  Scenario Outline: List instances
    Given I have some instances registered with the ECS cluster
    When I list instances
    Then the list instances response contains all the registered instances

  Scenario Outline: List instances with cluster filter
    Given I have some instances registered with the ECS cluster
    When I list instances with cluster filter set to the ECS cluster name
    Then the list instances response contains all the instances registered with the cluster

  Scenario: List instances with invalid status filter
    When I try to list instances with an invalid status filter
    Then I get a ListInstancesBadRequest instance exception
    And the instance exception message contains "Invalid status"

  Scenario: List instances with invalid cluster filter
    When I try to list instances with an invalid cluster filter
    Then I get a ListInstancesBadRequest instance exception
    And the instance exception message contains "Invalid cluster ARN or name"

  Scenario: List instances with redundant filters
    When I try to list instances with redundant filters
    Then I get a ListInstancesBadRequest instance exception
    And the instance exception message contains "At least one of the filters provided is specified multiple times"
