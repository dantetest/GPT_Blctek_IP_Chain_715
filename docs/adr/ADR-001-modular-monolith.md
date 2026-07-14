# ADR-001: Modular Monolith

- Status: Accepted
- Date: 2026-07-15

## Context

The product has coupled identity, dataset, order, payment, delivery, dispute, and evidence workflows, while the core team is small.

## Decision

Use one modular Go API and one Worker. Enforce package boundaries and domain-owned repositories. Split independent processes only for Data Agent, private Tracker, quality/dedup engines, and MCP.

## Consequences

This minimizes deployment and distributed-transaction risk. Modules may be extracted only after a measured scaling, security, or failure-isolation requirement appears.
