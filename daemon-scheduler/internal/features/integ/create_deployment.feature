@integ @deployment @create-deployment

Feature: CreateDeployment

    Scenario: CreateDeployment without token should succeed for active environment
        Given A cluster named "env.ECS_CLUSTER"
        And a registered "sleep" task-definition
        When I create an environment with name "blox-test-cd-1" in the cluster using the task-definition
        And I call CreateDeployment API
        Then GetDeployment with created deployment should succeed

    Scenario: CreateDeployment fails when created with the same token as an existing deployment
        Given A cluster named "env.ECS_CLUSTER"
        And a registered "sleep" task-definition
        When I create an environment with name "blox-test-cd-2" in the cluster using the task-definition
        And I call CreateDeployment API
        Then creating another deployment with the same token should fail
