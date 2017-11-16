@ignore
@dataservice
@describe-environment
Feature: Describe environment

  Scenario: Describe an existing environment
    Given I create an environment
    When I describe the environment
    Then the created and described environments match

  Scenario: Describe a non-existent environment
    When I describe an environment
    Then there should be an EnvironmentNotFoundException thrown

  #TODO: Add invalid parameter tests