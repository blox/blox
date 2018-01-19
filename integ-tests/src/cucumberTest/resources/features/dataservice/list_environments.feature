@dataservice
@list-environments
Feature: List environments

  @ignore
  Scenario: List environments on a non-existent cluster
    When I list environments in a cluster that does not exist
    Then 0 environments are returned

  Scenario: List environments when only one environment exists on that cluster
    Given I create an environment named "EnvironmentA" in the cluster "ClusterA"
    And I create an environment named "EnvironmentB" in the cluster "ClusterB"
    When I list environments in cluster "ClusterA"
    Then these environments are returned
      | environmentName | accountId    | cluster     |
      | EnvironmentA    | 123456789012 | ClusterA    |

  Scenario: List environments when multiple environments exist on that cluster
    Given I create an environment named "EnvironmentA" in the cluster "ClusterA"
    And I create an environment named "EnvironmentB" in the cluster "ClusterA"
    And I create an environment named "EnvironmentC" in the cluster "ClusterB"
    When I list environments in cluster "ClusterA"
    Then these environments are returned
      | environmentName | accountId    | cluster     |
      | EnvironmentA    | 123456789012 | ClusterA    |
      | EnvironmentB    | 123456789012 | ClusterA    |

  Scenario: List environments with specified name prefix when multiple environments exist on that cluster
    Given I create an environment named "EnvironmentA" in the cluster "ClusterA"
    And I create an environment named "EnvironmentB" in the cluster "ClusterA"
    And I create an environment named "EnvsA" in the cluster "ClusterA"
    When I list environments in cluster "ClusterA" with name prefix "Environment"
    Then these environments are returned
      | environmentName | accountId    | cluster     |
      | EnvironmentA    | 123456789012 | ClusterA    |
      | EnvironmentB    | 123456789012 | ClusterA    |

  #TODO: add next token
  #Scenario: List environments when the number of environments exceeds max results

  # TODO: Add invalid parameter tests
