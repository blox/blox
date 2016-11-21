@integ @deployment

Feature: Integration tests of Get/List Deployment API

    Scenario: GetDeployment should return deployment details
        Given A cluster "env.ECS_CLUSTER" and asg "env.ECS_CLUSTER_ASG"
        And a registered "sleep" task-definition
        And I update the desired-capacity of cluster to 1 instances and wait for a max of 300 seconds
        And I stop the tasks running in cluster
        And I create an environment with name "sulu-test-gd-1" in the cluster using the task-definition
        And I call CreateDeployment API
        Then GetDeployment with created deployment should succeed

    Scenario: GetDeployment with non-existent environment should fail with NotFound
        When I call GetDeployment with environment "sulu-test-missing", it should fail with NotFound

    Scenario: GetDeployment with invalid deploymentID should return NotFound
        Given a cluster "ecs-daemon-scheduler-test"
        And a registered "sleep" task-definition
        And I create an environment with name "sulu-test-gd-2" in the cluster using the task-definition
        When I call GetDeployment with id "invalid-deployment-id", it should fail with NotFound
