@pgeventstore
Feature: Event Store

  Scenario: No max for new aggregate
    Given a new aggregate instance
    When we check the max version in the event store
    Then the max version is 0

  Scenario:
    Given a new aggregate instance
    When we get the max version from the event store
    And the max version is greater than the aggregate version
    Then a concurrency error is return on aggregate store