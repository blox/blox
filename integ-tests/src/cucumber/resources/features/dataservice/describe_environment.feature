@ignore
@dataservice
@describe-environment
Feature: Describe environment

  Scenario: Describe a newly created environment
    Given I create an environment named "test"
    When I describe the environment named "test"
    Then the created and described environments match

  Scenario: Describe an updated environment
    Given I create an environment named "test"
    And I update the environment named "test"
    When I describe the environment named "test"
    Then the updated and described environments match

  Scenario: Describe a non-existent environment
    When I try to describe a non-existent environment named "non-existent"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment"
    And the resourceId contains "non-existent"

  Scenario: Describe a deleted environment
    Given I create an environment named "test"
    And I delete the created environment named "test"
    When I try to describe the environment named "test"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment"
    And the resourceId contains "test"

  #TODO: Add invalid parameter tests