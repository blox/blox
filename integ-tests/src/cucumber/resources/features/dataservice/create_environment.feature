@ignore
@dataservice
@create-environment
Feature: Create environment

  Scenario: Create an environment
    When I create an environment
    Then the created environment response is valid

  Scenario: Create an environment that already exists
    Given I create an environment
    When I try to create another environment with the same name in the same cluster
    Then there should be an EnvironmentExistsException thrown

  Scenario: Create an environment that has the same name as another environment in a different cluster
    Given I create an environment
    When I try to create another environment with the same name in a different cluster
    Then the created environment response is valid

  #TODO: Add invalid parameter tests