@e2e @instance @stream-instances
Feature: Stream Instances

  Scenario: Get instance stream
    When I start streaming all instance events
    And I start 1 task in the ECS cluster
    And I get instance where the task was started
    Then the stream instances response contains at least 1 instance
    And the stream instances response contains the instance where the task was started

  Scenario: Get instance stream with past entity version
    Given I start 1 task in the ECS cluster
    And I get instance where the task was started
    And I stop the 1 task in the ECS cluster
    When I start streaming all instance events with past entity version
    Then the stream instances response contains at least 1 instance
    And the stream instances response contains the instance where the task was started
