# ADR-003: Transactional Outbox and Provider Inbox

- Status: Accepted
- Date: 2026-07-15

## Context

Payment, settlement, evidence, and notification calls can time out or deliver duplicate callbacks.

## Decision

Write outbound events in the same transaction as domain changes. Persist inbound Provider callbacks before processing. Both paths are idempotent, observable, retryable, and support a dead-letter state.

## Consequences

External calls are eventually consistent. UI and APIs must expose processing states rather than claim immediate completion.
