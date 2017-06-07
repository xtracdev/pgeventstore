@pgeventpub
Feature: Event Publishing
  Scenario:
    Given an evironment with event publishing disabled
    When I store an aggregate
    Then no events are written to the publish table

  Scenario:
    Given an environment with event publishing enabled
    When I store a new aggregate
    Then the events are written to the publish table

  Scenario:
    Given an environment with event publishing enabled
    When I republish the events
    Then all the events are written to the publish table