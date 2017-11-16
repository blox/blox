@ignore
@dataservice
@create-environment-revision
Feature: Create environment revision

  Scenario: Create an environment revision
    Given I create an environment
    When I create an environment revision
    Then the created environment revision response is valid

  Scenario: Create an environment revision for a non-existent environment
    When I create an environment revision
    Then there should be an EnvironmentNotFoundException thrown

  Scenario: Create an environment revision that already exists
    Given I create an environment
    And I create an environment revision
    When I try to create another environment revision with the same id and version
    Then there should be an EnvironmentRevisionExistsException thrown

  #TODO: Add invalid parameter tests