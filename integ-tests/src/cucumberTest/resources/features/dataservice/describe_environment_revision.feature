@ignore
@dataservice
@describe-environment-revision
Feature: Describe environment revision

  Scenario: Describe a newly created environment revision
    Given I create an environment named "test":
      | taskDefinition | cluster       |
      | "testTaskDef"  | "testCluster" |
    And I list the environment revisions for environment named "test"
    When I describe the first environment revision with the returned environment revision id
    Then the described environment revision has:
      | taskDefinition | cluster       |
      | "testTaskDef"  | "testCluster" |

  Scenario: Describe an updated environment revision
    Given I create an environment named "test":
      | taskDefinition | cluster       |
      | "testTaskDef"  | "testCluster" |
    And I update the environment named "test":
      | taskDefinition | cluster        |
      | "testTaskDef2"  | "testCluster" |
    And I list the environment revisions for environment named "test"
    When I describe the latest environment revision with the returned environment revision id
    Then the described environment revision has:
      | taskDefinition | cluster       |
      | "testTaskDef2"  | "testCluster" |

  Scenario: Describe a non-existent environment revision
    When I try to describe a non-existent environment revision with id "non-existent"
    Then there should be a ResourceNotFoundException thrown
    And the resourceType is "environment revision"
    And the resourceId is "non-existent"

  #TODO: Add invalid parameter tests