@integ @ping

Feature: Integration tests of Ping API

    Scenario: Ping should return OK
        When I make a Ping call
        Then the Ping response indicates that the service is healthy
