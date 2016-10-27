@e2e @instance @get-instance
Feature: Get Instance

  Scenario: Get instance with a valid instance ARN
    Given I have an instance registered with the ECS cluster
    When I get instance with the instance ARN
    Then I get an instance that matches the registered instance
