@integ @deployment @create-deployment

Feature: CreateDeployment

    Scenario: CreateDeployment without token should succeed for active environment
        Given A cluster "env.ECS_CLUSTER" and asg "env.ECS_CLUSTER_ASG"
        And a registered "sleep" task-definition
        And I update the desired-capacity of cluster to 1 instances and wait for a max of 300 seconds
        And I stop the tasks running in cluster
        When I create an environment with name "sulu-test-cd-1" in the cluster using the task-definition
        And I call CreateDeployment API
        Then GetDeployment with created deployment should succeed
        And the deployment should have 1 task running within 300 seconds

    Scenario: CreateDeployment fails when created with the same token as an existing deployment
        Given A cluster "env.ECS_CLUSTER" and asg "env.ECS_CLUSTER_ASG"
        And a registered "sleep" task-definition
        And I update the desired-capacity of cluster to 1 instances and wait for a max of 300 seconds
        And I stop the tasks running in cluster
        When I create an environment with name "sulu-test-cd-2" in the cluster using the task-definition
        And I call CreateDeployment API
        Then creating another deployment with the same token should fail
