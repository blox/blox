@integ @environment @list-environments

Feature: Integration tests of ListEnvironment API

    Scenario: ListEnvironment should return environments
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        When I create an environment with name "sulu-test-le-1" in the cluster using the task-definition
        Then the environment should be returned in ListEnvironments call

#TODO: add tests with next token
