# Contributing to Blox
Contributions to Blox should be made via GitHub pull requests and discussed using GitHub issues.

### Before you start
If you would like to make a significant change, it's a good idea to first open an issue to discuss it.

### Making the request
Development takes place against the dev branch of this repository and pull requests should be opened against that branch.

### Code Style
This project follows the [Google Java Styleguide](https://google.github.io/styleguide/javaguide.html), and this style is enforced as part of the `check` task. We recommend you install [the `google-java-format` plugin for your IDE](https://github.com/google/google-java-format), or use the `gradle spotlessApply` task to format code before checking in.

### Testing
Any contributions should pass all tests. You can run all tests by running `gradle check` from the project root.

### Licensing
Blox is released under an Apache 2.0 license. Any code that you submit will be released under that license.
For significant changes, we may ask you to sign a Contributor License Agreement.
