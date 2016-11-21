@integ @deployment

Feature: Integration tests of ListDeployment API

    Scenario: ListDeployment should return deployment details
        Given A cluster "env.ECS_CLUSTER" and asg "env.ECS_CLUSTER_ASG"
        And a registered "sleep" task-definition
        And I update the desired-capacity of cluster to 1 instances and wait for a max of 300 seconds
        And I stop the tasks running in cluster
        When I create an environment with name "sulu-test-ld-1" in the cluster using the task-definition
        And I call CreateDeployment API 3 times
        Then ListDeployments should return all the 3 deployments

    Scenario: ListDeployments with non-existent environment should fail with NotFound
        When I call ListDeployments with environment "sulu-test-missing", it should fail with NotFound
