@integ @environment @create_environment

Feature: Integration tests of CreateEnvironment API

    Scenario: CreateEnvironment with valid cluster and task-definition
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "sulu-test-ce-1" in the cluster using the task-definition
        Then GetEnvironment should succeed

    Scenario: If environment with name already exists then CreateEnvironment should fail
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "sulu-test-ce-2" in the cluster using the task-definition
        Then creating the same environment should fail with BadRequest

    Scenario: CreateEnvironment API should fail for inactive cluster
        Given a cluster "ecs-daemon-scheduler-test-inactive"
        And a registered "sleep" task-definition
        And I delete cluster
        When I create an environment with name "sulu-test-invalid-cluster" it should fail with NotFound

    Scenario: CreateEnvironment API should fail for inactive ECS task-definition
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep-to-deregister" task-definition
        And I deregister task-definition
        When I create an environment with name "sulu-test-invalid-td" it should fail with NotFound

    Scenario: CreateEnvironment API should fail for empty name
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name " " it should fail with BadRequest
