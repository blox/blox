@e2e @instance @stream-instances
Feature: Stream Instances

  Scenario: Get instance stream
    When I start streaming all instance events
    And I start 1 task in the ECS cluster
    Then the stream instances response contains at least 1 instance
    And the stream instances response contains the cluster where the task was started