Feature: Deploying a Daemon environment to a cluster with Replace After Terminate

  Background:
    Given a cluster named "TestCluster"
    And a Daemon environment named "DaemonEnvironment":
      | clusterName | deploymentMethod      | taskDefinitionArn |
      | TestCluster | ReplaceAfterTerminate | v2                |

  Scenario: Scheduling on a cluster with tasks that don't match the environment version
    Given the cluster has the following instances and tasks:
      | instance | tasks          |
      | i-1      | t-1:v1:RUNNING |
      | i-2      | t-2:v2:RUNNING |
    When the scheduler runs
    Then it should stop the following tasks:
      | clusterName | task           | reason                |
      | TestCluster | v1             | Stopped by deployment to DaemonEnvironment@v2 |
    And it should not take any further actions

  Scenario: Scheduling on a cluster with multiple actions
    Given the cluster has the following instances and tasks:
      | instance | tasks          |
      | i-1      | t-1:v1:RUNNING |
      | i-2      | t-2:v2:RUNNING |
      | i-3      |                |
    When the scheduler runs
    Then it should start the following tasks:
      | containerInstanceArn | taskDefinitionArn | group             |
      | i-3                  | v2                | DaemonEnvironment |
    Then it should stop the following tasks:
      | clusterName | task           | reason                |
      | TestCluster | v1             | Stopped by deployment to DaemonEnvironment@v2 |
    And it should not take any further actions