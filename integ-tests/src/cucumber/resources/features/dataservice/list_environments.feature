@ignore
@dataservice
@list-environments
Feature: List environments

  Scenario: List environments on a non-existent cluster
    When I list environments on a cluster that does not exist
    Then 0 environments are returned

  Scenario: List environments when only one environment exists on that cluster
    Given I create an environment
    And I create another environment with a different cluster
    When I list environments on the first cluster
    Then 1 environment is returned

  Scenario: List environments when multiple environments exist on that cluster
    Given I create an environment
    And I create another environment with the same cluster
    And I create another environment with a different cluster
    When I list environments on the first cluster
    Then 2 environments are returned

  #TODO: add next token
  #Scenario: List environments when the number of environments exceeds max results

  # TODO: Add invalid parameter tests