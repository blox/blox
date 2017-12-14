@ignore
@dataservice
@delete-environment
Feature: Delete environment

  Scenario: Delete an environment
    Given I create an environment
    When I delete the created environment
    Then the delete environment response is valid

  Scenario: Delete a non-existent environment
    When I try to delete a non-existent environment named "non-existent"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment"
    And the resourceId contains "non-existent"

  Scenario: Delete a deleted environment
    Given I create an environment
    When I delete the created environment
    When I try to delete the environment
    Then the delete environment response is valid

  #TODO: Add invalid parameter tests