@ignore
@dataservice
@list-clusters
Feature: List clusters

  Scenario: List clusters when none exist
    When I list clusters
    Then 0 clusters are returned

  Scenario: List clusters when multiple environments exist on one cluster
    Given I create an environment
    And I create another environment with the same cluster
    When I list clusters
    Then 1 cluster is returned

  Scenario: List clusters when multiple environments exist on different clusters
    Given I create an environment
    And I create another environment with a different cluster
    When I list clusters
    Then 2 clusters are returned

  #TODO: add next token
  #Scenario: List clusters when the number of clusters exceeds max results

  # TODO: Add invalid parameter tests