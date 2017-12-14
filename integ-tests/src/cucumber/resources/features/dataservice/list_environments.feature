@ignore
@dataservice
@list-environments
Feature: List environments

  Scenario: List environments on a non-existent cluster
    When I list environments on a cluster that does not exist
    Then 0 environments are returned

  Scenario: List environments when only one environment exists on that cluster
    Given I create an environment named "test" in the cluster "testCluster"
    And I create another environment named "anotherenv" in the cluster "differentCluster"
    When I list environments on the first cluster
    Then 1 environment is returned

  Scenario: List environments when multiple environments exist on that cluster
    Given I create an environment named "test" in the cluster "testCluster"
    And I create another environment named "anotherenv" in the cluster "testCluster"
    And I create another environment named "anotherenv" in the cluster "differentCluster"
    When I list environments on the first cluster
    Then 2 environments are returned

  #TODO: add next token
  #Scenario: List environments when the number of environments exceeds max results

  # TODO: Add invalid parameter tests