@dataservice
@list-clusters
Feature: List clusters

  Background:
    Given I am using account ID 1234567890

  Scenario: List clusters when none exist
    # Given I have no environments
    When I list clusters
    Then no clusters are returned

  Scenario: List clusters when multiple environments exist on one cluster
    Given I create an environment named "EnvironmentA" in the cluster "Cluster"
    And I create an environment named "EnvironmentB" in the cluster "Cluster"
    When I list clusters
    Then these clusters are returned:
      | accountId  | clusterName |
      | 1234567890 | Cluster     |

  Scenario: List clusters when multiple environments exist on different clusters
    Given I create an environment named "EnvironmentA" in the cluster "ClusterA"
    And I create an environment named "EnvironmentB" in the cluster "ClusterB"
    When I list clusters
    Then these clusters are returned:
      | accountId  | clusterName |
      | 1234567890 | ClusterA    |
      | 1234567890 | ClusterB    |

  #TODO: add next token
  #Scenario: List clusters when the number of clusters exceeds max results

  # TODO: Add invalid parameter tests