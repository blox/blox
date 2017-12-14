@ignore
@dataservice
@list-environment-revision
Feature: List environment revisions

  Scenario: List environment revisions for a non-existent environment
    When I list environment revisions for a non-existent environment named "non-existent"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment"
    And the resourceId contains "non-existent"

  Scenario: List environment revisions when only one environment revision exists for that environment
    Given I create an environment named "test"
    When I list environment revisions for the created environment
    Then 1 environment revision is returned

  Scenario: List environments when multiple environment revisions exist for that environment
    Given I create an environment named "test"
    And I update the environment named "test"
    When I list the environment revisions for environment named "test"
    Then 2 environment revisions are returned

  #TODO: add next token
  #Scenario: List environments when the number of environments exceeds max results

  # TODO: Add invalid parameter tests