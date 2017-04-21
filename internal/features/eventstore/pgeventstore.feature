@pgeventstore
Feature: Event Store

  Scenario: No max for new aggregate
    Given a new aggregate instance
    When we check the max version in the event store
    Then the max version is 0

