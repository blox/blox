@e2e @instance @list-instances
Feature: List Instances

  Scenario Outline: List instances
    Given I have some instances registered with the ECS cluster
    When I list instances
    Then the list instances response contains all the registered instances
