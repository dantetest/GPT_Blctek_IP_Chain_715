# ADR-002: Provider Interface Pattern

- Status: Accepted
- Date: 2026-07-15

## Context

支付宝、KYC、存证、邮件和对象存储 credentials are unavailable during early construction, but product flows must remain testable.

## Decision

Each external system is accessed through a narrow application interface with Mock and production adapters. Provider callbacks enter through a persistent Inbox. Production startup rejects Mock financial, KYC, and evidence providers.

## Consequences

End-to-end tests remain deterministic. Production adapters cannot silently change domain states outside the shared application services.
