@integ @environment @get_environment

Feature: Integration tests of GetEnvironment API

    Scenario: GetEnvironment should return environment details
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "blox-test-ge-1" in the cluster using the task-definition
        Then GetEnvironment should succeed

    Scenario: GetEnvironment with invalid environment should return NotFound
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "blox-test-ge-2" in the cluster using the task-definition
        Then GetEnvironment with name "blox-test-missing" should fail with NotFound
