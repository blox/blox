# Blox: frontend-service
This is a Spring MVC/REST webservice, running as a single Lambda function, fronted by API Gateway. For more details on how this is structured, see the [documentation](../docs/frontend_design.md).

### Testing
To run the unit tests for only this project, run the following from the repository root:

```
./gradlew frontend-service:check
```
