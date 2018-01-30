@dataservice
@create-environment
Feature: Create environment

  Scenario: Create an environment
    When I create an environment
    Then the created environment response is valid

  Scenario: Create an environment that already exists
    Given I create an environment named "test" in the cluster "testCluster"
    When I try to create another environment with the name "test" in the cluster "testCluster"
    Then there should be a ResourceExistsException thrown
    And its resourceType is "environment"
    And its resourceId contains "test"

  Scenario: Create an environment that has the same name as another environment in a different cluster
    Given I create an environment named "test"
    When I try to create another environment with the name "test" in the cluster "anotherCluster"
    Then the created environment response is valid

  #TODO: Add invalid parameter tests