@ignore
@dataservice
@describe-environment-revision
Feature: Describe environment revision

  Scenario: Describe an existing environment revision
    Given I create an environment
    And I create an environment revision
    When I describe the environment revision
    Then the created and described environment revisions match

  Scenario: Describe a non-existent environment revision
    When I describe an environment revision
    Then there should be an EnvironmentRevisionNotFoundException thrown

  #TODO: Add invalid parameter tests