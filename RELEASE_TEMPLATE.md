## Subject Line
[Blox] Deploy release - v$$blox_version$$

## Activity Details
This activity is to: Release Blox version v$$blox_version$$ by merging the release-$$blox_version$$ branch into the dev and master branches on GitHub, and pushing new Blox images up to Docker Hub.

The purpose of this change is to: Release the latest features and bug fixes for Blox to allow consumers to start using the new functionality. For more details about the specific changes being released, refer to the $$blox_version$$ release notes [here](https://github.com/blox/blox/blob/dev/CHANGELOG.md).

The Blox framework version is: v$$blox_version$$ with a git hash of $$github_hash$$

## Impact Details
What will happen if this release doesn't happen?  
Blox consumers will not be able to take advantage of the new features and bug fixes in Blox.

Why is this the correct time/day to complete the release?  
This corresponds to the completion date of all GitHub issues planned for this release.

Are there any related, prerequisite changes upon which this release hinges?  
No.

## Worst Case Scenario
What could happen if this change causes impact?  
Blox consumers could start running this new version of the Blox framework and encounter random errors leading to unexpected behavior.

Where are the most likely places this change will fail?  
Errors encountered building the new Blox images or pushing them up to Docker Hub.

## Hostname or Service
GitHub: https://github.com/blox/blox  
Docker Hub: https://hub.docker.com/u/bloxoss/

## Timeline / Activity Plan
Times are relative to the start of the release.

- 00:00 Change the release status to "In Progress".
- 00:05 Perform the Activity Checklist steps below.
- 00:25 Perform the Validation Checklist steps below.
- 01:00 Activity complete. Change the release status to "Complete".

Refer to the Validation Activities section for more details about each activity and validation checklist step.

## Activity Checklist
TODO
