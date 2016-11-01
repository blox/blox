@e2e @instance @get-instance
Feature: Get Instance

  Scenario: Get instance with a valid instance ARN
    Given I have an instance registered with the ECS cluster
    When I get instance with the cluster name and instance ARN
    Then I get an instance that matches the registered instance

  Scenario: Get non-existent instance
    When I try to get instance with a non-existent ARN
    Then I get a GetInstanceNotFound instance exception
