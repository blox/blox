@integ @environment @list-environments

Feature: Integration tests of ListEnvironment API

    Scenario: ListEnvironments should return environments
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "sulu-test-le-1" in the cluster using the task-definition
        Then the environment should be returned in ListEnvironments call

    Scenario: ListEnvironments with cluster filter should return environments
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        And I create an environment with name "sulu-test-fe-1" in the cluster using the task-definition
        And another cluster "ecs-daemon-scheduler-test2"
        And a registered "sleep" task-definition
        And I create an environment with name "sulu-test-fe-2" in the cluster using the task-definition
        Then there should be at least 1 environment returned when I call ListEnvironments with cluster filter set to the second cluster
        And all the environments in the response should correspond to the second cluster
        And second environment should be one of the environments in the response

#TODO: add tests with next token
