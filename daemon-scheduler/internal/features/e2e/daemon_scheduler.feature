@e2e

Feature:

    Scenario: Deployments should start tasks
        Given A cluster "env.ECS_CLUSTER" and asg "env.ECS_CLUSTER_ASG"
        And a registered "sleep" task-definition
        When I create an environment with name "sulu-test-daemon" in the cluster using the task-definition
        And I stop the tasks running in cluster
        And I update the desired-capacity of cluster to 1 instances and wait for a max of 300 seconds
        And I call CreateDeployment API
        Then the deployment should have 1 task running within 300 seconds
        And the deployment should complete in 100 seconds
        When I update the desired-capacity of cluster to 2 instances and wait for a max of 300 seconds
        Then the deployment should have 2 tasks running within 120 seconds
        When I stop the tasks running in cluster
        Then the deployment should have 2 tasks running within 120 seconds
        When I update the desired-capacity of cluster to 10 instances and wait for a max of 300 seconds
        Then the deployment should have 10 tasks running within 120 seconds
        When I update the desired-capacity of cluster to 0 instances and wait for a max of 300 seconds
        Then the deployment should have 0 tasks running within 120 seconds
