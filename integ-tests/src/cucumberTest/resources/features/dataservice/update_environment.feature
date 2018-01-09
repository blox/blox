@ignore
@dataservice
@update-environment
Feature: Update environment

  Scenario: Update an environment
    Given I create an environment named "test"
    When I update the created environment named "test"
    Then the update environment response is valid

  Scenario: Update a non-existent environment
    When I try to update a non-existent environment named "test"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment"
    And the resourceId contains "test"

  #TODO: Add invalid parameter tests