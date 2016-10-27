@e2e @instance @filter-instances
Feature: Filter Instances

  Scenario Outline: Filter instances by cluster
    Given I have some instances registered with the ECS cluster
    When I filter instances by the same ECS cluster name
    Then the filter instances response contains all the instances registered with the cluster
