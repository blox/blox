@integ @deployment @delete-environment

Feature: DeleteEnvironment

    Scenario: Delete existing environment
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        And I create an environment with name "blox-test-daemon-de-1" in the cluster using the task-definition
        When I delete the environment
        Then get environment should return empty

    Scenario: Delete non-existing environment
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        And I create an environment with name "blox-test-daemon-de-2" in the cluster using the task-definition
        When I delete the environment
        Then get environment should return empty
        And deleting the environment again should succeed
