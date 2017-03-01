@integ @deployment @list-deployments

Feature: Integration tests of ListDeployment API

    Scenario: ListDeployment should return deployment details
        Given A cluster named "env.ECS_CLUSTER"
        And a registered "sleep" task-definition
        When I create an environment with name "blox-test-ld-1" in the cluster using the task-definition
        # TODO: add multiple deployments when updatedEnvironment is available
        And I call CreateDeployment API
        Then ListDeployments should return 1 deployment

    Scenario: ListDeployments with non-existent environment should fail with NotFound
        When I call ListDeployments with environment "blox-test-missing", it should fail with NotFound
