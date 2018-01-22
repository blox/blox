Feature: Deploying a Daemon environment to a cluster with Replace After Terminate

  Background:
    Given a cluster named "TestCluster"
    And a Daemon environment named "DaemonEnvironment":
      | clusterName | deploymentMethod      | taskDefinitionArn |
      | TestCluster | ReplaceAfterTerminate | v1                |


  Scenario: Scheduling on a cluster with no tasks
    Given the cluster has the following instances and tasks:
      | instance | tasks |
      | i-1      |       |
      | i-2      |       |
    When the scheduler runs
    Then it should start the following tasks:
      | clusterName | containerInstanceArn | taskDefinitionArn | group             |
      | TestCluster | i-1                  | v1                | DaemonEnvironment |
      | TestCluster | i-2                  | v1                | DaemonEnvironment |
    And it should not take any further actions

  Scenario: Scheduling on a cluster with all tasks up to date
    Given the cluster has the following instances and tasks:
      | instance | tasks          |
      | i-1      | t-1:v1:PENDING |
      | i-2      | t-2:v1:RUNNING |
    When the scheduler runs
    Then it should not take any further actions

  Scenario: Scheduling on a cluster with tasks that failed
    Given the cluster has the following instances and tasks:
      | instance | tasks          |
      | i-1      | t-1:v1:STOPPED |
      | i-2      | t-2:v1:RUNNING |
    When the scheduler runs
    Then it should start the following tasks:
      | containerInstanceArn | taskDefinitionArn | group             |
      | i-1                  | v1                | DaemonEnvironment |
    And it should not take any further actions
