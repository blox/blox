@ignore
@create-environment
Feature: Create environment

  Scenario: Create an environment
    When I create an environment
    Then the created environment response is valid
